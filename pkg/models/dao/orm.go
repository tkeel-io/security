package dao

import (
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

func SetUp() {
	var err error
	dsn := "root:123456@tcp(139.198.108.153:3306)/tkeelauth?charset=utf8mb4&parseTime=True&loc=Local"
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
