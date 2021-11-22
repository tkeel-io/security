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
	"time"

	"github.com/tkeel-io/security/pkg/apiserver/response"
	"github.com/tkeel-io/security/pkg/errcode"
	"github.com/tkeel-io/security/pkg/logger"

	"github.com/emicklei/go-restful"
	"github.com/go-oauth2/oauth2/v4/server"
)

var _log = logger.NewLogger("auth.apirouter.oauthV1")

type oauthHandler struct {
	operator *server.Server
}

func newOauthHandler(srv *server.Server) *oauthHandler {
	return &oauthHandler{
		operator: srv,
	}
}

func (h *oauthHandler) Token(req *restful.Request, resp *restful.Response) {
	err := h.operator.HandleTokenRequest(resp, req.Request)
	if err != nil {
		_log.Error(err)
		return
	}
}

func (h *oauthHandler) Authorize(req *restful.Request, resp *restful.Response) {
	err := h.operator.HandleAuthorizeRequest(resp, req.Request)
	if err != nil {
		_log.Error(err)
		return
	}
}
func (h *oauthHandler) Login(req *restful.Request, resp *restful.Response) {
	err := h.operator.HandleAuthorizeRequest(resp, req.Request)
	if err != nil {
		_log.Error(err)
		return
	}
}

func (h *oauthHandler) OnCode(req *restful.Request, resp *restful.Response) {
	response.SrvErrWithRest(resp, errcode.SuccessServe, req.Request.RequestURI)
}

func (h *oauthHandler) CheckAuth(req *restful.Request, resp *restful.Response) {
	token, err := h.operator.ValidationBearerToken(req.Request)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
		return
	}

	cli, err := h.operator.Manager.GetClient(req.Request.Context(), token.GetClientID())
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(time.Until(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn())).Seconds()),
		"user_id":    token.GetUserID(),
		"client_id":  token.GetClientID(),
		"scope":      token.GetScope(),
		"domain":     cli.GetDomain(),
	}

	response.SrvErrWithRest(resp, errcode.SuccessServe, data)
}

func (h *oauthHandler) Authenticate(req *restful.Request, resp *restful.Response) {
	token, err := h.operator.ValidationBearerToken(req.Request)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
		return
	}

	cli, err := h.operator.Manager.GetClient(req.Request.Context(), token.GetClientID())
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(time.Until(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn())).Seconds()),
		"user_id":    token.GetUserID(),
		"client_id":  token.GetClientID(),
		"scope":      token.GetScope(),
		"domain":     cli.GetDomain(),
	}

	response.SrvErrWithRest(resp, errcode.SuccessServe, data)
}
