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
	"errors"

	"gorm.io/gorm"
)

type Tenant struct {
	ID     int    `json:"id" gorm:"primaryKey;type:int;comment:租户ID"`
	Title  string `json:"title" gorm:"type:varchar(128);comment:租户标题; not null;index"`
	Remark string `json:"remark" gorm:"type:varchar(255);comment:备注"`
}

func (Tenant) TableName() string {
	return "sys_t_tenant"
}

func (o *Tenant) Create() (err error) {
	err = _db.Create(o).Error
	return err
}

func (o *Tenant) Existed() (existed bool) {
	tenant := &Tenant{}
	_db.Where("title", o.Title).First(tenant)
	return tenant.ID != 0
}

func (o *Tenant) List(page *Page) (tenants []*Tenant, err error) {
	db := _db.Model(o)
	if page != nil {
		formatPage(db, page)
	}
	if o.ID != 0 {
		err = db.First(o).Error
		tenants = append(tenants, o)
	} else if o.Title != "" {
		err = db.Where("title LIKE ?", "%"+o.Title+"%").Find(&tenants).Error
	} else {
		err = db.Find(&tenants).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (o *Tenant) Delete() error {
	return _db.Delete(o).Error
}
