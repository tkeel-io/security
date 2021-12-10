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
	"github.com/tkeel-io/security/apirouter"
	"github.com/tkeel-io/security/apiserver/config"
	"github.com/tkeel-io/security/constants"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func RegisterToRestContainer(c *restful.Container, conf *config.EntityConfig, authConf *config.OAuth2Config) error {
	webservice := apirouter.GetWebserviceWithPatch(c, "/v1/entity")
	handler := newEntityHandler(conf)

	webservice.Route(webservice.GET("/{entity_type}/{entity_id}/token").
		To(handler.Token).
		Doc("get a entity token").
		Param(webservice.PathParameter("entity_type", "EntityType").Required(true)).
		Param(webservice.PathParameter("entity_id", "Entity's ID").Required(true)).
		Param(webservice.QueryParameter("owner", "user`s ID").Required(true)).
		Param(webservice.QueryParameter("expires_in", "invalid period( seconds );default 365 days").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagEntity}))

	webservice.Route(webservice.POST("/token/valid").
		To(handler.TokenValid).
		Doc("valid a entity token").
		Reads(TokenValidIn{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.APITagEntity}))

	return nil
}
