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

	"github.com/tkeel-io/security/constants"
)

func AddTenantPlugin(tenantID string, pluginID string) (ok bool, err error) {
	role := fmt.Sprintf("%srole", tenantID)
	policy := []string{role, constants.SysTenant, pluginID, constants.SysPluginActionUse}
	ok, err = _enforcer.AddPolicy(policy)
	if err != nil {
		err = fmt.Errorf("%w", err)
	}
	return
}

func DeleteTenantPlugin(tenantID string, pluginID string) (ok bool, err error) {
	role := fmt.Sprintf("%srole", tenantID)
	ok, err = _enforcer.RemovePolicy(role, constants.SysTenant, pluginID, constants.SysPluginActionUse)
	if err != nil {
		err = fmt.Errorf("%w", err)
	}
	return
}

func TenantPluginPermissible(tenantID string, pluginID string) (ok bool, err error) {
	user := fmt.Sprintf("%suser", tenantID)
	ok, err = _enforcer.Enforce(user, constants.SysTenant, pluginID, constants.SysPluginActionUse)
	if err != nil {
		err = fmt.Errorf("%w", err)
	}
	return
}

func TenantPlugins(tenantID string) (pluginIDs []string) {
	user := fmt.Sprintf("%suser", tenantID)
	permissions := _enforcer.GetPermissionsForUserInDomain(user, constants.SysTenant)
	for i := range permissions {
		plugin := permissions[i][2]
		pluginIDs = append(pluginIDs, plugin)
	}
	return pluginIDs
}
