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
	"fmt"

	"github.com/casbin/casbin/v2"
)

type TenantPluginMgr interface {
	AddTenantPlugin(tenantID, pluginID string) (bool, error)
	DeleteTenantPlugin(tenantID, pluginID string) (bool, error)
	ListTenantPlugins(tenantID string) []string
	TenantPluginPermissible(tenantID, pluginID string) (bool, error)
	OnCreateTenant(tenantID string) (bool, error)
}

const (
	// SysTenant Abstract tenant.
	SysTenant          = "_systenant"
	SysPluginActionUse = "_use"
	SysRole            = "_role"
	sysUser            = "_user"
)

type TenantPluginOperator struct {
	RBACOperator *casbin.SyncedEnforcer
}

func NewTenantPluginOperator(opt *casbin.SyncedEnforcer) TenantPluginMgr {
	return &TenantPluginOperator{RBACOperator: opt}
}

func (t *TenantPluginOperator) AddTenantPlugin(tenantID, pluginID string) (ok bool, err error) {
	role := fmt.Sprintf("%s%s", SysRole, tenantID)
	policy := []string{role, SysTenant, pluginID, SysPluginActionUse}
	ok, err = t.RBACOperator.AddPolicy(policy)
	if err != nil {
		err = fmt.Errorf("%w", err)
	}
	return
}

func (t *TenantPluginOperator) DeleteTenantPlugin(tenantID, pluginID string) (ok bool, err error) {
	role := fmt.Sprintf("%s%s", SysRole, tenantID)
	ok, err = t.RBACOperator.RemovePolicy(role, SysTenant, pluginID, SysPluginActionUse)
	if err != nil {
		err = fmt.Errorf("%w", err)
	}
	return
}

func (t *TenantPluginOperator) ListTenantPlugins(tenantID string) (pluginIDs []string) {
	user := fmt.Sprintf("%s%s", sysUser, tenantID)
	permissions := t.RBACOperator.GetPermissionsForUserInDomain(user, SysTenant)
	for i := range permissions {
		plugin := permissions[i][2]
		pluginIDs = append(pluginIDs, plugin)
	}
	return pluginIDs
}

func (t *TenantPluginOperator) TenantPluginPermissible(tenantID, pluginID string) (ok bool, err error) {
	user := fmt.Sprintf("%s%s", sysUser, tenantID)
	ok, err = t.RBACOperator.Enforce(user, SysTenant, pluginID, SysPluginActionUse)
	if err != nil {
		err = fmt.Errorf("%w", err)
	}
	return
}

func (t *TenantPluginOperator) OnCreateTenant(tenantID string) (ok bool, err error) {
	gpolicy := []string{
		fmt.Sprintf("%s%s", sysUser, tenantID),
		fmt.Sprintf("%s%s", SysRole, tenantID),
		SysTenant,
	}
	ok, err = t.RBACOperator.AddGroupingPolicy(gpolicy)
	return
}
