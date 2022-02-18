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
	"github.com/tkeel-io/security/utils"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(128);not null;uniqueIndex:role_tenant"`
	TenantID    string `json:"tenant_id" gorm:"type:varchar(32);default:'';uniqueIndex:role_tenant"`
	Description string `json:"description" gorm:"type:varchar(256);default:''"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Role) TableName() string {
	return "sys_t_role"
}

func (dao *Role) BeforeCreate(_ *gorm.DB) error {
	if dao.ID == "" {
		roleID, err := GenRoleID()
		if err != nil {
			return err
		}
		dao.ID = roleID
	}
	return nil
}

func (dao *Role) Create(db *gorm.DB) error {
	return db.Create(dao).Error
}

func (dao *Role) IsExisted(db *gorm.DB, where map[string]interface{}) (bool, error) {
	roles := []Role{}
	err := db.Where(where).Find(&roles).Error
	if err != nil {
		return false, err
	}
	return len(roles) != 0, nil
}

func (dao *Role) List(db *gorm.DB, where map[string]interface{}, page *Page, keywords string) (total int64, roles []*Role, err error) {
	if where != nil {
		db = db.Where(where)
	}
	if keywords != "" {
		db = db.Where("name like ? or description like ?", "%"+keywords+"%", "%"+keywords+"%")
	}
	if page != nil {
		db = FormatPage(db, page)
	}
	err = db.Find(&roles).Count(&total).Error
	return
}

func (dao *Role) Update(db *gorm.DB, where map[string]interface{}, updates map[string]interface{}) (affected int64, err error) {
	result := db.Table(dao.TableName()).Where(where).Updates(updates)
	return result.RowsAffected, result.Error
}

func (dao *Role) Delete(db *gorm.DB, where map[string]interface{}) (affected int64, err error) {
	result := db.Delete(dao, where)
	return result.RowsAffected, result.Error
}

func GenRoleID() (string, error) {
	return utils.RandStringWithPrefix("role", 12)
}
