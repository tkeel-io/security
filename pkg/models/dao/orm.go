package dao

import (
	"fmt"
	"github.com/tkeel-io/security/pkg/apiserver/config"
	"sync"

	"github.com/tkeel-io/security/pkg/logger"

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
	_log.Info(dsn)
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
