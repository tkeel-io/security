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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// RandStringWithPrefix .
func RandStringWithPrefix(prefix string, l int) (string, error) {
	if l == 0 {
		l = 16
	}
	uuid := make([]byte, l)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", fmt.Errorf("generate an rand string failed, %w ", err)
	}
	return fmt.Sprintf("%s-%x", prefix, uuid), nil
}

// RandBase64String The final length is not the byte length, but the base64-encoded length.
func RandBase64String(bytesLen int) (string, error) {
	b := make([]byte, bytesLen)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("rand string %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
