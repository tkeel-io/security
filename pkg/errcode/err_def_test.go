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
package errcode

import (
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	type args struct {
		code int
		msg  string
	}
	tests := []struct {
		name string
		args args
		want *SrvError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.code, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSrvError_Code(t *testing.T) {
	type fields struct {
		code    int
		msg     string
		details []string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SrvError{
				code:    tt.fields.code,
				msg:     tt.fields.msg,
				details: tt.fields.details,
			}
			if got := e.Code(); got != tt.want {
				t.Errorf("Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSrvError_Details(t *testing.T) {
	type fields struct {
		code    int
		msg     string
		details []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SrvError{
				code:    tt.fields.code,
				msg:     tt.fields.msg,
				details: tt.fields.details,
			}
			if got := e.Details(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Details() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSrvError_Error(t *testing.T) {
	type fields struct {
		code    int
		msg     string
		details []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := SrvError{
				code:    tt.fields.code,
				msg:     tt.fields.msg,
				details: tt.fields.details,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSrvError_Msg(t *testing.T) {
	type fields struct {
		code    int
		msg     string
		details []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SrvError{
				code:    tt.fields.code,
				msg:     tt.fields.msg,
				details: tt.fields.details,
			}
			if got := e.Msg(); got != tt.want {
				t.Errorf("Msg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSrvError_Msgf(t *testing.T) {
	type fields struct {
		code    int
		msg     string
		details []string
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SrvError{
				code:    tt.fields.code,
				msg:     tt.fields.msg,
				details: tt.fields.details,
			}
			if got := e.Msgf(tt.args.args); got != tt.want {
				t.Errorf("Msgf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSrvError_WithDetails(t *testing.T) {
	type fields struct {
		code    int
		msg     string
		details []string
	}
	type args struct {
		details []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *SrvError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SrvError{
				code:    tt.fields.code,
				msg:     tt.fields.msg,
				details: tt.fields.details,
			}
			if got := e.WithDetails(tt.args.details...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}
