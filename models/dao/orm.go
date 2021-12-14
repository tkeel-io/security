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
	"sync"

	"github.com/tkeel-io/security/apiserver/config"
	"github.com/tkeel-io/security/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	_db     *gorm.DB
	_dbOnce sync.Once
	_log    = logger.NewLogger("auth.models.dao")
)

func SetUp(conf *config.MysqlConf) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
	_dbOnce.Do(func() {
		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			_log.Error(err)
			return
		}
	})

	_db.AutoMigrate(new(User))
	_db.AutoMigrate(new(Tenant))
}
