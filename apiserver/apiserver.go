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

package apiserver

import (
	"context"
	"net/http"

	entityrouter "github.com/tkeel-io/security/apirouter/entity/v1"
	oauthV1 "github.com/tkeel-io/security/apirouter/oauth/v1"
	rbacrouter "github.com/tkeel-io/security/apirouter/rbac/v1"
	tenantrouter "github.com/tkeel-io/security/apirouter/tenant/v1"
	"github.com/tkeel-io/security/apiserver/config"
	"github.com/tkeel-io/security/apiserver/filters"
	"github.com/tkeel-io/security/logger"
	"github.com/tkeel-io/security/models/dao"
	"github.com/tkeel-io/security/models/entity"
	"github.com/tkeel-io/security/tools/swagger"

	"github.com/emicklei/go-restful"
)

var (
	_log = logger.NewLogger("auth.apiserver")
)

type APIServer struct {
	Config        *config.Config
	Server        *http.Server
	restContainer *restful.Container
}

func (s *APIServer) Run(ctx context.Context) (err error) {
	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		_ = s.Server.Shutdown(shutdownCtx)
		_log.Warn("shutdown server...")
	}()

	_log.Infof("start listening on %s", s.Server.Addr)
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}

	return
}

func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	s.restContainer = restful.NewContainer()
	s.restContainer.Router(restful.CurlyRouter{})
	dao.SetUp(s.Config.Database.Mysql)
	s.installApis()
	_log.Infof(s.Config.Oauth2.AccessGenerate.AccessTokenExp)
	for _, webservice := range s.restContainer.RegisteredWebServices() {
		_log.Infof("registered root patch: %s", webservice.RootPath())
	}
	s.Server.Handler = s.restContainer

	//
	swagger.GenerateSwaggerJSON(s.restContainer, swagger.GenAuthSwaggerObjectFunc(), "./_example/auth/api/auth-openapi-spec/swagger.json")
	return nil
}

func (s *APIServer) installApis() {
	daprC, err := entity.NewGPRCClient(5, "10s", s.Config.DaprClient.GRPCPort)
	must(err)
	entityTokenOperator := entity.NewEntityTokenOperator(s.Config.DaprClient.StoreName, daprC)
	s.restContainer.Filter(filters.GlobalLog())
	must(oauthV1.RegisterToRestContainer(s.restContainer, s.Config.Oauth2))
	must(rbacrouter.RegisterToRestContainer(s.restContainer, s.Config.RBAC, s.Config.Oauth2))
	must(tenantrouter.RegisterToRestContainer(s.restContainer))
	must(entityrouter.RegisterToRestContainer(s.restContainer, s.Config.Entity, entityTokenOperator))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
