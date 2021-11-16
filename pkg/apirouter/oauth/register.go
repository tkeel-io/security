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
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/tkeel-io/security/pkg/constants"
	"github.com/tkeel-io/security/pkg/models/oauth"

	"github.com/emicklei/go-restful"
)

func AddToRestContainer(c *restful.Container) error {

	webservice := &restful.WebService{}
	webservice.Path("/oauth").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	oauth.OperatorSetup()
	oauthOperator := oauth.GetOauthOperator()
	handler := newOauthHandler(oauthOperator)

	webservice.Route(webservice.GET("authorize").
		To(handler.Authorize).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagOauth}))

	webservice.Route(webservice.GET("token").
		To(handler.Token).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagOauth}))

	webservice.Route(webservice.GET("authenticate").
		To(handler.Authenticate).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagOauth}))

	webservice.Route(webservice.GET("on_code").
		To(handler.OnCode).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagOauth}))
	c.Add(webservice)
	return nil
}
