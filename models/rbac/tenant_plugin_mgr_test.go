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

package rbac

import (
	"os"
	"testing"

	"github.com/tkeel-io/security/apiserver/config"

	"github.com/stretchr/testify/assert"
)

var (
	dbname   = ""
	user     = ""
	password = ""
	host     = ""
	port     = ""
)

func TestMain(m *testing.M) {
	adapter := &config.MysqlConf{
		DBName:   dbname,
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
	_, err := NewRBACOperator(adapter)
	if err != nil {
		os.Exit(-1)
	}
	gpolicy := &GroupingPolicy{
		Subject: "t1user",
		Role:    "t1role",
		Domain:  "systenant",
	}
	AddGroupingPolicy(gpolicy)
	m.Run()
}

func TestAddPlugin(t *testing.T) {
	plugin := "plugin"
	tenantID := "t1"
	gotOk, err := AddTenantPlugin(tenantID, plugin)
	assert.Nil(t, err)
	assert.Equal(t, true, gotOk)
}

func TestDeletePlugin(t *testing.T) {
	tenantID := "t1"
	pluginID := "plugin"
	gotOk, err := DeleteTenantPlugin(tenantID, pluginID)
	assert.Nil(t, err)
	assert.Equal(t, true, gotOk)
}

func TestPermissible(t *testing.T) {
	tenantID := "t1"
	pluginID := "plugin"
	gotOk, err := TenantPluginPermissible(tenantID, pluginID)
	assert.Nil(t, err)
	assert.Equal(t, true, gotOk)
}

func TestTenantPlugins(t *testing.T) {
	tenantID := "t1"
	gotPluginIDs := TenantPlugins(tenantID)
	t.Log(gotPluginIDs)
}
