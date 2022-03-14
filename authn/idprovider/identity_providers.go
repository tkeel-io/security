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

import "errors"

// tenantID:provider.
var (
	// _providers all providers with tenantID:provider.
	_providers = make(map[string]Provider)
	// _providerFactories all provider factory with type:factory.
	_providerFactories = make(map[string]ProviderFactory)
	// ErrIdentityProviderNotFound error in not found identity provider.
	ErrIdentityProviderNotFound = errors.New("identity provider not found")
)

// RegisterProviderFactory  registers ProviderFactory with the specified type.
func RegisterProviderFactory(factory ProviderFactory) {
	_providerFactories[factory.Type()] = factory
}

//nolint
// GetIdentityProvider returns identity Provider with key.
func GetIdentityProvider(key string) (Provider, error) {
	if provider, ok := _providers[key]; ok {
		return provider, nil
	}
	return nil, ErrIdentityProviderNotFound
}

//nolint
// RegisterIdentityProvider register Provider with key.
func RegisterIdentityProvider(key string, provider Provider) {
	_providers[key] = provider
}
