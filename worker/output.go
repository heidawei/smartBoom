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
	options xlsx.DateTimeOptions
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
	l, _ := time.LoadLocation("Local")
	options := xlsx.DateTimeOptions{Location: l, ExcelTimeFormat: "h:mm:ss"}
	return &OutPut{f: f, sheet: sheet, path: path, options: options}
}

func (o *OutPut) Write(f *Finalize) {
	r := o.sheet.AddRow()

	// timestamp
	cell := r.AddCell()
	cell.SetDateWithOptions(f.TimeStamp, o.options)
	// TPS
	cell = r.AddCell()
	cell.SetFloat(f.TPS)
	// avg
	cell = r.AddCell()
	cell.SetFloat(f.AvgDelay)
	// success
	cell = r.AddCell()
	cell.SetInt64(f.Success)
	// fail
	cell = r.AddCell()
	cell.SetInt64(f.Err)
	// TP10
	cell = r.AddCell()
	cell.SetFloat(f.TP10)
	// TP25
	cell = r.AddCell()
	cell.SetFloat(f.TP25)
	// TP50
	cell = r.AddCell()
	cell.SetFloat(f.TP50)
	// TP75
	cell = r.AddCell()
	cell.SetFloat(f.TP75)
	// TP90
	cell = r.AddCell()
	cell.SetFloat(f.TP90)
	// TP95
	cell = r.AddCell()
	cell.SetFloat(f.TP95)
	// TP99
	cell = r.AddCell()
	cell.SetFloat(f.TP99)
	fmt.Printf("====>>TPS: %f, avgDelay: %fms, TP99: %fms\n", f.TPS, f.AvgDelay, f.TP99)
}

func (o *OutPut) Save() {
	err := o.f.Save(path.Join(o.path, fmt.Sprintf("output_%s.xlsx", time.Now().Format(time.RFC3339))))
	if err != nil {
		fmt.Println("save output file failed ", err)
		os.Exit(-1)
	}
}


