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

package constants

// API Tag.
const (
	// APITagRBAC swagger tag RBAC.
	APITagRBAC = "RBAC"
	// APITagTenant swagger tag Tenant.
	APITagTenant = "Tenant"
	// APITagOauth swagger tag Oauth.
	APITagOauth = "Oauth"
	// APITagEntity swagger tag Entity.
	APITagEntity = "Entity"
)

// IdentityProvider Type.
const (
	// ProviderOIDC  OIDC identity provider type.
	ProviderOIDC = "OIDC"
	// ProviderOIDC LDAP identity provider type.
	ProviderLDAP = "LDAP"
)

const (
	// OauthClient clientID.
	OauthClient = "000000"
	// OauthClientSecurity .
	OauthClientSecurity = "999999"
	// OauthClientDomain .
	OauthClientDomain = "http://localhost"
)
