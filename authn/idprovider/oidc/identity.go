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

type oidcIdentity struct {
	// TenantID tenant id.
	TenantID string `json:"tenant_id"`
	// Subject - Identifier for the End-User at the Issuer.
	Sub string `json:"sub"`
	// Shorthand name by which the End-User wishes to be referred to at the RP,
	// such as janedoe or j.doe. This value MAY be any valid JSON string including special characters such as @, /, or whitespace.
	// The RP MUST NOT rely upon this value being unique.
	PreferredUsername string `json:"preferred_username"`
	// End-User's preferred e-mail address.
	// Its value MUST conform to the RFC 5322 [RFC5322] addr-spec syntax.
	// The RP MUST NOT rely upon this value being unique.
	Email string `json:"email"`
}

func (o oidcIdentity) GetTenantID() string {
	return ""
}

func (o oidcIdentity) GetExternalID() string {
	return o.Sub
}

func (o oidcIdentity) GetExtra() map[string]interface{} {
	return nil
}

func (o oidcIdentity) GetUserID() string {
	return o.Sub
}

func (o oidcIdentity) GetUsername() string {
	return o.PreferredUsername
}

func (o oidcIdentity) GetEmail() string {
	return o.Email
}
