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
	"errors"
	"fmt"

	"github.com/tkeel-io/security/authn/idprovider"

	gocas "gopkg.in/cas.v2"
)

const _casIdentityProvider = "CASIdentityProvider"

var _ idprovider.Provider = &casProvider{}

type casProvider struct {
	RedirectURL        string `json:"redirect_url" yaml:"redirectURL"`    //nolint
	CASServerURL       string `json:"cas_server_url" yaml:"casServerURL"` //nolint
	InsecureSkipVerify bool   `json:"insecure_skip_verify" yaml:"insecureSkipVerify"`
	client             *gocas.RestClient
}

func (c casProvider) AuthCodeURL(state, nonce string) string {
	//TODO implement me
	panic("implement me")
}

func (c casProvider) Type() string {
	return _casIdentityProvider
}

//nolint
func (c casProvider) AuthenticateCode(ticket string) (idprovider.Identity, error) {
	resp, err := c.client.ValidateServiceTicket(gocas.ServiceTicket(ticket))
	if err != nil {
		return nil, fmt.Errorf("cas validate service ticket failed: %w", err)
	}
	return &casIdentity{User: resp.User}, nil
}

//nolint
func (c casProvider) Authenticate(username string, password string) (idprovider.Identity, error) {
	return nil, errors.New("unsupported authenticate with username password")
}
