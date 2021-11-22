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

	"github.com/tkeel-io/security/pkg/apiserver/config"
	"github.com/tkeel-io/security/pkg/apiserver/filters"
	"github.com/tkeel-io/security/pkg/constants"
	"github.com/tkeel-io/security/pkg/models/oauth"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func AddToRestContainer(c *restful.Container, conf *config.OAuth2Config) error {
	var webservice *restful.WebService
	for _, v := range c.RegisteredWebServices() {
		if v.RootPath() == "v1" {
			webservice = v
			break
		}
	}
	if webservice == nil {
		webservice = &restful.WebService{}
		webservice.Path("v1").
			Produces(restful.MIME_JSON).
			Filter(filters.Auth())

		c.Add(webservice)
	}

	oauthOperator, err := oauth.NewOperator(conf)
	if err != nil {
		_log.Error(err)
		return fmt.Errorf("oauth new operate %w", err)
	}
	handler := newOauthHandler(oauthOperator)

	webservice.Route(webservice.GET("oauth/authorize").
		To(handler.Authorize).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagOauth}))

	webservice.Route(webservice.GET("oauth/token").
		To(handler.Token).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagOauth}))

	webservice.Route(webservice.GET("oauth/authenticate").
		To(handler.Authenticate).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagOauth}))

	webservice.Route(webservice.GET("oauth/on_code").
		To(handler.OnCode).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagOauth}))

	return nil
}
