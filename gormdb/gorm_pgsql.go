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
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormPgSQLByConfig set up gorm db with pg config.
func GormPgSQLByConfig(conf DBConfig) (*gorm.DB, error) {
	var err error
	if conf.Dbname == "" {
		return nil, ErrDBNameNotFound
	}
	pgsqlConfig := postgres.Config{
		DSN:                  conf.PGDsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	if _db, err = gorm.Open(postgres.New(pgsqlConfig), logConfig(conf.LogLevel)); err != nil {
		return _db, fmt.Errorf("set up gorm db %w", err)
	}
	sqlDB, _ := _db.DB()
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	return _db, nil
}
