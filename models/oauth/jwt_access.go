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

package oauth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// EntityClaims jwt claims.
type EntityClaims struct {
	Audience   string `json:"aud,omitempty"`
	ExpiresAt  int64  `json:"exp,omitempty"`
	ID         string `json:"jti,omitempty"`
	IssuedAt   int64  `json:"iat,omitempty"`
	Issuer     string `json:"iss,omitempty"`
	NotBefore  int64  `json:"nbf,omitempty"`
	Subject    string `json:"sub,omitempty"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

// Valid claims verification.
func (a *EntityClaims) Valid() error {
	if time.Unix(a.ExpiresAt, 0).Before(time.Now()) {
		return errors.New("invalid access")
	}
	return nil
}

// NewJWTAccessGenerate create to generate the jwt access token instance.
func NewJWTAccessGenerate(kid string, key []byte, method jwt.SigningMethod) *JWTAccessGenerate {
	return &JWTAccessGenerate{
		SignedKeyID:  kid,
		SignedKey:    key,
		SignedMethod: method,
	}
}

// JWTAccessGenerate generate the jwt access token.
type JWTAccessGenerate struct {
	SignedKeyID  string
	SignedKey    []byte
	SignedMethod jwt.SigningMethod
}

// Token based on the UUID generated token.
func (a *JWTAccessGenerate) Token(ctx context.Context, data *jwt.MapClaims, isGenRefresh bool) (string, string, error) {
	token := jwt.NewWithClaims(a.SignedMethod, data)
	if a.SignedKeyID != "" {
		token.Header["kid"] = a.SignedKeyID
	}
	var key interface{}
	if a.isEs() {
		v, err := jwt.ParseECPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", fmt.Errorf("ParseECPrivateKeyFromPEM %w", err)
		}
		key = v
	} else if a.isRsOrPS() {
		v, err := jwt.ParseRSAPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", fmt.Errorf("ParseRSAPrivateKeyFromPEM %w", err)
		}
		key = v
	} else if a.isHs() {
		key = a.SignedKey
	} else {
		return "", "", errors.New("unsupported sign method")
	}

	access, err := token.SignedString(key)
	if err != nil {
		return "", "", fmt.Errorf("SignedString %w", err)
	}
	refresh := ""

	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}

func (a *JWTAccessGenerate) Valid(tokenStr string) (claims interface{}, err error) {
	var ok bool
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		_log.Info(a.SignedKey)
		return a.SignedKey, nil
	})
	_log.Info(token, err)
	if claims, ok = token.Claims.(jwt.MapClaims); ok && token.Valid {
		_log.Info(claims)
		return claims, nil
	}

	return
}

func (a *JWTAccessGenerate) isEs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "ES")
}

func (a *JWTAccessGenerate) isRsOrPS() bool {
	isRs := strings.HasPrefix(a.SignedMethod.Alg(), "RS")
	isPs := strings.HasPrefix(a.SignedMethod.Alg(), "PS")
	return isRs || isPs
}

func (a *JWTAccessGenerate) isHs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "HS")
}
