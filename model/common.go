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

package model

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	// DefaultPageNum default page num.
	DefaultPageNum = 15
)

type Page struct {
	PageNum      int    `json:"page_num"`
	PageSize     int    `json:"page_size"`
	OrderBy      string `json:"order_by"`
	IsDescending bool   `json:"is_descending"` // false:ascending;true:descending
}

func FormatPage(db *gorm.DB, page *Page) *gorm.DB {
	if page.PageNum <= 0 {
		page.PageNum = 1
	}
	if page.PageSize > 0 {
		db = db.Offset((page.PageNum - 1) * page.PageSize).Limit(page.PageSize)
	}

	if page.OrderBy != "" {
		if page.IsDescending {
			desc := fmt.Sprintf("%s desc", page.OrderBy)
			db = db.Order(desc)
		} else {
			db = db.Order(page.OrderBy)
		}
	}
	return db
}
