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
	"time"
)

var startTime = time.Now()

// now returns time.Duration using stdlib time
func now() time.Duration { return time.Since(startTime) }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type interim struct {
	avgTotal float64
	lats     []float64
	numRes   int64
	successCount int64
	sizeTotal int64
	errCount int64
}

func NewInterim() *interim {
	return &interim{lats: make([]float64, 0, 100000)}
}

func (i *interim) reset() {
	if i.lats != nil {
		i.lats = i.lats[:0]
	}
	i.numRes = 0
	i.avgTotal = 0.0
	i.successCount = 0
	i.errCount = 0
	i.sizeTotal = 0
}

func (i *interim) finalize(total time.Duration) *Finalize {
	tps := float64(i.numRes) / total.Seconds()
	l := len(i.lats)
	average := i.avgTotal / float64(l)
	ls := latencies(i.lats)

	f := &Finalize{
		TimeStamp: time.Now().Format(time.RFC3339),
		TPS: tps,
		AvgDelay: average,
		Success: i.successCount,
		Err: i.errCount,
		Size: i.sizeTotal,
	}
	for _, lat := range ls {
		switch lat.Percentage {
		case 10:
			f.TP10 = lat.Latency
		case 25:
			f.TP25 = lat.Latency
		case 50:
			f.TP50 = lat.Latency
		case 75:
			f.TP75 = lat.Latency
		case 90:
			f.TP90 = lat.Latency
		case 95:
			f.TP95 = lat.Latency
		case 99:
			f.TP99 = lat.Latency
		}
	}
	return f
}

type Finalize struct {
	TimeStamp string        `json:"timestamp"`
	TPS       float64       `json:"tps"`
	AvgDelay  float64       `json:"avg_delay"`
	Success   int64         `json:"success"`
	Err       int64         `json:"err"`
	Size      int64         `json:"size"`
	TP10      float64       `json:"tp10"`
	TP25      float64       `json:"tp25"`
	TP50      float64       `json:"tp50"`
	TP75      float64       `json:"tp75"`
	TP90      float64       `json:"tp90"`
	TP95      float64       `json:"tp95"`
	TP99      float64       `json:"tp99"`
}

func latencies(lats []float64) []LatencyDistribution {
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]float64, len(pctls))
	j := 0
	for i := 0; i < len(lats) && j < len(pctls); i++ {
		current := i * 100 / len(lats)
		if current >= pctls[j] {
			data[j] = lats[i]
			j++
		}
	}
	res := make([]LatencyDistribution, len(pctls))
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			res[i] = LatencyDistribution{Percentage: pctls[i], Latency: data[i]}
		}
	}
	return res
}

type LatencyDistribution struct {
	Percentage int
	Latency    float64
}
