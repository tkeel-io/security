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
	"github.com/tkeel-io/security/pkg/logger"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
)

var (
	_log      = logger.NewLogger("auth.models.rbac")
	_enforcer *casbin.SyncedEnforcer
)

func NewSyncedEnforcer() (enforcer *casbin.SyncedEnforcer, err error) {
	adapter, err := xormadapter.NewAdapter("mysql", "root:123456@tcp(139.198.108.153:3306)/tkeelauth", true)
	if err != nil {
		_log.Error(err)
		return
	}
	casbinModel, err := model.NewModelFromString(_textCasbinModel)
	if err != nil {
		_log.Error(err)
		return
	}

	_enforcer, err = casbin.NewSyncedEnforcer(casbinModel, adapter)
	if err != nil {
		_log.Error(err)
		return
	}
	err = _enforcer.LoadPolicy()
	if err != nil {
		_log.Error(err)
		return
	}
	_enforcer.EnableLog(true)

	return _enforcer, err
}

func AddGroupingPolicy(gPolicy *GroupingPolicy) (ok bool, err error) {
	err = gPolicy.Valid()
	if err != nil {
		_log.Error(err)
		return
	}

	params := []string{gPolicy.Subject, gPolicy.Role, gPolicy.Domain}
	ok, err = _enforcer.AddGroupingPolicy(params)
	return
}

func AddPolicy(pPolicy *Policy) (ok bool, err error) {
	err = pPolicy.Valid()
	if err != nil {
		_log.Error(err)
		return
	}
	params := []string{pPolicy.Role, pPolicy.Domain, pPolicy.Object, pPolicy.Action}
	ok, err = _enforcer.AddPolicy(params)
	return
}

func Enforce(r *RequestPolicy) (ok bool, err error) {
	err = r.Valid()
	if err != nil {
		_log.Error(err)
		return
	}
	params := []string{r.Subject, r.Domain, r.Object, r.Action}
	ok, err = _enforcer.Enforce(params)
	return
}
