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

	return &OutPut{f: f, sheet: sheet, path: path}
}

func (o *OutPut) Write(f *Finalize) {
	r := o.sheet.AddRow()
	cell := r.AddCell()
	// timestamp
	cell.Value = f.TimeStamp
	cell = r.AddCell()
	// TPS
	cell.Value = fmt.Sprintf("%6.2f", f.TPS)
	cell = r.AddCell()
	// avg
	cell.Value = fmt.Sprintf("%6.5f", f.AvgDelay)
	cell = r.AddCell()
	// success
	cell.Value = fmt.Sprintf("%d", f.Success)
	cell = r.AddCell()
	// fail
	cell.Value = fmt.Sprintf("%d", f.Err)
	cell = r.AddCell()
	// TP10
	cell.Value = fmt.Sprintf("%g", f.TP10)
	cell = r.AddCell()
	// TP25
	cell.Value = fmt.Sprintf("%g", f.TP25)
	cell = r.AddCell()
	// TP50
	cell.Value = fmt.Sprintf("%g", f.TP50)
	cell = r.AddCell()
	// TP75
	cell.Value = fmt.Sprintf("%g", f.TP75)
	cell = r.AddCell()
	// TP90
	cell.Value = fmt.Sprintf("%g", f.TP90)
	cell = r.AddCell()
	// TP95
	cell.Value = fmt.Sprintf("%g", f.TP95)
	cell = r.AddCell()
	// TP99
	cell.Value = fmt.Sprintf("%g", f.TP99)
}

func (o *OutPut) Save() {
	err := o.f.Save(path.Join(o.path, fmt.Sprintf("output_%s.xlsx", time.Now().Format(time.RFC3339))))
	if err != nil {
		fmt.Println("save output file failed ", err)
	}
}


