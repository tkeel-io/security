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
	"strconv"
	"time"

	"github.com/tkeel-io/security/pkg/apiserver/config"
	"github.com/tkeel-io/security/pkg/apiserver/response"
	"github.com/tkeel-io/security/pkg/errcode"
	"github.com/tkeel-io/security/pkg/logger"
	"github.com/tkeel-io/security/pkg/models/oauth"

	"github.com/emicklei/go-restful"
	"github.com/golang-jwt/jwt"
)

var (
	_log = logger.NewLogger("auth.apirouter.entityV1")
)

type entityHandler struct {
	tokenOperator *oauth.JWTAccessGenerate
}

func newEntityHandler(conf *config.EntityConfig) *entityHandler {
	operator := oauth.NewJWTAccessGenerate("", []byte(conf.SecurityKey), jwt.SigningMethodHS512)
	return &entityHandler{operator}
}

func (h *entityHandler) Token(req *restful.Request, resp *restful.Response) {
	var (
		err             error
		expiresDuration time.Duration
		out             *Token
	)
	entityType := req.PathParameter("entity_type")
	entityID := req.PathParameter("entity_id")
	expiresIn, _ := strconv.Atoi(req.QueryParameter("expires_in"))
	owner := req.Attribute("userID")
	if expiresIn == 0 {
		expiresDuration = time.Hour * 24 * 365
	} else {
		expiresDuration = time.Second * time.Duration(expiresIn)
	}
	expiresAt := time.Now().Add(expiresDuration).Unix()
	claims := jwt.MapClaims{}
	claims["entity_type"] = entityType
	claims["entity_id"] = entityID
	claims["owner"] = owner
	claims["exp"] = expiresAt
	out = &Token{}
	out.Token, _, err = h.tokenOperator.Token(req.Request.Context(), &claims, false)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrGenToken, nil)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, out)
}
func (h *entityHandler) TokenValid(req *restful.Request, resp *restful.Response) {
	var (
		err error
		in  *TokenValidIn
	)
	err = req.ReadEntity(&in)
	if err != nil || in == nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
		return
	}

	claims, err := h.tokenOperator.Valid(in.EntityToken)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, claims)
}
