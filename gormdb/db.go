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

package gormdb

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	_db *gorm.DB
	// ErrDBNotSetUp get gorm db but not set up.
	ErrDBNotSetUp = errors.New("gorm db not set up")
	// ErrDBNameNotFound config db name not found.
	ErrDBNameNotFound = errors.New("db config dbname not found")
)

func GetGormDB() (*gorm.DB, error) {
	if _db == nil {
		return nil, ErrDBNotSetUp
	}
	return _db, nil
}

func SetUp(conf DBConfig) (*gorm.DB, error) {
	switch strings.ToLower(conf.Type) {
	case "mysql":
		return GormMySQLByConfig(conf)
	case "pgsql":
		return GormPgSQLByConfig(conf)
	default:
		return GormMySQLByConfig(conf)
	}
}
