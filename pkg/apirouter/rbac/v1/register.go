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
	"github.com/tkeel-io/security/pkg/constants"
	"net/http"

	"github.com/tkeel-io/security/pkg/apiserver/filters"
	"github.com/tkeel-io/security/pkg/errcode"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func AddToRestContainer(c *restful.Container) error {

	webservice := &restful.WebService{}
	webservice.Path("/rbac").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Filter(filters.Auth())

	handler := newRBACHandler()

	webservice.Route(webservice.POST("/{tenant_id}/roles").
		To(handler.AddRoleInDomain).
		Doc(" Add a role in tenant ").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Reads(AddPermissionIn{}).
		Returns(http.StatusOK, errcode.ErrMsgOK, nil).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagRBAC}))

	webservice.Route(webservice.GET("/{tenant_id}/roles").
		To(handler.RolesInDomain).
		Doc("Get role list in tenant").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagRBAC}))

	webservice.Route(webservice.DELETE("/{tenant_id}/roles/{role}").
		To(handler.DeleteRoleInDomain).
		Doc("delete a role in tenant").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Param(webservice.PathParameter("role", "role'ID")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagRBAC}))

	webservice.Route(webservice.POST("/{tenant_id}/{role}/permissions").
		To(handler.AddPermissionInRole).
		Doc("delete a role in tenant").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Param(webservice.PathParameter("role", "role'ID")).
		Reads(AddPermissionIn{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagRBAC}))

	webservice.Route(webservice.GET("/{tenant_id}/users/{user}/permissions").
		To(handler.PermissionsInUser).
		Doc("delete a role in tenant").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Param(webservice.PathParameter("user", "user'ID")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagRBAC}))

	webservice.Route(webservice.DELETE("/{tenant_id}/{role}/permissions/{permission_id}").
		To(handler.DeletePermissionInRole).
		Doc("delete a role in tenant").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Param(webservice.PathParameter("permission_id", "permissions'ID")).
		Param(webservice.PathParameter("role", "role'ID")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagRBAC}))

	webservice.Route(webservice.POST("/permission/check").
		To(handler.PermissionCheck).
		Doc("delete a role in tenant").
		Reads(PermissionCheck{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagRBAC}))

	c.Add(webservice)
	return nil
}
