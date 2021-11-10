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
	"github.com/tkeel-io/security/pkg/apiserver/response"
	"github.com/tkeel-io/security/pkg/errcode"
	"github.com/tkeel-io/security/pkg/logger"
	"github.com/tkeel-io/security/pkg/models/rbac"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/emicklei/go-restful"
)

var _log = logger.NewLogger("auth.models.rbac")

type rbacHandler struct {
	operator *casbin.SyncedEnforcer
}

func newRBACHandler() *rbacHandler {
	e, err := rbac.NewSyncedEnforcer()

	if err != nil {
		_log.Error(err)
		return nil
	}
	return &rbacHandler{operator: e}
}

func (h *rbacHandler) AddRoleInDomain(req *restful.Request, resp *restful.Response) {
	var (
		err error
		in  AddRoleInDomainIn
	)
	tenantID := req.PathParameter("tenant_id")
	userID := req.Attribute("userID").(string)
	err = req.ReadEntity(&in)
	if err != nil || len(tenantID) == 0 || len(in.Role) == 0 {
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
		return
	}
	if h.operator == nil {
		_log.Error("nil")
	}
	ok, err := h.operator.AddRoleForUserInDomain(userID, in.Role, tenantID)
	if ok || err == nil {
		response.SrvErrWithRest(resp, errcode.SuccessServe, ok)
		return
	}
	_log.Error(err)
	response.SrvErrWithRest(resp, errcode.ErrInternalServer, nil)
}

func (h *rbacHandler) DeleteRoleInDomain(req *restful.Request, resp *restful.Response) {
	tenantID := req.PathParameter("tenant_id")
	role := req.PathParameter("role")
	users := h.operator.GetUsersForRoleInDomain(role, tenantID)
	for _, user := range users {
		_, err := h.operator.DeleteRoleForUserInDomain(user, role, tenantID)
		if err != nil {
			_log.Error(err)
		}
	}

	response.SrvErrWithRest(resp, errcode.SuccessServe, nil)
}

func (h *rbacHandler) RolesInDomain(req *restful.Request, resp *restful.Response) {
	tenantID := req.PathParameter("tenant_id")
	userID := req.Attribute("userID").(string)
	_log.Info(h.operator.GetAllRoles())
	roles := h.operator.GetRolesForUserInDomain(userID, tenantID)
	response.SrvErrWithRest(resp, errcode.SuccessServe, roles)
}

func (h *rbacHandler) AddPermissionInRole(req *restful.Request, resp *restful.Response) {
	var (
		err error
		in  AddPermissionIn
		ok  bool
	)
	tenantID := req.PathParameter("tenant_id")
	role := req.PathParameter("role")
	err = req.ReadEntity(&in)
	if err != nil || len(in.PermissionAction) == 0 || len(in.PermissionObject) == 0 {
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
	}
	ok, err = h.operator.AddPolicy(role, tenantID, in.PermissionObject, in.PermissionAction)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInternalServer, nil)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, ok)
}

func (h *rbacHandler) PermissionsInUser(req *restful.Request, resp *restful.Response) {
	userID := req.PathParameter("user_id")
	tenantID := req.PathParameter("tenant_id")
	permissions := h.operator.GetPermissionsForUserInDomain(userID, tenantID)
	response.SrvErrWithRest(resp, errcode.SuccessServe, permissions)
}

func (h *rbacHandler) DeletePermissionInRole(req *restful.Request, resp *restful.Response) {

}

func (h *rbacHandler) AddRoleToUser(req *restful.Request, resp *restful.Response) {

}

func (h *rbacHandler) DeleteRoleOnUser(req *restful.Request, resp *restful.Response) {

}

func (h *rbacHandler) PermissionCheck(req *restful.Request, resp *restful.Response) {
	var (
		err error
		in  PermissionCheck
		out *PermissionCheckResult
		ok  bool
	)
	err = req.ReadEntity(&in)
	if err != nil || len(in.PermissionObject) == 0 || len(in.UserID) == 0 {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
		return
	}
	splits := strings.Split(in.UserID, "-")
	ok, err = h.operator.Enforce(in.UserID, splits[1], in.PermissionObject, in.PermissionAction)
	out = &PermissionCheckResult{Allowed: ok}
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, out)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, out)
}
