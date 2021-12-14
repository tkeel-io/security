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
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt"
)

func TestJWTAccessGenerate_Token(t *testing.T) {
	type fields struct {
		SignedKeyID  string
		SignedKey    []byte
		SignedMethod jwt.SigningMethod
	}
	type args struct {
		ctx          context.Context
		data         *jwt.MapClaims
		isGenRefresh bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &JWTAccessGenerate{
				SignedKeyID:  tt.fields.SignedKeyID,
				SignedKey:    tt.fields.SignedKey,
				SignedMethod: tt.fields.SignedMethod,
			}
			got, got1, err := a.Token(tt.args.ctx, tt.args.data, tt.args.isGenRefresh)
			if (err != nil) != tt.wantErr {
				t.Errorf("Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Token() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Token() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestJWTAccessGenerate_Valid(t *testing.T) {
	type fields struct {
		SignedKeyID  string
		SignedKey    []byte
		SignedMethod jwt.SigningMethod
	}
	type args struct {
		tokenStr string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantClaims interface{}
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &JWTAccessGenerate{
				SignedKeyID:  tt.fields.SignedKeyID,
				SignedKey:    tt.fields.SignedKey,
				SignedMethod: tt.fields.SignedMethod,
			}
			gotClaims, err := a.Valid(tt.args.tokenStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotClaims, tt.wantClaims) {
				t.Errorf("Valid() gotClaims = %v, want %v", gotClaims, tt.wantClaims)
			}
		})
	}
}
