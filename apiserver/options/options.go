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

package options

import (
	"net/http"

	"github.com/tkeel-io/security/apiserver"
	"github.com/tkeel-io/security/apiserver/config"
)

type ServerRunOptions struct {
	ConfigFile string
	*config.Config
}

func NewServerRunOptions() *ServerRunOptions {
	conf, _ := config.TryLoadFromDisk()
	opts := &ServerRunOptions{
		Config: conf,
	}
	return opts
}

func (opts *ServerRunOptions) NewAPIServer(stopCh <-chan struct{}) (*apiserver.APIServer, error) {
	apiServer := &apiserver.APIServer{
		Config: opts.Config,
	}
	server := &http.Server{
		Addr: opts.Config.Server.Address,
	}
	apiServer.Server = server
	return apiServer, nil
}
