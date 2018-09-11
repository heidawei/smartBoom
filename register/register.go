// Copyright 2018 The hedawei Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
package register

import (
	"sync"
	"github.com/heidawei/smartBoom/executor"
)

type CreateExecutor func(config map[string]interface{}) executor.Executor

type Register struct {
	sync.RWMutex
	cache    map[string]CreateExecutor
}

var register *Register

func RegisterExecutor(name string, new CreateExecutor) bool {
	register.Lock()
	defer register.Unlock()
	if _, found := register.cache[name]; found {
		return false
	}
	register.cache[name] = new
	return true
}

func GetExecutor(name string) (CreateExecutor, bool) {
	register.RLock()
	defer register.RUnlock()
	if e, found := register.cache[name]; found {
		return e, true
	}
	return nil, false
}

func init() {
	register = &Register{cache: make(map[string]CreateExecutor)}
}
