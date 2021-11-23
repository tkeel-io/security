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
	"net/http"

	"github.com/tkeel-io/security/pkg/apirouter"
	"github.com/tkeel-io/security/pkg/constants"
	"github.com/tkeel-io/security/pkg/errcode"
	"github.com/tkeel-io/security/pkg/models/dao"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func AddToRestContainer(c *restful.Container) error {
	webservice := apirouter.GetWebserviceWithPatch(c, "/v1/tenant")

	handler := newTenantHandler()

	webservice.Route(webservice.POST("").
		To(handler.Create).
		Doc(" Create a tenant").
		Reads(TenantCreteIn{}).
		Returns(http.StatusOK, errcode.ErrMsgOK, TenantCreateOut{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagTenant}))

	webservice.Route(webservice.GET("").
		To(handler.Query).
		Doc("get tenants").
		Param(webservice.QueryParameter("tenant_id", "").Required(false)).
		Param(webservice.QueryParameter("title", "").Required(false)).
		Returns(http.StatusOK, errcode.ErrMsgOK, []*dao.Tenant{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagTenant}))

	webservice.Route(webservice.DELETE("/tenant/{tenant_id}").
		To(handler.Delete).
		Doc("delete a tenant").
		Param(webservice.PathParameter("tenant_id", "tenant's ID").Required(true)).
		Returns(http.StatusOK, errcode.ErrMsgOK, nil).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagTenant}))

	webservice.Route(webservice.POST("/tenant/users").
		To(handler.UserCreate).
		Doc("create a user").
		Reads(UserCreateIn{}).
		Returns(http.StatusOK, errcode.ErrMsgOK, dao.User{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagTenant}))

	webservice.Route(webservice.GET("/tenant/users").
		To(handler.UserQuery).
		Doc("get users").
		Returns(http.StatusOK, errcode.ErrMsgOK, []dao.User{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagTenant}))

	webservice.Route(webservice.DELETE("/tenant/users/{user_id}").
		To(handler.UserDelete).
		Doc("delete a  users").
		Param(webservice.PathParameter("user_id", "").Required(true)).
		Returns(http.StatusOK, errcode.ErrMsgOK, nil).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagTenant}))

	return nil
}
