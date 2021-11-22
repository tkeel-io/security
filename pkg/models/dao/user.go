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

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         string `json:"user_id" gorm:"primaryKey"`
	ExternalID string `json:"external_id" gorm:"type:varchar(128)"`
	TenantID   int    `json:"tenant_id" gorm:"type:int;column:tenant_id;not null;uniqueIndex:user_tenant"`
	UserName   string `json:"username" gorm:"type:varchar(64);column:username;not null;uniqueIndex:user_tenant"`
	Password   string `json:"-" gorm:"type:varchar(128);column:password;not null"`
	NickName   string `json:"nick_name" gorm:"type:varchar(128);comment:昵称"`
	Avatar     string `json:"avatar" gorm:"type:varchar(255);comment:头像"`
	Email      string `json:"email" gorm:"type:varchar(128);column:email;comment:邮箱"`
}

func (User) TableName() string {
	return "sys_t_user"
}

func (u *User) Encrypt() (err error) {
	if u.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost); err != nil {
		return
	}
	u.Password = string(hash)

	return
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	return u.Encrypt()
}

func (u *User) Existed() (existed bool, err error) {
	err = _db.First(u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	if err == nil {
		existed = true
	}
	return
}

func (u *User) Create() error {
	if u.ID == "" {
		return errors.New("user id empty")
	}
	return _db.Create(u).Error
}

func (u *User) Delete() error {
	if u.ID == "" {
		return errors.New("empty user_id on delete")
	}
	err := _db.Delete(u).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

// QueryByCondition query by condition (todo fix page).
func (u *User) QueryByCondition(condition map[string]interface{}) (users []*User, err error) {
	if condition == nil {
		return nil, errors.New("query user condition is empty")
	}
	userID, ok := condition["user_id"]
	if ok {
		user := &User{ID: userID.(string)}
		_db.First(user)
	}
	err = _db.Model(&User{}).Where(condition).Find(&users).Error
	return
}
