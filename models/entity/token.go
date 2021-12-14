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

import (
	"bytes"
	"context"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Token struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Owner      string `json:"owner"`
	CreatedAt  int64  `json:"created_at"`
	ExpiredAt  int64  `json:"expired_at"`
}

func (token *Token) MD5ID(i *int) string {
	*i++
	buf := bytes.NewBufferString(token.EntityID)
	buf.WriteString(token.EntityType)
	buf.WriteString(token.Owner)
	buf.WriteString(strconv.FormatInt(token.CreatedAt, 10))
	access := strings.ToUpper(base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes()).String())))
	access = strings.TrimRight(access, "=")
	return access
}

type TokenOperator interface {
	CreateToken(ctx context.Context, entity *Token) (token string, err error)
	GetEntityInfo(ctx context.Context, token string) (entity *Token, err error)
}
