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
package executor

import (
	"time"
	"sync"
)

type Result struct {
	Err           error
	StatusCode    int
	Duration      time.Duration
	ContentLength int64
	// default 1
	Count         int
}

type Executor interface {
	Init()
	// return num of message do
	Do(base, index, n int) *Result
}


var resultsPool = &sync.Pool{
	New: func() interface{} {
		return make([]*Result, 0, 100000)
	},
}

func GetResultsFromPool() []*Result {
	return resultsPool.Get().([]*Result)
}

func PutResultsToPool(res []*Result) {
	res = res[:0]
	resultsPool.Put(res)
}
