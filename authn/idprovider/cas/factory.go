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

package cas

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/tkeel-io/security/authn/idprovider"

	"github.com/mitchellh/mapstructure"
	gocas "gopkg.in/cas.v2"
)

func init() {
	idprovider.RegisterProviderFactory(&casProviderFactory{})
}

type casProviderFactory struct {
}

func (f casProviderFactory) Type() string {
	return "CASIdentityProvider"
}

//nolint
func (f casProviderFactory) Create(options map[string]interface{}) (idprovider.Provider, error) {
	var cas casProvider
	if err := mapstructure.Decode(options, &cas); err != nil {
		return nil, err
	}
	casURL, err := url.Parse(cas.CASServerURL)
	if err != nil {
		return nil, err
	}
	redirectURL, err := url.Parse(cas.RedirectURL)
	if err != nil {
		return nil, err
	}
	cas.client = gocas.NewRestClient(&gocas.RestOptions{
		CasURL:     casURL,
		ServiceURL: redirectURL,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: cas.InsecureSkipVerify},
			},
		},
		URLScheme: nil,
	})
	return &cas, nil
}
