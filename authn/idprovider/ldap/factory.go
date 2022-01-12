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

package ldap

import (
	"github.com/mitchellh/mapstructure"
	"github.com/tkeel-io/security/authn/idprovider"
)

func init() {
	idprovider.RegisterProviderFactory(&ldapProviderFactory{})
}

type ldapProviderFactory struct {
}

func (l *ldapProviderFactory) Type() string {
	return _ldapIdentityProvider
}

func (l *ldapProviderFactory) Create(options map[string]interface{}) (idprovider.Provider, error) {
	var ldapProvider ldapProvider
	if err := mapstructure.Decode(options, &ldapProvider); err != nil {
		return nil, err
	}
	if ldapProvider.ReadTimeout <= 0 {
		ldapProvider.ReadTimeout = _defaultReadTimeout
	}
	return &ldapProvider, nil
}
