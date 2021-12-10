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

func UUIDWithPrefix(prefix string, length int) string {
	if length == 0 {
		length = 16
	}
	uuid := make([]byte, length)
	_, err := rand.Read(uuid)
	if err != nil {
		panic("generate an uuid failed, error: " + err.Error())
	}
	return fmt.Sprintf("%s-%x", prefix, uuid)
}

func RandString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("rand string %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
