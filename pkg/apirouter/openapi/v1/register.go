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
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/tkeel-io/security/pkg/constants"
	"github.com/tkeel-io/security/pkg/errcode"
	"net/http"
)

func AddToRestContainer(c *restful.Container) error {
	webservice := &restful.WebService{}
	webservice.Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	handler := newOpenApiHandler()

	webservice.Route(webservice.GET("identify").
		To(handler.Identify).
		Doc("identify for plugin register ").
		Returns(http.StatusOK, errcode.ErrMsgOK, identifyResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagPluginRequired}))

	webservice.Route(webservice.GET("status").
		To(handler.Status).
		Doc(" status for plugin register ").
		Returns(http.StatusOK, errcode.ErrMsgOK, status{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ApiTagPluginRequired}))
	c.Add(webservice)
	return nil
}
