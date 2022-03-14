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

package oidc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/tkeel-io/security/authn/idprovider"
	"github.com/tkeel-io/security/utils"

	"github.com/coreos/go-oidc"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"
)

func init() {
	idprovider.RegisterProviderFactory(&oidcProviderFactory{})
}

type oidcProviderFactory struct {
}

func (f *oidcProviderFactory) Type() string {
	return _oidcIdentityType
}

//nolint
func (f *oidcProviderFactory) Create(options map[string]interface{}) (idprovider.Provider, error) {
	var oidcProvider OIDCProvider
	if err := mapstructure.Decode(options, &oidcProvider); err != nil {
		return nil, fmt.Errorf("mapstructure decode provider options %w", err)
	}
	if oidcProvider.Issuer != "" {
		ctx := context.TODO()
		if oidcProvider.InsecureSkipVerify {
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true, // nolint
					},
				},
			}
			ctx = oidc.ClientContext(ctx, client)
		}
		provider, err := oidc.NewProvider(ctx, oidcProvider.Issuer)
		if err != nil {
			return nil, fmt.Errorf("failed to create oidc provider: %w", err)
		}
		var providerJSON map[string]interface{}
		if err = provider.Claims(&providerJSON); err != nil {
			return nil, fmt.Errorf("failed to decode oidc provider claims: %w", err)
		}
		oidcProvider.Endpoint.AuthURL, _ = providerJSON["authorization_endpoint"].(string)
		oidcProvider.Endpoint.TokenURL, _ = providerJSON["token_endpoint"].(string)
		oidcProvider.Endpoint.UserInfoURL, _ = providerJSON["userinfo_endpoint"].(string)
		oidcProvider.Endpoint.JWKSURL, _ = providerJSON["jwks_uri"].(string)
		oidcProvider.Endpoint.EndSessionURL, _ = providerJSON["end_session_endpoint"].(string)
		oidcProvider.Provider = provider
		oidcProvider.Verifier = provider.Verifier(&oidc.Config{
			// TODO: support HS256.
			ClientID: oidcProvider.ClientID,
		})
		options["endpoint"] = map[string]interface{}{
			"auth_url":        oidcProvider.Endpoint.AuthURL,
			"token_url":       oidcProvider.Endpoint.TokenURL,
			"user_info_url":   oidcProvider.Endpoint.UserInfoURL,
			"jwksurl":         oidcProvider.Endpoint.JWKSURL,
			"end_session_url": oidcProvider.Endpoint.EndSessionURL,
		}
	}
	scopes := []string{oidc.ScopeOpenID}
	if !utils.StringsInclude(oidcProvider.Scopes, oidc.ScopeOpenID) {
		scopes = append(scopes, oidcProvider.Scopes...)
	}

	oidcProvider.Scopes = scopes
	oidcProvider.OAuth2Config = &oauth2.Config{
		ClientID:     oidcProvider.ClientID,
		ClientSecret: oidcProvider.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: oidcProvider.Endpoint.TokenURL,
			AuthURL:  oidcProvider.Endpoint.AuthURL,
		},
		RedirectURL: oidcProvider.RedirectURL,
		Scopes:      oidcProvider.Scopes,
	}
	return &oidcProvider, nil
}
