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
	"fmt"
	"github.com/tkeel-io/security/pkg/models/rbac"
	"strconv"

	"github.com/tkeel-io/security/pkg/apiserver/response"
	"github.com/tkeel-io/security/pkg/errcode"
	"github.com/tkeel-io/security/pkg/logger"
	"github.com/tkeel-io/security/pkg/models/dao"

	"github.com/emicklei/go-restful"
)

var _log = logger.NewLogger("apirouter.tenant")

type tenantHandler struct{}

func newTenantHandler() *tenantHandler {
	return &tenantHandler{}
}

func (h *tenantHandler) Create(req *restful.Request, resp *restful.Response) {
	var (
		err error
		in  TenantCreteIn
		out *TenantCreateOut
	)
	if err = req.ReadEntity(&in); err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
		return
	}
	tenant := &dao.Tenant{
		Title:  in.Title,
		Remark: in.Remark,
	}
	if tenant.Existed() {
		response.SrvErrWithRest(resp, errcode.ErrInExistedResource, nil)
		return
	}
	if err = tenant.Create(); err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInternalServer, nil)
		return
	}
	out = &TenantCreateOut{
		ID:     tenant.ID,
		Title:  tenant.Title,
		Remark: tenant.Remark,
	}

	if in.Admin != nil {
		user := dao.User{
			ID:       dao.GenUserID(tenant.ID),
			UserName: in.Admin.UserName,
			Password: in.Admin.Password,
			Avatar:   in.Admin.Avatar,
			Email:    in.Admin.Email,
			NickName: in.Admin.NickName,
		}
		err = user.Create()
		if err == nil {
			out.Admin.UserName = user.UserName
			gPolicy := &rbac.GroupingPolicy{
				Subject: user.ID,
				Role:    "admin",
				Domain:  fmt.Sprintf("%d", tenant.ID),
			}
			policy := &rbac.Policy{
				Role:   "admin",
				Domain: fmt.Sprintf("%d", tenant.ID),
				Object: "*",
				Action: "*",
			}
			rbac.AddGroupingPolicy(gPolicy)
			rbac.AddPolicy(policy)
		}
	}

	response.SrvErrWithRest(resp, errcode.SuccessServe, out)
}

func (h *tenantHandler) Delete(req *restful.Request, resp *restful.Response) {
	var (
		err    error
		tenant *dao.Tenant
	)
	tenantID, err := strconv.Atoi(req.PathParameter("tenant_id"))
	if err != nil || tenantID == 0 {
		response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
		return
	}
	tenant = &dao.Tenant{
		ID: tenantID,
	}
	err = tenant.Delete()
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInDatabase, nil)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, nil)
}

func (h *tenantHandler) Query(req *restful.Request, resp *restful.Response) {
	var (
		err    error
		tenant *dao.Tenant
		out    []*dao.Tenant
	)
	tenantID, _ := strconv.Atoi(req.QueryParameter("tenant_id"))
	title := req.QueryParameter("title")
	tenant = &dao.Tenant{ID: tenantID, Title: title}
	out, err = tenant.List(nil)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInDatabase, nil)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, out)
}

func (h *tenantHandler) UserCreate(req *restful.Request, resp *restful.Response) {
	var (
		err error
		in  UserCreateIn
	)
	if err = req.ReadEntity(&in); err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
		return
	}
	condition := make(map[string]interface{})
	condition["tenant_id"] = in.TenantID
	condition["username"] = in.UserName
	user := &dao.User{}
	users, err := user.QueryByCondition(condition)
	_log.Info(users, err)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInDatabase, nil)
		return
	}
	if len(users) != 0 {
		response.SrvErrWithRest(resp, errcode.ErrInExistedResource, nil)
		return
	}
	user.ID = dao.GenUserID(in.TenantID)
	user.Email = in.Email
	user.Avatar = in.Avatar
	user.NickName = in.NickName
	user.UserName = in.UserName
	user.TenantID = in.TenantID
	user.Password = in.Password
	existed, err := user.Existed()
	if err != nil || existed {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInDatabase, nil)
		return
	}
	err = user.Create()
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInDatabase, nil)
		return
	}

	response.SrvErrWithRest(resp, errcode.SuccessServe, user)
}

func (h *tenantHandler) UserQuery(req *restful.Request, resp *restful.Response) {
	var (
		err  error
		user *dao.User

		out []*dao.User
	)
	tenantID, _ := strconv.Atoi(req.QueryParameter("tenant_id"))
	userID := req.QueryParameter("user_id")
	keyWords := req.QueryParameter("key_words")
	searchKey := req.QueryParameter("searchKey")
	condition := make(map[string]interface{})
	if tenantID != 0 {
		condition["tenant_id"] = tenantID
	}
	if len(userID) != 0 {
		condition["user_id"] = userID
	}
	if len(searchKey) != 0 && len(keyWords) != 0 {
		condition[searchKey] = keyWords
	}
	out, err = user.QueryByCondition(condition)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInDatabase, nil)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, out)
}
func (h *tenantHandler) UserDelete(req *restful.Request, resp *restful.Response) {
	var (
		err  error
		user *dao.User
	)
	userID := req.PathParameter("user_id")
	if len(userID) == 0 {
		response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
		return
	}
	user = &dao.User{
		ID: userID,
	}
	err = user.Delete()
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInDatabase, nil)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, nil)
}
