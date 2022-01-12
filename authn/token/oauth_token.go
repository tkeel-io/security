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

package token

import (
	"net/http"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/store"
)

type TokenConf struct {
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration

	TokenType                   string                // token type.
	AllowGetAccessRequest       bool                  // to allow GET requests for the token.
	AllowedResponseTypes        []oauth2.ResponseType // allow the authorization type.
	AllowedGrantTypes           []oauth2.GrantType    // allow the grant type.
	AllowedCodeChallengeMethods []oauth2.CodeChallengeMethod
	ForcePKCE                   bool
}

type (
	// ClientInfoHandler get client info from request
	ClientInfoHandler func(r *http.Request) (clientID, clientSecret string, err error)

	// ClientAuthorizedHandler check the client allows to use this authorization grant type
	ClientAuthorizedHandler func(clientID string, grant oauth2.GrantType) (allowed bool, err error)

	// ClientScopeHandler check the client allows to use scope
	ClientScopeHandler func(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error)

	// UserAuthorizationHandler get user id from request authorization
	UserAuthorizationHandler func(w http.ResponseWriter, r *http.Request) (userID string, err error)

	// PasswordAuthorizationHandler get user id from username and password
	PasswordAuthorizationHandler func(username, password string) (userID string, err error)

	// RefreshingScopeHandler check the scope of the refreshing token
	RefreshingScopeHandler func(tgr *oauth2.TokenGenerateRequest, oldScope string) (allowed bool, err error)

	// RefreshingValidationHandler check if refresh_token is still valid. eg no revocation or other
	RefreshingValidationHandler func(ti oauth2.TokenInfo) (allowed bool, err error)

	// ResponseErrorHandler response error handing
	ResponseErrorHandler func(re *Response)

	// InternalErrorHandler internal error handing
	InternalErrorHandler func(err error) (re *Response)

	// AuthorizeScopeHandler set the authorized scope
	AuthorizeScopeHandler func(w http.ResponseWriter, r *http.Request) (scope string, err error)

	// AccessTokenExpHandler set expiration date for the access token
	AccessTokenExpHandler func(w http.ResponseWriter, r *http.Request) (exp time.Duration, err error)

	// ExtensionFieldsHandler in response to the access token with the extension of the field
	ExtensionFieldsHandler func(ti oauth2.TokenInfo) (fieldsValue map[string]interface{})

	// ResponseTokenHandler response token handing
	ResponseTokenHandler func(w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) error
)

type OauthTokenServer struct {
	Config                       TokenConf
	Manager                      *manage.Manager
	ClientInfoHandler            ClientInfoHandler
	ClientAuthorizedHandler      ClientAuthorizedHandler
	ClientScopeHandler           ClientScopeHandler
	UserAuthorizationHandler     UserAuthorizationHandler
	PasswordAuthorizationHandler PasswordAuthorizationHandler
	RefreshingValidationHandler  RefreshingValidationHandler
	RefreshingScopeHandler       RefreshingScopeHandler
	ResponseErrorHandler         ResponseErrorHandler
	InternalErrorHandler         InternalErrorHandler
	ExtensionFieldsHandler       ExtensionFieldsHandler
	AccessTokenExpHandler        AccessTokenExpHandler
	AuthorizeScopeHandler        AuthorizeScopeHandler
	ResponseTokenHandler         ResponseTokenHandler
}

func NewOauthTokenServer(conf *TokenConf, tokenStore oauth2.TokenStore, generator oauth2.AccessGenerate, client oauth2.ClientInfo) *OauthTokenServer {
	manager := manage.NewDefaultManager()
	tokenConf := manage.DefaultAuthorizeCodeTokenCfg
	if conf != nil {
		tokenConf.AccessTokenExp = conf.AccessTokenExp
		tokenConf.RefreshTokenExp = conf.RefreshTokenExp
	}

	clientStore := store.NewClientStore()
	if client == nil {
		client = &models.Client{ID: DefaultClient, Secret: DefaultClientSecurity, Domain: DefaultClientDomain}
	}
	clientStore.Set(client.GetID(), client)
	if tokenStore == nil {
		tokenStore, _ = store.NewMemoryTokenStore()
	}
	if generator == nil {
		generator = generates.NewAccessGenerate()
	}

	manager.SetPasswordTokenCfg(tokenConf)
	manager.MapClientStorage(clientStore)
	manager.MapTokenStorage(tokenStore)
	manager.MapAccessGenerate(generator)

	tokenServer := &OauthTokenServer{Manager: manager}
	return tokenServer
}
