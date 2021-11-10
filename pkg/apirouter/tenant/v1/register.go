/*
Copyright 2021 The tKeel Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/tkeel-io/security/pkg/constants"
	"github.com/tkeel-io/security/pkg/errcode"
	"github.com/tkeel-io/security/pkg/models/dao"
	"net/http"
)

func AddToRestContainer(c *restful.Container) error {

	webservice := &restful.WebService{}
	webservice.Path("/tenant").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	handler := newTenantHandler()

	webservice.Route(webservice.POST("/").
		To(handler.Create).
		Doc(" Create a tenant").
		Reads(TenantCreteIn{}).
		Returns(http.StatusOK, errcode.ErrMsgOK, TenantCreateOut{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagTenant}))

	webservice.Route(webservice.GET("/").
		To(handler.Query).
		Doc("get tenants").
		Param(webservice.QueryParameter("tenant_id", "").Required(false)).
		Param(webservice.QueryParameter("title", "").Required(false)).
		Returns(http.StatusOK, errcode.ErrMsgOK, []*dao.Tenant{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagTenant}))

	webservice.Route(webservice.DELETE("/{tenant_id}").
		To(handler.Delete).
		Doc("delete a tenant").
		Param(webservice.PathParameter("tenant_id", "tenant's ID").Required(true)).
		Returns(http.StatusOK, errcode.ErrMsgOK, nil).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagTenant}))

	webservice.Route(webservice.POST("/users").
		To(handler.UserCreate).
		Doc("create a user").
		Reads(UserCreateIn{}).
		Returns(http.StatusOK, errcode.ErrMsgOK, dao.User{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagTenant}))

	webservice.Route(webservice.GET("/users").
		To(handler.UserQuery).
		Doc("get users").
		Returns(http.StatusOK, errcode.ErrMsgOK, []dao.User{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagTenant}))

	webservice.Route(webservice.DELETE("/users/{user_id}").
		To(handler.UserDelete).
		Doc("delete a  users").
		Param(webservice.PathParameter("user_id", "").Required(true)).
		Returns(http.StatusOK, errcode.ErrMsgOK, nil).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagTenant}))

	c.Add(webservice)
	return nil
}
