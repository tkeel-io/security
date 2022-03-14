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

package casbin

import "errors"

var (
	errInvalidParam = errors.New("invalid param")

	_textCasbinModel = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && (r.obj == p.obj || p.obj == '*') && ( r.act == p.act || p.act == '*') 
`
)

type Policy struct {
	Role   string
	Domain string
	Object string
	Action string
}

type GroupingPolicy struct {
	Subject string
	Role    string
	Domain  string
}

type RequestPolicy struct {
	Subject string
	Domain  string
	Object  string
	Action  string
}

func (p *Policy) Valid() error {
	if p.Role == "" || p.Domain == "" || p.Object == "" || p.Action == "" {
		return errInvalidParam
	}
	return nil
}
func (g *GroupingPolicy) Valid() error {
	if g.Role == "" || g.Domain == "" || g.Subject == "" {
		return errInvalidParam
	}
	return nil
}

func (r *RequestPolicy) Valid() error {
	if r.Object == "" || r.Domain == "" || r.Subject == "" || r.Action == "" {
		return errInvalidParam
	}
	return nil
}
