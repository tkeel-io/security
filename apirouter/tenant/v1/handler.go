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
	"strconv"

	"github.com/tkeel-io/security/apiserver/response"
	"github.com/tkeel-io/security/errcode"
	"github.com/tkeel-io/security/logger"
	dao2 "github.com/tkeel-io/security/models/dao"
	rbac2 "github.com/tkeel-io/security/models/rbac"

	"github.com/emicklei/go-restful"
)

var _log = logger.NewLogger("auth.apirouter.tenantV1")

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
	tenant := &dao2.Tenant{
		Title:  in.Title,
		Remark: in.Remark,
	}
	if tenant.Existed() {
		response.SrvErrWithRest(resp, errcode.ErrInExistedResource, nil)
		return
	}
	if err = tenant.Create(); err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}
	out = &TenantCreateOut{
		ID:     tenant.ID,
		Title:  tenant.Title,
		Remark: tenant.Remark,
	}

	if in.Admin != nil {
		user := dao2.User{
			ID:       dao2.GenUserID(tenant.ID),
			TenantID: tenant.ID,
			UserName: in.Admin.UserName,
			Password: in.Admin.Password,
			Avatar:   in.Admin.Avatar,
			Email:    in.Admin.Email,
			NickName: in.Admin.NickName,
		}
		err = user.Create()
		if err == nil {
			out.Admin.UserName = user.UserName
			gPolicy := &rbac2.GroupingPolicy{
				Subject: user.ID,
				Role:    "admin",
				Domain:  fmt.Sprintf("%d", tenant.ID),
			}
			policy := &rbac2.Policy{
				Role:   "admin",
				Domain: fmt.Sprintf("%d", tenant.ID),
				Object: "*",
				Action: "*",
			}
			rbac2.AddGroupingPolicy(gPolicy)
			rbac2.AddPolicy(policy)
		}
	}

	response.SrvErrWithRest(resp, errcode.SuccessServe, out)
}

func (h *tenantHandler) Delete(req *restful.Request, resp *restful.Response) {
	var (
		err    error
		tenant *dao2.Tenant
	)
	tenantID, err := strconv.Atoi(req.PathParameter("tenant_id"))
	if err != nil || tenantID == 0 {
		response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
		return
	}
	tenant = &dao2.Tenant{
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
		tenant *dao2.Tenant
		out    []*dao2.Tenant
	)
	tenantID, _ := strconv.Atoi(req.QueryParameter("tenant_id"))
	title := req.QueryParameter("title")
	tenant = &dao2.Tenant{ID: tenantID, Title: title}
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
	user := &dao2.User{}
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
	user.ID = dao2.GenUserID(in.TenantID)
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
		user *dao2.User

		out []*dao2.User
	)
	tenantID, _ := strconv.Atoi(req.QueryParameter("tenant_id"))
	userID := req.QueryParameter("user_id")
	keyWords := req.QueryParameter("key_words")
	searchKey := req.QueryParameter("search_key")
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
		user *dao2.User
	)
	userID := req.PathParameter("user_id")
	if len(userID) == 0 {
		response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
		return
	}
	user = &dao2.User{
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
