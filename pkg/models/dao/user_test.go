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
	"reflect"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestSetUp(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestTenant_Create(t *testing.T) {
	type fields struct {
		ID     int
		Title  string
		Remark string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Tenant{
				ID:     tt.fields.ID,
				Title:  tt.fields.Title,
				Remark: tt.fields.Remark,
			}
			if err := o.Create(); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTenant_Delete(t *testing.T) {
	type fields struct {
		ID     int
		Title  string
		Remark string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Tenant{
				ID:     tt.fields.ID,
				Title:  tt.fields.Title,
				Remark: tt.fields.Remark,
			}
			if err := o.Delete(); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTenant_Existed(t *testing.T) {
	type fields struct {
		ID     int
		Title  string
		Remark string
	}
	tests := []struct {
		name        string
		fields      fields
		wantExisted bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Tenant{
				ID:     tt.fields.ID,
				Title:  tt.fields.Title,
				Remark: tt.fields.Remark,
			}
			if gotExisted := o.Existed(); gotExisted != tt.wantExisted {
				t.Errorf("Existed() = %v, want %v", gotExisted, tt.wantExisted)
			}
		})
	}
}

func TestTenant_List(t *testing.T) {
	type fields struct {
		ID     int
		Title  string
		Remark string
	}
	type args struct {
		page *Page
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantTenants []*Tenant
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Tenant{
				ID:     tt.fields.ID,
				Title:  tt.fields.Title,
				Remark: tt.fields.Remark,
			}
			gotTenants, err := o.List(tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTenants, tt.wantTenants) {
				t.Errorf("List() gotTenants = %v, want %v", gotTenants, tt.wantTenants)
			}
		})
	}
}

func TestTenant_TableName(t *testing.T) {
	type fields struct {
		ID     int
		Title  string
		Remark string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := Tenant{
				ID:     tt.fields.ID,
				Title:  tt.fields.Title,
				Remark: tt.fields.Remark,
			}
			if got := te.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	type fields struct {
		ID         string
		ExternalID string
		TenantID   int
		UserName   string
		Password   string
		NickName   string
		Avatar     string
		Salt       string
		Email      string
	}
	type args struct {
		in0 *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:         tt.fields.ID,
				ExternalID: tt.fields.ExternalID,
				TenantID:   tt.fields.TenantID,
				UserName:   tt.fields.UserName,
				Password:   tt.fields.Password,
				NickName:   tt.fields.NickName,
				Avatar:     tt.fields.Avatar,
				Email:      tt.fields.Email,
			}
			if err := u.BeforeCreate(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("BeforeCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_Create(t *testing.T) {
	type fields struct {
		ID         string
		ExternalID string
		TenantID   int
		UserName   string
		Password   string
		NickName   string
		Avatar     string
		Salt       string
		Email      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:         tt.fields.ID,
				ExternalID: tt.fields.ExternalID,
				TenantID:   tt.fields.TenantID,
				UserName:   tt.fields.UserName,
				Password:   tt.fields.Password,
				NickName:   tt.fields.NickName,
				Avatar:     tt.fields.Avatar,
				Email:      tt.fields.Email,
			}
			if err := u.Create(); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	type fields struct {
		ID         string
		ExternalID string
		TenantID   int
		UserName   string
		Password   string
		NickName   string
		Avatar     string
		Salt       string
		Email      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:         tt.fields.ID,
				ExternalID: tt.fields.ExternalID,
				TenantID:   tt.fields.TenantID,
				UserName:   tt.fields.UserName,
				Password:   tt.fields.Password,
				NickName:   tt.fields.NickName,
				Avatar:     tt.fields.Avatar,
				Email:      tt.fields.Email,
			}
			if err := u.Delete(); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_Encrypt(t *testing.T) {
	type fields struct {
		ID         string
		ExternalID string
		TenantID   int
		UserName   string
		Password   string
		NickName   string
		Avatar     string
		Salt       string
		Email      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{{fields: fields{Password: "123456"}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:         tt.fields.ID,
				ExternalID: tt.fields.ExternalID,
				TenantID:   tt.fields.TenantID,
				UserName:   tt.fields.UserName,
				Password:   tt.fields.Password,
				NickName:   tt.fields.NickName,
				Avatar:     tt.fields.Avatar,
				Email:      tt.fields.Email,
			}
			if err := u.Encrypt(); (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(u.Password)
			err := bcrypt.CompareHashAndPassword([]byte("$2a$10$DTMnHdYFEB5/n.CevILL8OgTxs6mvTdCLwYYlCFantvdhpC/5KCWy"), []byte("123456"))
			err2 := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("123456"))
			t.Log(err)
			t.Log(err2)
		})
	}
}

func TestUser_Existed(t *testing.T) {
	type fields struct {
		ID         string
		ExternalID string
		TenantID   int
		UserName   string
		Password   string
		NickName   string
		Avatar     string
		Salt       string
		Email      string
	}
	tests := []struct {
		name        string
		fields      fields
		wantExisted bool
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:         tt.fields.ID,
				ExternalID: tt.fields.ExternalID,
				TenantID:   tt.fields.TenantID,
				UserName:   tt.fields.UserName,
				Password:   tt.fields.Password,
				NickName:   tt.fields.NickName,
				Avatar:     tt.fields.Avatar,
				Email:      tt.fields.Email,
			}
			gotExisted, err := u.Existed()
			if (err != nil) != tt.wantErr {
				t.Errorf("Existed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExisted != tt.wantExisted {
				t.Errorf("Existed() gotExisted = %v, want %v", gotExisted, tt.wantExisted)
			}
		})
	}
}

func TestUser_TableName(t *testing.T) {
	type fields struct {
		ID         string
		ExternalID string
		TenantID   int
		UserName   string
		Password   string
		NickName   string
		Avatar     string
		Salt       string
		Email      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := User{
				ID:         tt.fields.ID,
				ExternalID: tt.fields.ExternalID,
				TenantID:   tt.fields.TenantID,
				UserName:   tt.fields.UserName,
				Password:   tt.fields.Password,
				NickName:   tt.fields.NickName,
				Avatar:     tt.fields.Avatar,
				Email:      tt.fields.Email,
			}
			if got := us.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
