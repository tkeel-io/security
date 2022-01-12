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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandStringWithPrefix(t *testing.T) {
	type args struct {
		prefix string
		len    int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "1", args: args{prefix: "usr", len: 16}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := RandStringWithPrefix(tt.args.prefix, tt.args.len); got != "" {
				t.Logf("RandStringWithPrefix() = %v", got)
			}
		})
	}
}

func TestRandString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"a", args{4}, false},
		{"b", args{6}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandBase64String(tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("RandBase64String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
