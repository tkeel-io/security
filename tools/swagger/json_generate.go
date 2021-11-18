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
package swagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/tkeel-io/security/pkg/logger"
	"github.com/tkeel-io/security/pkg/version"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

var (
	_log = logger.NewLogger("auth.tools.swagger")
)

func GenerateSwaggerJSON(c *restful.Container, f restfulspec.PostBuildSwaggerObjectFunc, output string) []byte {
	config := restfulspec.Config{
		WebServices:                   c.RegisteredWebServices(),
		PostBuildSwaggerObjectHandler: f}
	swagger := restfulspec.BuildSwagger(config)
	swagger.Info.Extensions = make(spec.Extensions)
	swagger.Info.Extensions.Add("x-tagGroups", []struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}{})

	data, _ := json.MarshalIndent(swagger, "", "  ")
	if len(output) != 0 {
		err := ioutil.WriteFile(output, data, 420)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("successfully written to %s", output)
	return data
}

func GenAuthSwaggerObjectFunc() restfulspec.PostBuildSwaggerObjectFunc {
	return func(swo *spec.Swagger) {
		swo.Info = &spec.Info{
			InfoProps: spec.InfoProps{
				Title:       "Auth",
				Description: "Auth OpenAPI",
				Version:     version.Get().AppVersion,
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
}

func ValidateSpec(apiSpec []byte) error {
	swaggerDoc, err := loads.Analyzed(apiSpec, "")
	if err != nil {
		return fmt.Errorf("swagger dpc analyzed err : %w", err)
	}
	// Attempts to report about all errors
	validate.SetContinueOnErrors(true)
	v := validate.NewSpecValidator(swaggerDoc.Schema(), strfmt.Default)
	result, _ := v.Validate(swaggerDoc)
	if result.HasWarnings() {
		_log.Info("see warnings below:\n")
		for _, desc := range result.Warnings {
			_log.Warnf("- WARNING: %s\n", desc.Error())
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
