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

package token

import (
	"net/http"
	"time"

	"github.com/go-oauth2/oauth2/v4"
)

const (
	// DefaultClient default clientID.
	DefaultClient = "tkeel"
	// DefaultClientSecurity default config .
	DefaultClientSecurity = "tkeel"
	// DefaultClientDomain default config.
	DefaultClientDomain = "https://tkeel.io"
)

// AuthorizeRequest authorization request.
type AuthorizeRequest struct {
	ResponseType        oauth2.ResponseType
	ClientID            string
	Scope               string
	RedirectURI         string
	State               string
	UserID              string
	CodeChallenge       string
	CodeChallengeMethod oauth2.CodeChallengeMethod
	AccessTokenExp      time.Duration
	Request             *http.Request
}
