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

	"github.com/tkeel-io/security/apiserver/config"
	"github.com/tkeel-io/security/apiserver/response"
	"github.com/tkeel-io/security/errcode"
	"github.com/tkeel-io/security/logger"
	"github.com/tkeel-io/security/models/entity"
	"github.com/tkeel-io/security/models/oauth"

	"github.com/emicklei/go-restful"
	"github.com/golang-jwt/jwt"
)

var (
	_log = logger.NewLogger("auth.apirouter.entityV1")
)

type entityHandler struct {
	tokenOperator       *oauth.JWTAccessGenerate
	entityTokenOperator entity.TokenOperator
}

func newEntityHandler(conf *config.EntityConfig, entityOperator entity.TokenOperator) *entityHandler {
	operator := oauth.NewJWTAccessGenerate("", []byte(conf.SecurityKey), jwt.SigningMethodHS512)
	return &entityHandler{operator, entityOperator}
}

func (h *entityHandler) Token(req *restful.Request, resp *restful.Response) {
	var (
		err             error
		expiresDuration time.Duration
		out             *Token
	)
	entityType := req.PathParameter("entity_type")
	entityID := req.PathParameter("entity_id")
	owner := req.QueryParameter("owner")
	if len(entityType) == 0 || len(entityID) == 0 || len(owner) == 0 {
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
		return
	}
	expiresIn, _ := strconv.Atoi(req.QueryParameter("expires_in"))
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
		response.SrvErrWithRest(resp, errcode.ErrInGenToken, nil)
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

func (h *entityHandler) CreateEntityToken(req *restful.Request, resp *restful.Response) {
	in := &EntityTokenIn{}
	err := req.ReadEntity(in)
	if err != nil || in.EntityID == "" || in.EntityType == "" {
		_log.Error("create entity token invalid params: ", in)
		response.SrvErrWithRest(resp, errcode.ErrInvalidParam, nil)
		return
	}
	now := time.Now()
	if in.ExpiresIn == 0 {
		in.ExpiresIn = now.Add(time.Hour * 24 * 365).Unix()
	} else {
		in.ExpiresIn = now.Add(time.Second * time.Duration(in.ExpiresIn)).Unix()
	}
	entityInfo := &entity.Token{
		EntityID:   in.EntityID,
		EntityType: in.EntityType,
		Owner:      in.Owner,
		CreatedAt:  now.Unix(),
		ExpiredAt:  in.ExpiresIn,
	}
	token, err := h.entityTokenOperator.CreateToken(req.Request.Context(), entityInfo)
	if err != nil {
		_log.Error(err)
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, Token{token})
}

func (h *entityHandler) GetEntityInfo(req *restful.Request, resp *restful.Response) {
	token := req.PathParameter("token")
	entityInfo, err := h.entityTokenOperator.GetEntityInfo(req.Request.Context(), token)
	if err != nil {
		_log.Error(err, token)
		response.SrvErrWithRest(resp, errcode.ErrInUnexpected, nil)
		return
	}
	response.SrvErrWithRest(resp, errcode.SuccessServe, entityInfo)
}
