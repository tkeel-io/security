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

package idprovider

type Provider interface {
	// Type unique type of the provider.
	Type() string
	// AuthenticateCode authenticate identity with code from remote server.
	AuthenticateCode(code string) (Identity, error)
	// Authenticate basic authn username password.
	Authenticate(username string, password string) (Identity, error)
	AuthCodeURL(state, nonce string) string
}

type ProviderFactory interface {
	// Type unique type of the provider.
	Type() string
	// Apply the  options from external identity provider config.
	Create(options map[string]interface{}) (Provider, error)
}
