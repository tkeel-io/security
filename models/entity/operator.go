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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/tkeel-io/security/logger"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	_log = logger.NewLogger("security.models.entity")
)

type TokenOp struct {
	storeName string
	operator  dapr.Client
}

func NewEntityTokenOperator(storeName string, client dapr.Client) *TokenOp {
	if client == nil || storeName == "" {
		return nil
	}
	return &TokenOp{storeName: storeName, operator: client}
}

func (e *TokenOp) CreateToken(ctx context.Context, entity *Token) (token string, err error) {
	value, err := json.Marshal(entity)
	if err != nil {
		return "", fmt.Errorf("marshal entity token %w", err)
	}
	i := 0
	key := entity.MD5ID(&i)
	resultKey := ""
	var item *dapr.StateItem
	for item, _ = e.operator.GetState(ctx, e.storeName, key); item.Value == nil && i < 4; key = entity.MD5ID(&i) {
		err = e.operator.SaveBulkState(ctx, e.storeName, &dapr.SetStateItem{Key: key, Value: value, Etag: &dapr.ETag{Value: item.Etag},
			Options: &dapr.StateOptions{Concurrency: dapr.StateConcurrencyFirstWrite, Consistency: dapr.StateConsistencyStrong}})
		if err != nil {
			return "", fmt.Errorf("create token save state %w", err)
		}
		resultKey = key
		i = 4
	}
	if item.Value != nil {
		return "", errors.New("entity hash three repetitions on key ")
	}

	return resultKey, nil
}

func (e *TokenOp) GetEntityInfo(ctx context.Context, key string) (entity *Token, err error) {
	item, err := e.operator.GetState(ctx, e.storeName, key)
	if err != nil {
		return nil, fmt.Errorf("entity info get state %w", err)
	}
	entity = &Token{}
	if err = json.Unmarshal(item.Value, entity); err != nil {
		return nil, fmt.Errorf("unmarshal entity info  %w", err)
	}
	return
}

func NewGPRCClient(retry int, interval, gprcPort string) (dapr.Client, error) {
	var daprGRPCClient dapr.Client
	var err error
	inval, err := time.ParseDuration(interval)
	if err != nil {
		return nil, fmt.Errorf("error parse interval(%s): %w", interval, err)
	}
	if retry < 1 {
		retry = 1
	}
	for i := 0; i < retry; i++ {
		daprGRPCClient, err = dapr.NewClientWithPort(gprcPort)
		if err == nil {
			break
		}

		time.Sleep(inval)
		_log.Debugf("error new client: %s retry: %d", err, i)
	}
	if err != nil {
		return nil, fmt.Errorf("error new client with port(%s): %w", gprcPort, err)
	}
	return daprGRPCClient, nil
}
