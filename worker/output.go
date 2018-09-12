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
	"fmt"
	"os"
	"path"
	"time"

	"github.com/tealeg/xlsx"
)

var Titles = []string{"timestamp", "TPS", "avg latency", "total success", "total fail",
                      "TP10", "TP25", "TP50", "TP75", "TP90", "TP95", "TP99"}

type OutPut struct {
	path string
	f    *xlsx.File
	sheet *xlsx.Sheet
	fs   []*Finalize
}

func NewOutPut(path string) *OutPut {
	f := xlsx.NewFile()
	sheet, err := f.AddSheet("statis")
	if err != nil {
		fmt.Println("xlsx add sheet failed ", err)
		os.Exit(-1)
	}
	// title
	r := sheet.AddRow()
	for _, title := range Titles {
		cell := r.AddCell()
		cell.Value = title
	}

	return &OutPut{f: f, sheet: sheet, path: path, fs: make([]*Finalize, 0, 1000)}
}

func (o *OutPut) Write(f *Finalize) {
	r := o.sheet.AddRow()
	cell := r.AddCell()
	// timestamp
	cell.SetDateTime(f.TimeStamp)
	cell = r.AddCell()
	// TPS
	cell.SetFloat(f.TPS)
	cell = r.AddCell()
	// avg
	cell.SetFloat(f.AvgDelay)
	cell = r.AddCell()
	// success
	cell.SetInt64(f.Success)
	cell = r.AddCell()
	// fail
	cell.SetInt64(f.Err)
	cell = r.AddCell()
	// TP10
	cell.SetFloat(f.TP10)
	cell = r.AddCell()
	// TP25
	cell.SetFloat(f.TP25)
	cell = r.AddCell()
	// TP50
	cell.SetFloat(f.TP50)
	cell = r.AddCell()
	// TP75
	cell.SetFloat(f.TP75)
	cell = r.AddCell()
	// TP90
	cell.SetFloat(f.TP90)
	cell = r.AddCell()
	// TP95
	cell.SetFloat(f.TP95)
	cell = r.AddCell()
	// TP99
	cell.SetFloat(f.TP99)
}

func (o *OutPut) Save() {
	err := o.f.Save(path.Join(o.path, fmt.Sprintf("output_%s.xlsx", time.Now().Format(time.RFC3339))))
	if err != nil {
		fmt.Println("save output file failed ", err)
		os.Exit(-1)
	}
}


