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
package filters

import (
	"strings"

	"github.com/tkeel-io/security/pkg/apiserver/config"
	"github.com/tkeel-io/security/pkg/apiserver/response"
	"github.com/tkeel-io/security/pkg/errcode"
	"github.com/tkeel-io/security/pkg/models/oauth"
	"github.com/tkeel-io/security/pkg/models/rbac"

	"github.com/emicklei/go-restful"
)

func AuthFilter(conf *config.OAuth2Config, roles ...string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var err error
		if conf.AuthType == "demo" {
			req.SetAttribute("userID", "usr-1-demoUser")
			chain.ProcessFilter(req, resp)
			return
		}
		operator := oauth.GetOauthOperator()
		if operator == nil {
			operator, err = oauth.NewOperator(conf)
			if err != nil {
				_log.Error(err)
				response.SrvErrWithRest(resp, errcode.ErrServiceUnavailable, nil)
				return
			}
		}
		token, err := operator.ValidationBearerToken(req.Request)
		if err != nil {
			_log.Error(err)
			response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
			return
		}
		domain := strings.Split(token.GetUserID(), "-")[1]
		req.SetAttribute("userID", token.GetUserID())
		req.SetAttribute("tenantID", domain)

		if len(roles) == 0 {
			chain.ProcessFilter(req, resp)
			return
		}
		hasRole := false
		for i := range roles {
			_log.Info(token.GetUserID(), roles[i], domain)
			if rbac.HasRoleInDomain(token.GetUserID(), roles[i], domain) {
				hasRole = true
				break
			}
		}
		_log.Info(hasRole)
		if !hasRole {
			response.SrvErrWithRest(resp, errcode.ErrInForbiddenAccess, nil)
			return
		}
		chain.ProcessFilter(req, resp)
	}
}
