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

package authorizer

const (
	// DecisionDeny means that an authorizer decided to deny the action.
	DecisionDeny Decision = iota
	// DecisionAllow means that an authorizer decided to allow the action.
	DecisionAllow
	// DecisionNoOpinion means that an authorizer has no opinion on whether to allow or deny an action.
	DecisionNoOpinion
)

type Decision int

type AuthZReq interface {
	GetUser() string
	GetTenant() string
	OptObject() string
	OptAction() string
}

type Authorizer interface {
	Authorize(req AuthZReq) (Decision, error)
}
