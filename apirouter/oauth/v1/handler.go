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
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tkeel-io/security/apiserver/config"
	"github.com/tkeel-io/security/apiserver/response"
	"github.com/tkeel-io/security/constants"
	"github.com/tkeel-io/security/errcode"
	"github.com/tkeel-io/security/logger"
	"github.com/tkeel-io/security/models/dao"
	oidc2 "github.com/tkeel-io/security/models/idprovider/oidc"
	"github.com/tkeel-io/security/utils"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/emicklei/go-restful"
	oauth2V4 "github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

var _log = logger.NewLogger("auth.apirouter.oauthV1")

type oauthHandler struct {
	operator *server.Server
	conf     *config.OAuth2Config
}

func newOauthHandler(srv *server.Server, conf *config.OAuth2Config) *oauthHandler {
	return &oauthHandler{
		operator: srv,
		conf:     conf,
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
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}
}
func (h *oauthHandler) Login(req *restful.Request, resp *restful.Response) {
	err := h.operator.HandleAuthorizeRequest(resp, req.Request)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
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

//
func (h *oauthHandler) Connections(req *restful.Request, resp *restful.Response) {
	idpType := req.PathParameter("idp_type")
	tenantID := req.PathParameter("tenant_id")
	switch idpType {
	case constants.ProviderOIDC:
		state, err := utils.RandString(16)
		if err != nil {
			_log.Error(err)
			response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
			return
		}
		nonce, err := utils.RandString(16)
		if err != nil {
			_log.Error(err)
			response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
			return
		}

		oidcProvider := oidc2.GetOIDCProvider(tenantID)
		if oidcProvider == nil {
			_log.Errorf("tenant %s not register oidc provider", tenantID)
			response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
			return
		}
		setCookie(resp, req, "state", state)
		setCookie(resp, req, "nonce", nonce)
		http.Redirect(resp.ResponseWriter, req.Request, oidcProvider.OAuth2Config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)

	case constants.ProviderLDAP:
	default:
		_log.Errorf("identity type %s not supported", idpType)
		response.SrvErrWithRest(resp, errcode.ErrInIdentityProviderType, nil)
	}
}

func (h *oauthHandler) OnConnections(req *restful.Request, resp *restful.Response) {
	idpType := req.PathParameter("idp_type")
	tenantID := req.PathParameter("tenant_id")
	switch idpType {
	case constants.ProviderOIDC:
		h.onConnOIDC(req, resp, tenantID)
	default:
		_log.Errorf("unknown identity provider type: %s", idpType)
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
	}
}

func (h *oauthHandler) Register(req *restful.Request, resp *restful.Response) {
	idpType := req.PathParameter("idp_type")
	tenantID := req.PathParameter("tenant_id")
	switch idpType {
	case constants.ProviderOIDC:
		in := oidc2.Provider{}
		err := req.ReadEntity(&in)
		if err != nil {
			_log.Errorf("unknown identity provider type: %s", idpType)
			response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
			return
		}
		err = oidc2.RegisterOIDCProvider(tenantID, in)
		if err != nil {
			_log.Errorf("register oidc provider ", err)
			response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
			return
		}
		response.SrvErrWithRest(resp, errcode.SuccessServe, nil)
	default:
		_log.Errorf("unknown identity provider type: %s", idpType)
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
	}
}

func (h *oauthHandler) onConnOIDC(req *restful.Request, resp *restful.Response, key string) {
	provider := oidc2.GetOIDCProvider(key)
	state, err := req.Request.Cookie("state")
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
		return
	}
	if req.QueryParameter("state") != state.Value {
		_log.Error("state not match")
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
		return
	}

	claims, err := ClaimsOnCode(req.Request.Context(), provider, req.QueryParameter("code"))
	if err != nil {
		_log.Errorf("claims on code ", err)
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		_log.Errorf("missing required claim sub")
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}
	email, _ := claims["email"].(string)
	preferredUsername, _ := claims["preferred_username"].(string)
	if preferredUsername == "" {
		preferredUsername = subject
	}

	user, err := dao.MappingFromExternal(subject, preferredUsername, email, key)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInDatabase, nil)
		return
	}

	tgr := &oauth2V4.TokenGenerateRequest{
		ClientID:     constants.OauthClient,
		ClientSecret: constants.OauthClientSecurity,
		Request:      req.Request,
		UserID:       user.ID,
	}

	ti, err := h.operator.GetAccessToken(req.Request.Context(), oauth2V4.PasswordCredentials, tgr)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}
	tokenData := h.operator.GetTokenData(ti)
	respData := struct {
		User       *dao.User              `json:"user"`
		UserAccess map[string]interface{} `json:"user_access"`
	}{user, tokenData}

	response.SrvErrWithRest(resp, errcode.SuccessServe, respData)
}

func ClaimsOnCode(ctx context.Context, provider *oidc2.Provider, code string) (claims jwt.MapClaims, err error) {
	oauth2Token, err := provider.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		err = fmt.Errorf("failed to exchange token: %w", err)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		err = errors.New("no idtoken field in oauth2 token")
		return
	}
	idToken, err := provider.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		err = fmt.Errorf("failed to verify id token %w", err)
		return
	}

	if err = idToken.Claims(&claims); err != nil {
		err = fmt.Errorf("failed to decode id token claims %w ", err)
		return
	}
	userInfo, err := provider.Provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		err = fmt.Errorf("failed to get userinfo: %w ", err)
		return
	}
	if err = userInfo.Claims(&claims); err != nil {
		err = fmt.Errorf("failed to decode userinfo claims %w ", err)
		return
	}
	oauth2Token.AccessToken = "*REDACTED*"
	return claims, err
}

func setCookie(resp *restful.Response, r *restful.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.Request.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(resp, c)
}
