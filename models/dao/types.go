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

package dao

import (
	"fmt"

	"github.com/tkeel-io/security/utils"

	"gorm.io/gorm"
)

const (
	// DefaultPageNum 15.
	DefaultPageNum = 15
	_userPrefix    = "usr"
	_defaultLen    = 16
)

type Page struct {
	PageIndex  int                    `json:"page_index"`
	PageNum    int                    `json:"page_num"`
	SortBy     string                 `json:"sort_by"`
	Descending bool                   `json:"descending"`
	Condition  map[string]interface{} `json:"condition"`
}

func formatPage(db *gorm.DB, page *Page) {
	if page.Condition != nil {
		for k, v := range page.Condition {
			db.Where(k, v)
		}
	}

	if page.PageNum == 0 {
		page.PageNum = DefaultPageNum
	}
	db.Offset((page.PageIndex - 1) * page.PageNum).Limit(page.PageNum)
	if page.SortBy != "" {
		if page.Descending {
			desc := fmt.Sprintf("%s desc", page.SortBy)
			db.Order(desc)
		} else {
			db.Order(page.SortBy)
		}
	}
}

func GenUserID(tenantID int) string {
	return utils.UUIDWithPrefix(fmt.Sprintf("%s-%d", _userPrefix, tenantID), _defaultLen)
}
