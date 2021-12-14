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
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	oidcProviders = make(map[string]*Provider)
)

type Provider struct {
	// Defines how Clients dynamically discover information about OpenID Providers
	// See also, https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfig.
	Issuer string `json:"issuer" yaml:"issuer"`

	// ClientID is the application's ID.
	ClientID string `json:"client_id" yaml:"clientID"`

	// ClientSecret is the application's secret.
	ClientSecret string `json:"client_secret" yaml:"clientSecret"`

	// Endpoint contains the resource server's token endpoint URLs.
	// These are constants specific to each server and are often available via site-specific packages,
	// such as google.Endpoint or github.Endpoint.
	Endpoint endpoint `json:"endpoint" yaml:"endpoint"`

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string `json:"redirect_url" yaml:"redirectURL"`

	// Scope specifies optional requested permissions.
	Scopes []string `json:"scopes" yaml:"scopes"`

	Provider     *oidc.Provider        `json:"-" yaml:"-"`
	OAuth2Config *oauth2.Config        `json:"-" yaml:"-"`
	Verifier     *oidc.IDTokenVerifier `json:"-" yaml:"-"`
}

// endpoint represents an OAuth 2.0 provider's authorization and token
// endpoint URLs.
type endpoint struct {
	// URL of the OP's OAuth 2.0 Authorization Endpoint [OpenID.Core](https://openid.net/specs/openid-connect-discovery-1_0.html#OpenID.Core).
	AuthURL string `json:"auth_url" yaml:"authURL"`
	// URL of the OP's OAuth 2.0 Token Endpoint [OpenID.Core](https://openid.net/specs/openid-connect-discovery-1_0.html#OpenID.Core).
	// This is REQUIRED unless only the Implicit Flow is used.
	TokenURL string `json:"token_url" yaml:"tokenURL"`
	// URL of the OP's UserInfo Endpoint [OpenID.Core](https://openid.net/specs/openid-connect-discovery-1_0.html#OpenID.Core).
	// This URL MUST use the https scheme and MAY contain port, path, and query parameter components.
	UserInfoURL string `json:"userinfo_url" yaml:"userInfoURL"`
	//  URL of the OP's JSON Web Key Set [JWK](https://openid.net/specs/openid-connect-discovery-1_0.html#JWK) document.
	JWKSURL string `json:"jwks_url"`
	// URL at the OP to which an RP can perform a redirect to request that the End-User be logged out at the OP.
	// This URL MUST use the https scheme and MAY contain port, path, and query parameter components.
	// https://openid.net/specs/openid-connect-rpinitiated-1_0.html#OPMetadata
	EndSessionURL string `json:"end_session_url"`
}

func GetOIDCProvider(key string) *Provider {
	if provider, ok := oidcProviders[key]; ok {
		return provider
	}
	return nil
}

func RegisterOIDCProvider(key string, providerConf Provider) error {
	ctx := context.TODO()

	provider, err := oidc.NewProvider(ctx, providerConf.Issuer)
	if err != nil {
		return fmt.Errorf("oidc new provider %w", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     providerConf.ClientID,
		ClientSecret: providerConf.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  providerConf.RedirectURL,
		Scopes:       providerConf.Scopes,
	}

	oidcConfig := &oidc.Config{
		ClientID: providerConf.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)
	providerConf.Provider = provider
	providerConf.OAuth2Config = oauth2Config
	providerConf.Verifier = verifier
	oidcProviders[key] = &providerConf
	return nil
}
