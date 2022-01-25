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
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Type         string `mapstructure:"type" json:"type" yaml:"type"`                             // 类型 mysql/postgresql.
	Host         string `mapstructure:"host" json:"host" yaml:"host"`                             // 服务器地址.
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                             // 端口.
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                       // 高级配置.
	Dbname       string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`                       // 数据库名.
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                 // 数据库用户名.
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                 // 数据库密码.
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"maxIdleConns"` // nolint
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"maxOpenConns"` // nolint
	LogLevel     string `mapstructure:"log_level" json:"log_level" yaml:"logLevel"`               // nolint
}

func (conf *DBConfig) MysqlDsn() string {
	if conf.Config != "" {
		return conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Dbname + "?" + conf.Config
	}
	return conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Dbname + "?" + "charset=utf8mb4&parseTime=True"
}

func (conf *DBConfig) PGDsn() string {
	return "host=" + conf.Host + " user=" + conf.Username + " password=" + conf.Password + " dbname=" + conf.Dbname + " port=" + conf.Port + " " + conf.Config
}

func logConfig(level string) *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	defaultLog := logger.New(newWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})
	switch strings.ToLower(level) {
	case "silent":
		config.Logger = defaultLog.LogMode(logger.Silent)
	case "error":
		config.Logger = defaultLog.LogMode(logger.Error)
	case "warn":
		config.Logger = defaultLog.LogMode(logger.Warn)
	case "info":
		config.Logger = defaultLog.LogMode(logger.Info)
	default:
		config.Logger = defaultLog.LogMode(logger.Info)
	}
	return config
}

func newWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

type writer struct {
	logger.Writer
}
