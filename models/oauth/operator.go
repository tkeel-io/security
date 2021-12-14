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

package oauth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tkeel-io/security/apiserver/config"
	"github.com/tkeel-io/security/apiserver/response"
	"github.com/tkeel-io/security/constants"
	"github.com/tkeel-io/security/errcode"
	"github.com/tkeel-io/security/logger"
	"github.com/tkeel-io/security/models/dao"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	_log           = logger.NewLogger("auth.models.oauth")
	_oauthOperator *server.Server
)

func GetOauthOperator() *server.Server {
	return _oauthOperator
}
func NewOperator(conf *config.OAuth2Config) (*server.Server, error) {
	var err error
	manager := manage.NewDefaultManager()
	tokenConf := &manage.Config{}
	if conf.AccessGenerate.AccessTokenExp == "" {
		tokenConf = manage.DefaultAuthorizeCodeTokenCfg
	} else {
		tokenConf.AccessTokenExp, err = time.ParseDuration(conf.AccessGenerate.AccessTokenExp)
		_log.Error(err)
		tokenConf.RefreshTokenExp = tokenConf.AccessTokenExp * 10
		tokenConf.IsGenerateRefresh = true
	}

	// manager.SetAuthorizeCodeTokenCfg(tokenConf).
	manager.SetPasswordTokenCfg(tokenConf)
	// token store.
	redisStore := oredis.NewRedisStore(&redis.Options{
		Addr:     conf.Redis.Addr,
		DB:       conf.Redis.DB,
		Password: conf.Redis.Password,
	})
	manager.MapTokenStorage(redisStore)
	// generate jwt access token.
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(conf.AccessGenerate.SecurityKey), jwt.SigningMethodHS512))

	clientStore := store.NewClientStore()
	clientStore.Set(constants.OauthClient, &models.Client{
		ID:     constants.OauthClient,
		Secret: constants.OauthClientSecurity,
		Domain: constants.OauthClientDomain,
	})
	manager.MapClientStorage(clientStore)
	_oauthOperator = server.NewServer(server.NewConfig(), manager)
	_oauthOperator.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		splits := strings.Split(username, "-")
		tenantID, err := strconv.Atoi(splits[0])
		if err != nil {
			_log.Error(err)
			return "", fmt.Errorf("atoi tenant id %w ", err)
		}
		user := &dao.User{}
		conditions := make(map[string]interface{})
		conditions["username"] = splits[1]
		conditions["tenant_id"] = tenantID
		users, err := user.QueryByCondition(conditions)
		if len(users) == 0 || err != nil {
			if err != nil {
				err = fmt.Errorf("user not found %w", err)
			} else {
				err = fmt.Errorf("user not found")
			}
			return "", err
		}
		err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password))
		if err != nil {
			_log.Error(err)
			return "", fmt.Errorf("compare hash and password %w", err)
		}
		userID = users[0].ID
		return
	})
	_oauthOperator.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		return
	})
	_oauthOperator.SetResponseTokenHandler(func(w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) error {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		for key := range header {
			w.Header().Set(key, header.Get(key))
		}

		if v, ok := data["error"]; ok {
			response.DefineWithWriter(w, http.StatusBadRequest, errcode.CodeInvalidParam, v.(string), nil)
			return fmt.Errorf(v.(string))
		}
		response.SrvErrWithWriter(w, errcode.SuccessServe, data)
		return nil
	})
	_oauthOperator.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		return &errors.Response{
			Error:       errors.New(errcode.ErrInUnexpected.Msg()),
			ErrorCode:   errcode.ErrInUnexpected.Code(),
			Description: errcode.ErrInUnexpected.Msg(),
			StatusCode:  http.StatusInternalServerError,
		}
	})

	_oauthOperator.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Implicit, oauth2.PasswordCredentials, oauth2.Refreshing, oauth2.ClientCredentials)

	_oauthOperator.SetAllowGetAccessRequest(true)

	_oauthOperator.SetClientInfoHandler(func(r *http.Request) (clientID, clientSecret string, err error) {
		return constants.OauthClient, constants.OauthClientSecurity, nil
	})

	return _oauthOperator, nil
}
