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
package httpE

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"golang.org/x/net/http2"
	"crypto/tls"
	"time"
	"net/url"
	"io"
	"os"
	"fmt"
	"strings"
	gourl "net/url"

	"github.com/heidawei/smartBoom/executor"
	"github.com/heidawei/smartBoom/register"

)

const MaxIdleConn = 2

var Name = "http"

type HttpE struct {
	cli    *http.Client
	request *http.Request
	config  map[string]interface{}
	body   []byte
}

func New(config map[string]interface{}) executor.Executor {
	var host string
	var proxyAddr *url.URL
	var disableCompression, disableKeepAlives, h2 bool
	var timeout int
	var err error
	if config != nil {
		if h, ok := config["host"]; ok {
			host = h.(string)
		}
		if p, ok := config["proxy"]; ok {
			proxyAddr, err = gourl.Parse(p.(string))
			if err != nil {
				os.Exit(-1)
			}
		}
		if d, ok := config["disableCompression"]; ok {
			disableCompression = d.(bool)
		}
		if d, ok := config["disableKeepAlives"]; ok {
			disableKeepAlives = d.(bool)
		}
		if h, ok := config["h2"]; ok {
			h2 = h.(bool)
		}
		if t, ok := config["timeout"]; ok {
			timeout = int(t.(float64))
		}
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		},
		MaxIdleConnsPerHost: MaxIdleConn,
		DisableCompression:  disableCompression,
		DisableKeepAlives:   disableKeepAlives,
		Proxy:               http.ProxyURL(proxyAddr),
	}
	if h2 {
		http2.ConfigureTransport(tr)
	} else {
		tr.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(timeout) * time.Second}
	return &HttpE{cli: client, config: config}
}

func(h *HttpE)Init() {
	var url, contentType, accept, method string
	var bodyAll []byte
	contentType = "text/html"
	method = "GET"
	if h.config != nil {
		if u, ok := h.config["url"]; ok {
			url = u.(string)
		} else {
			fmt.Println("ulr must not empty")
			os.Exit(-1)
		}
		if t, ok := h.config["Content-Type"]; ok {
			contentType = t.(string)
		}
		if m, ok := h.config["method"]; ok {
			method = strings.ToUpper(m.(string))
		}
		if a, ok := h.config["Accept"]; ok {
			accept = a.(string)
		}
		if b, ok := h.config["body"]; ok {
			bodyAll = []byte(b.(string))
			h.body = bodyAll
		}
	}

	// set content-type
	header := make(http.Header)
	header.Set("Content-Type", contentType)

	if accept != "" {
		header.Set("Accept", accept)
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		os.Exit(-1)
	}
	req.ContentLength = int64(len(bodyAll))
	req.Header = header
	h.request = req
	return
}

// return num of message do
func(h *HttpE)Do(base, index, n int) *executor.Result {
	s := now()
	var size int64
	var code int
	req := cloneRequest(h.request, h.body)
	resp, err := h.cli.Do(req)
	if err == nil {
		size = resp.ContentLength
		code = resp.StatusCode
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	t := now()
	finish := t - s
	return &executor.Result{
		StatusCode:    code,
		Duration:      finish,
		Err:           err,
		ContentLength: size,
		Count:         1,
	}
}

func cloneRequest(r *http.Request, body []byte) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	if len(body) > 0 {
		r2.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
	return r2
}

var startTime = time.Now()

// now returns time.Duration using stdlib time
func now() time.Duration { return time.Since(startTime) }

func init() {
	register.RegisterExecutor(Name, New)
}


