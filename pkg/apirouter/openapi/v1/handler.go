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
	"github.com/tkeel-io/security/pkg/version"
)

type openAPIHandler struct {
}

func newOpenApiHandler() *openAPIHandler {
	return &openAPIHandler{}
}

// TODO fix
func (h *openAPIHandler) Identify(req *restful.Request, resp *restful.Response) {
	success := identifyResponse{
		Code:         0,
		Msg:          "ok",
		PluginID:     "auth",
		Version:      version.Get().AppVersion,
		TkeelVersion: "v0.1.0",
	}
	resp.WriteEntity(success)
}

// TODO fix
func (h *openAPIHandler) Status(req *restful.Request, resp *restful.Response) {
	success := status{
		Code:   0,
		Msg:    "ok",
		Status: "ACTIVE",
	}
	resp.WriteEntity(success)
}

type identifyResponse struct {
	Code         int            `json:"code"`
	Msg          string         `json:"msg"`
	PluginID     string         `json:"plugin_id"`
	Version      string         `json:"version"`
	TkeelVersion string         `json:"tkeel_version"`
	AddonsPoints []*AddonsPoint `json:"addons_points,omitempty"`
	MainPlugins  []*MainPlugin  `json:"main_plugins,omitempty"`
}

type status struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
}
type AddonsPoint struct {
	AddonsPoint string `json:"addons_point"`
	Desc        string `json:"desc,omitempty"`
}

type AddonsEndpoint struct {
	AddonsPoint string `json:"addons_point"`
	Endpoint    string `json:"endpoint"`
}

type MainPlugin struct {
	ID        string            `json:"id"`
	Version   string            `json:"version,omitempty"`
	Endpoints []*AddonsEndpoint `json:"endpoints"`
}
