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
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	oauthrouter "github.com/tkeel-io/security/pkg/apirouter/oauth"
	openapirouter "github.com/tkeel-io/security/pkg/apirouter/openapi/v1"
	rbacrouter "github.com/tkeel-io/security/pkg/apirouter/rbac/v1"
	tenantrouter "github.com/tkeel-io/security/pkg/apirouter/tenant/v1"
	"github.com/tkeel-io/security/pkg/logger"
	"github.com/tkeel-io/security/pkg/version"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
)

var (
	output string
	_log   = logger.NewLogger("auth.api.doc")
)

func init() {
	flag.StringVar(&output, "output", "./api/auth-openapi-spec/swagger.json", "--output=./api.json")
}

func main() {
	flag.Parse()
	swaggerSpec := generateSwaggerJson()

	err := validateSpec(swaggerSpec)
	if err != nil {
		_log.Warn("Swagger specification has errors ", err)
	}
}

func validateSpec(apiSpec []byte) error {

	swaggerDoc, err := loads.Analyzed(apiSpec, "")
	if err != nil {
		return err
	}

	// Attempts to report about all errors
	validate.SetContinueOnErrors(true)

	v := validate.NewSpecValidator(swaggerDoc.Schema(), strfmt.Default)
	result, _ := v.Validate(swaggerDoc)

	if result.HasWarnings() {
		log.Printf("See warnings below:\n")
		for _, desc := range result.Warnings {
			log.Printf("- WARNING: %s\n", desc.Error())
		}

	}
	if result.HasErrors() {
		str := fmt.Sprintf("The swagger spec is invalid against swagger specification %s.\nSee errors below:\n", swaggerDoc.Version())
		for _, desc := range result.Errors {
			str += fmt.Sprintf("- %s\n", desc.Error())
		}
		log.Println(str)
		return errors.New(str)
	}

	return nil
}

func generateSwaggerJson() []byte {

	container := restful.NewContainer()

	(oauthrouter.AddToRestContainer(container))
	(openapirouter.AddToRestContainer(container))
	(rbacrouter.AddToRestContainer(container))
	(tenantrouter.AddToRestContainer(container))

	config := restfulspec.Config{
		WebServices:                   container.RegisteredWebServices(),
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}

	swagger := restfulspec.BuildSwagger(config)
	swagger.Info.Extensions = make(spec.Extensions)
	swagger.Info.Extensions.Add("x-tagGroups", []struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}{})

	data, _ := json.MarshalIndent(swagger, "", "  ")
	err := ioutil.WriteFile(output, data, 420)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully written to %s", output)

	return data
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Auth",
			Description: "Auth OpenAPI",
			Version:     version.Get().GitVersion,
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "tKeel",
					URL:   "https://tKeel.io/",
					Email: "auth@tkeel.com",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache 2.0",
					URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
				},
			},
		},
	}

	// setup security definitions
	swo.SecurityDefinitions = map[string]*spec.SecurityScheme{
		"jwt": spec.APIKeyAuth("Authorization", "header"),
	}
	swo.Security = []map[string][]string{{"jwt": []string{}}}
}
