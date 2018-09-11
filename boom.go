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
package main

import (
	"flag"
	"runtime"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/heidawei/smartBoom/worker"
	"github.com/dustin/gojson"
)

var (
	c = flag.Int("c", 50, "")
	n = flag.Int("n", 200, "")
	q = flag.Float64("q", 0, "")
	i = flag.Duration("i", time.Second, "")
	name = flag.String("name", "", "")
	config = flag.String("config", "", "")

	cpus = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")
)

var usage = `Usage: smartBoom [options...] <url>

Options:
  -n  Number of requests to run. Default is 200.
  -c  Number of requests to run concurrently. Total number of requests cannot
      be smaller than the concurrency level. Default is 50.
  -q  Rate limit, in queries per second (QPS). Default is no rate limit.

  -i  Interval of collector report, unit second. Default is 1 second

  -name Name of executor.
  -config Executor config json file.
  

  -cpus                 Number of used cpu cores.
                        (default for current machine is %d cores)
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}

	flag.Parse()

	runtime.GOMAXPROCS(*cpus)
	num := *n
	conc := *c
	q := *q
	executor := *name
	interval := *i

	if num <= 0 || conc <= 0 || interval <= time.Millisecond {
		usageAndExit("-n and -c cannot be smaller than 1.")
	}

	if interval <= time.Millisecond {
		usageAndExit("-i cannot be smaller than 1 ms")
	}

	if num < conc {
		usageAndExit("-n cannot be less than -c.")
	}

	if len(executor) == 0 {
		usageAndExit("-name cannot be empty.")
	}

	cfg := make(map[string]interface{})
	if len(*config) > 0 {
		f, err := os.Open(*config)
		if err != nil {
			usageAndExit(fmt.Sprintf("config file %s is invalid, err %v", *config, err))
		}
		err = json.NewDecoder(f).Decode(&cfg)
		f.Close()
		if err != nil {
			usageAndExit(fmt.Sprintf("config file %s is not json file, err %v", *config, err))
		}
	}

	w := &worker.Worker{
		N:                  num,
		C:                  conc,
		QPS:                q,
		Interval:           interval,
		ExecutorName:       executor,
		Config:             cfg,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		w.Stop()
	}()
	w.Run()
}

func errAndExit(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

