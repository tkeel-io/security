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

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormMySQLByConfig set up gorm db with mysql config.
func GormMySQLByConfig(conf DBConfig) (*gorm.DB, error) {
	var err error
	if conf.Dbname == "" {
		return nil, ErrDBNameNotFound
	}
	mysqlConfig := mysql.Config{
		DSN:                       conf.MysqlDsn(), // DSN data source name.
		DefaultStringSize:         191,             // string 类型字段的默认长度.
		SkipInitializeWithVersion: false,
	}
	if _db, err = gorm.Open(mysql.New(mysqlConfig), logConfig(conf.LogLevel)); err != nil {
		return nil, fmt.Errorf("set up gorm db %w", err)
	}
	sqlDB, _ := _db.DB()
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 10
	}
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 10
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	return _db, nil
}
