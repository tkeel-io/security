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

	"github.com/tkeel-io/security/pkg/apiserver/config"
	"github.com/tkeel-io/security/pkg/apiserver/filters"
	"github.com/tkeel-io/security/pkg/constants"
	"github.com/tkeel-io/security/pkg/errcode"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func RegisterToRestContainer(c *restful.Container, conf *config.RBACConfig) error {
	var webservice *restful.WebService
	for _, v := range c.RegisteredWebServices() {
		if v.RootPath() == "v1" {
			webservice = v
			break
		}
	}
	if webservice == nil {
		webservice = &restful.WebService{}
		webservice.Path("v1").
			Produces(restful.MIME_JSON).
			Filter(filters.Auth())

		c.Add(webservice)
	}

	handler := newRBACHandler(conf)

	webservice.Route(webservice.POST("rbac/{tenant_id}/roles").
		To(handler.AddRoleInDomain).
		Doc(" Add a role in tenant ").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Reads(AddPermissionIn{}).
		Returns(http.StatusOK, errcode.ErrMsgOK, nil).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagRBAC}))

	webservice.Route(webservice.GET("rbac/{tenant_id}/roles").
		To(handler.RolesInDomain).
		Doc("Get role list in tenant").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagRBAC}))

	webservice.Route(webservice.DELETE("rbac/{tenant_id}/roles/{role}").
		To(handler.DeleteRoleInDomain).
		Doc("delete a role in tenant").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Param(webservice.PathParameter("role", "role'ID")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagRBAC}))

	webservice.Route(webservice.POST("rbac/{tenant_id}/{role}/permissions").
		To(handler.AddPermissionInRole).
		Doc("add permission for role").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Param(webservice.PathParameter("role", "role'ID")).
		Reads(AddPermissionIn{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagRBAC}))

	webservice.Route(webservice.DELETE("rbac/{tenant_id}/{role}/permissions").
		To(handler.DeletePermissionInRole).
		Doc("delete a permission for role ").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Param(webservice.PathParameter("role", "role'ID")).
		Param(webservice.QueryParameter("permission_object", "permission object")).
		Param(webservice.QueryParameter("permission_action", "permission action")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagRBAC}))

	webservice.Route(webservice.POST("rbac/{tenant_id}/users/roles").
		To(handler.AddRoleToUser).
		Doc("add roles for users").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Reads(AddRoleInDomainIn{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagRBAC}))

	webservice.Route(webservice.DELETE("rbac/{tenant_id}/users/{user_id}/roles/{role}").
		To(handler.DeleteRoleOnUser).
		Doc("delete a role for user").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Param(webservice.PathParameter("user_id", "user's ID")).
		Param(webservice.PathParameter("role", "role")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagRBAC}))

	webservice.Route(webservice.GET("rbac/{tenant_id}/users/{user_id}/permissions").
		To(handler.UserPermissions).
		Doc("get user permissions ").
		Param(webservice.PathParameter("tenant_id", "tenant's ID")).
		Param(webservice.PathParameter("user_id", "user's ID")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagRBAC}))

	webservice.Route(webservice.POST("rbac/permission/check").
		To(handler.PermissionCheck).
		Doc("delete a role in tenant").
		Reads(PermissionCheck{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagRBAC}))

	return nil
}
