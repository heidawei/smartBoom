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
package worker

import (
	"io"
	"time"
	"sync"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heidawei/smartBoom/register"
	"github.com/heidawei/smartBoom/executor"
)

const maxRes = 10000000

func init() {
	// do nothing
}

type Worker struct {
	// N is the total number of requests to make.
	N int

	// C is the concurrency level, the number of concurrent workers to run.
	C int

	// Qps is the rate limit in queries per second.
	QPS float64

	// Sampling interval.
	Interval time.Duration

	// Writer is where results will be written. If nil, results are written to stdout.
	Writer io.Writer

	ExecutorName string
	Config     map[string]interface{}

	cells   []*Cell
	output   *OutPut
	stopCh   chan struct{}
	done     chan struct{}
	once     sync.Once
}

func (b *Worker) writer() io.Writer {
	if b.Writer == nil {
		return os.Stdout
	}
	return b.Writer
}

func (b *Worker) Run() {
	for i := 0; i < b.C; i++ {
		if e,found := register.GetExecutor(b.ExecutorName); !found {
			fmt.Println("invalid runner name")
			os.Exit(-1)
		} else {
			exe := e(b.Config)
			exe.Init()
			cell := NewCell(b.QPS, exe)
			b.cells = append(b.cells, cell)
		}
	}
	b.done = make(chan struct{})
	b.stopCh = make(chan struct{})
	b.output = NewOutPut(getCurrentDirectory())
	// Run the reporter first, it polls the result channel until it is closed.
	go func() {
		b.runReporter()
	}()
	b.runWorkers()
	b.Finish()
}

func (b *Worker) Stop() {
	b.once.Do(func() {
		// Send stop signal so that workers can stop gracefully.
		for _, c := range b.cells {
			c.stop()
		}
		close(b.stopCh)
		b.finish()
	})
}

func (b *Worker) Finish() {
	b.once.Do(func() {
		close(b.stopCh)
		b.finish()
	})
}

func (b *Worker) finish() {
	// Wait until the reporter is done.
	<-b.done
	// TODO report
    b.output.Save()
}

func (b *Worker) runWorkers() {
	var wg sync.WaitGroup
	wg.Add(b.C)
	// Ignore the case where b.N % b.C != 0.
	for i := 0; i < b.C; i++ {
		go func(index int) {
			b.cells[index].run(index, b.N/b.C)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func (b *Worker) runReporter() {
	r := NewInterim()
	start := now()
	rss := make([][]*executor.Result, len(b.cells))
	collector := func(total time.Duration) {
		for i, cell := range b.cells {
			rs := cell.reset()
			rss[i] = rs
		}
		for _, rs := range rss {
			for _, res := range rs {
				if res.Count == 0 {
					r.numRes++
				} else {
					r.numRes += int64(res.Count)
				}
				if res.Err != nil {
					r.errCount++
				} else {
					r.successCount++
					r.avgTotal += res.Duration.Seconds()
					if len(r.lats) < maxRes {
						r.lats = append(r.lats, res.Duration.Seconds())
					}
					if res.ContentLength > 0 {
						r.sizeTotal += res.ContentLength
					}
				}
			}
		}
        f := r.finalize(total)
        b.output.Write(f)
		r.reset()
		for _, rs := range rss {
			executor.PutResultsToPool(rs)
		}
	}
	for {
		select {
		case <-b.stopCh:
			collector(now() - start)
		    close(b.done)
			return
		case <-time.After(b.Interval):
			collector(now() - start)
		    start = now()
		}
	}
}

type Cell struct {
	sync.Mutex
	qps      float64
	stopCh   chan struct{}
	runner   executor.Executor
	results  []*executor.Result
}

func NewCell(qps float64, runner executor.Executor) *Cell {
	return &Cell{qps: qps, stopCh: make(chan struct{}), runner: runner, results: executor.GetResultsFromPool()}
}

func (c *Cell) run(base, n int) {
	var throttle <-chan time.Time
	if c.qps > 0 {
		throttle = time.Tick(time.Duration(1e6/(c.qps)) * time.Microsecond)
	}

	for i := 0; i < n; {
		// Check if application is stopped. Do not send into a closed channel.
		select {
		case <-c.stopCh:
			return
		default:
			if c.qps > 0 {
				<-throttle
			}
			res := c.runner.Do(base, i, n)
			i += res.Count
			c.Lock()
			c.results = append(c.results, res)
			c.Unlock()
		}
	}
}

func (c *Cell) stop() {
	close(c.stopCh)
}

func (c *Cell) reset() []*executor.Result {
	c.Lock()
	defer c.Unlock()
	rs := c.results
	c.results = executor.GetResultsFromPool()
	return rs
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("get file path failed ", err)
		os.Exit(-1)
	}
	return strings.Replace(dir, "\\", "/", -1)
}