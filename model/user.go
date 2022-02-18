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
	"errors"
	"fmt"
	"time"

	"github.com/tkeel-io/security/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	ExternalID string    `json:"external_id" gorm:"type:varchar(128);default:'';uniqueIndex:user_tenant"`
	TenantID   string    `json:"tenant_id" gorm:"type:varchar(32);column:tenant_id;not null;uniqueIndex:user_tenant"`
	UserName   string    `json:"username" gorm:"type:varchar(64);column:username;not null;uniqueIndex:user_tenant"` // nolint
	Password   string    `json:"-" gorm:"type:varchar(128);column:password;not null"`
	NickName   string    `json:"nick_name" gorm:"type:varchar(128);comment:昵称"`
	Avatar     string    `json:"avatar" gorm:"type:varchar(255);comment:头像"`
	Email      string    `json:"email" gorm:"type:varchar(128);column:email;comment:邮箱"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
	if u.ID == "" {
		usrID, err := GenUserID()
		if err != nil {
			return err
		}
		u.ID = usrID
	}
	return u.Encrypt()
}

func (u *User) Existed(db *gorm.DB) (existed bool, err error) {
	err = db.First(u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	if err == nil {
		existed = true
	}
	return
}

func (u *User) Create(db *gorm.DB) error {
	var err error
	if u.ID == "" {
		u.ID, err = GenUserID()
		if err != nil {
			return err
		}
	}
	return db.Create(u).Error
}

func (u *User) CountInTenant(db *gorm.DB, tenantID string) (int64, error) {
	var i int64
	err := db.Model(u).Where("tenant_id = ?", tenantID).Count(&i).Error
	return i, err
}

func (u *User) Delete(db *gorm.DB) error {
	if u.ID == "" {
		return errors.New("empty user_id on delete")
	}
	err := db.Delete(u).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (u *User) DeleteAllInTenant(db *gorm.DB, tenantID string) error {
	return db.Delete(u, "tenant_id = ?", tenantID).Error
}

// QueryByCondition query by condition (todo fix page).
func (u *User) QueryByCondition(db *gorm.DB, condition map[string]interface{}, page *Page, keyWords string) (total int64, users []*User, err error) {
	if condition == nil {
		return total, nil, errors.New("query user condition is empty")
	}
	db = db.Model(&User{})
	if keyWords != "" {
		db = db.Where("username like ? or nick_name like ?", "%"+keyWords+"%", "%"+keyWords+"%")
	}
	db = db.Where(condition).Count(&total)
	if page != nil {
		db = FormatPage(db, page)
	}
	err = db.Find(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}
func (u *User) Update(db *gorm.DB, tenantID, userID string, updates map[string]interface{}) error {
	return db.Table(u.TableName()).Where("id = ? and tenant_id = ?", userID, tenantID).Updates(updates).Error
}

func MappingFromExternal(db *gorm.DB, externalID, name, email, tenantID string) (*User, error) {
	users := make([]*User, 0)
	condition := map[string]interface{}{
		"tenant_id":   tenantID,
		"external_id": externalID,
	}
	err := db.Model(&User{}).Where(condition).Find(&users).Error
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		userID, err := GenUserID()
		if err != nil {
			return nil, err
		}
		user := &User{
			ID:         userID,
			UserName:   name,
			ExternalID: externalID,
			Email:      email,
			TenantID:   tenantID,
		}
		err = user.Create(db)
		return user, err
	}
	return users[0], nil
}

func AuthenticateUser(db *gorm.DB, tenantID, username, password string) (*User, error) {
	user := &User{}
	err := db.Model(user).Where("tenant_id = ? and username = ? ", tenantID, username).First(user).Error
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate user password %w", err)
	}
	return user, nil
}

func GenUserID() (string, error) {
	return utils.RandStringWithPrefix("usr", 14)
}

func (u *User) FirstOrAssignCreate(db *gorm.DB, where User, assign User) error {
	result := db.Where(where).Assign(assign).FirstOrCreate(u)
	return result.Error
}
