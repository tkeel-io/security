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
package v1

type AddRoleInDomainIn struct {
	Role string `json:"role"`
}

type AddPermissionIn struct {
	PermissionObject string `json:"permission_object"`
	PermissionAction string `json:"permission_action"`
}

type PermissionCheck struct {
	UserID           string `json:"user_id"`
	PermissionObject string `json:"permission_object"`
	PermissionAction string `json:"permission_action"`
}

type PermissionCheckResult struct {
	Allowed bool `json:"allowed"`
}

type AddRoleForUserIn struct {
	UserIDS []string `json:"user_ids"`
	Roles   []string `json:"roles"`
}

type DelRoleForUserIn struct {
	UserIDS []string `json:"user_ids"`
	Roles   []string `json:"roles"`
}

type UserPermission struct {
	Role             string `json:"role"`
	PermissionObject string `json:"permission_object"`
	PermissionAction string `json:"permission_action"`
}
