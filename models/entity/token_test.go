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

package entity

import "testing"

func TestEntityToken_MD5ID(t *testing.T) {
	type fields struct {
		EntityID   string
		EntityType string
		Owner      string
		CreatedAt  int64
		ExpiredAt  int64
	}
	tests := []struct {
		i      int
		name   string
		fields fields
	}{{4, "testmd5id", fields{"id", "type", "owner", 123456789, 987654321}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := &Token{
				EntityID:   tt.fields.EntityID,
				EntityType: tt.fields.EntityType,
				Owner:      tt.fields.Owner,
				CreatedAt:  tt.fields.CreatedAt,
				ExpiredAt:  tt.fields.ExpiredAt,
			}
			i := 0
			var got string
			for got = token.MD5ID(&i); got == "" && i < tt.i; got = token.MD5ID(&i) {
				t.Log(got, i)
			}
			t.Logf("get %s at %d", got, i)
		})
	}
}
