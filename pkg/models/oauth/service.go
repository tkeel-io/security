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
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"

	"github.com/tkeel-io/security/pkg/apiserver/response"
	"github.com/tkeel-io/security/pkg/errcode"
	"github.com/tkeel-io/security/pkg/logger"
	"github.com/tkeel-io/security/pkg/models/dao"

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
)

var (
	_log           = logger.NewLogger("auth.models.oauth")
	_oauthOperator *server.Server
)

func GetOauthOperator() *server.Server {
	return _oauthOperator
}
func OauthOperatorSetup() error {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	// token store
	redisStore := oredis.NewRedisStore(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   15,
	})
	manager.MapTokenStorage(redisStore)
	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))

	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)

	_oauthOperator = server.NewServer(server.NewConfig(), manager)
	_oauthOperator.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		splits := strings.Split(username, "-")
		tenantID, err := strconv.Atoi(splits[0])
		if err != nil {
			_log.Error(err)
			return "", err
		}
		user := &dao.User{}
		conditions := make(map[string]interface{})
		conditions["username"] = splits[1]
		conditions["tenant_id"] = tenantID
		users, err := user.QueryByCondition(conditions)
		if len(users) != 1 || err != nil {
			return "", err
		}
		err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password))
		if err != nil {
			_log.Error(err)
			return "", err
		}
		userID = users[0].ID
		return
	})
	_oauthOperator.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		return
	})
	_oauthOperator.SetAuthorizeScopeHandler(func(w http.ResponseWriter, r *http.Request) (scope string, err error) {
		return
	})
	_oauthOperator.SetResponseTokenHandler(func(w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) error {
		response.SrvErrWithWriter(w, errcode.SuccessServe, data)
		return nil
	})

	_oauthOperator.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		return &errors.Response{
			Error:       err,
			ErrorCode:   errcode.ErrInternalServer.Code(),
			Description: errcode.ErrInternalServer.Msg(),
			StatusCode:  http.StatusInternalServerError,
		}
	})

	_oauthOperator.SetResponseErrorHandler(func(re *errors.Response) {
		_log.Error("response error:", re.Error.Error())
	})
	// AuthorizationCode   GrantType = "authorization_code"
	// PasswordCredentials GrantType = "password"
	// ClientCredentials   GrantType = "client_credentials"
	// Refreshing          GrantType = "refresh_token"
	// Implicit
	_oauthOperator.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Implicit, oauth2.PasswordCredentials, oauth2.Refreshing, oauth2.ClientCredentials)
	// _oauthOperator.SetUserAuthorizationHandler(userAuthorizeHandler)
	_oauthOperator.SetAllowGetAccessRequest(true)

	_oauthOperator.SetClientAuthorizedHandler(func(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
		_log.Error(clientID, grant)
		if grant != oauth2.PasswordCredentials {

			return false, nil
		}
		return true, nil
	})

	_oauthOperator.SetClientInfoHandler(func(r *http.Request) (clientID, clientSecret string, err error) {
		return "000000", "999999", nil
	})

	_oauthOperator.SetClientScopeHandler(func(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
		return true, nil
	})

	_oauthOperator.SetRefreshingScopeHandler(func(tgr *oauth2.TokenGenerateRequest, oldScope string) (allowed bool, err error) {
		return false, nil
	})

	return nil
}
