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
	"fmt"
)

type SrvError struct {
	code    int
	msg     string
	details []string
}

var codes = map[int]struct{}{}

// NewError create a error.
func NewError(code int, msg string) *SrvError {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("code %d is exsit, please change one", code))
	}
	codes[code] = struct{}{}
	return &SrvError{code: code, msg: msg}
}

// SrvError return a error string.
func (e SrvError) Error() string {
	return fmt.Sprintf("code：%d, msg:：%s", e.Code(), e.Msg())
}

// Code Authenticate return error code.
func (e *SrvError) Code() int {
	return e.code
}

// Msg return error msg.
func (e *SrvError) Msg() string {
	return e.msg
}

// Msgf format error string.
func (e *SrvError) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

// Details return more error details.
func (e *SrvError) Details() []string {
	return e.details
}

// WithDetails return err with detail.
func (e *SrvError) WithDetails(details ...string) *SrvError {
	newError := *e
	newError.details = []string{}
	newError.details = append(newError.details, details...)

	return &newError
}
