// MIT License
// 
// (C) Copyright [2021] Hewlett Packard Enterprise Development LP
// 
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package base

import (
	"net/http"
	"testing"
	"os"
)

func TestGetServiceInstanceName(t *testing.T) {
	hname,err := os.Hostname()
	if (err != nil) {
		t.Logf("Warning, can't get os.Hostname()!")
	}

	inst,ierr := GetServiceInstanceName()
	if ((err == nil) && (ierr != nil)) {
		t.Errorf("Hostname() call worked, GetServiceInstanceName() failed: %v",
			ierr)
	}

	if ((err != nil) && (ierr == nil)) {
		t.Errorf("Hostname() call failed, GetServiceInstanceName() worked: %v",
			err)
	}

	if (hname != inst) {
		t.Errorf("Hostname mismatch: hostname: '%s' instname: '%s'",hname,inst)
	}
}

func Test_SetHTTPUserAgent(t *testing.T) {
	var hkey string
	expval := "xyzzy"

	req,_ := http.NewRequest("GET","http://alfred_e_newman.com",nil)
	SetHTTPUserAgent(req,expval)
	hkey = req.Header.Get(USERAGENT)
	if (hkey == "") {
		t.Errorf("%s key not present!",USERAGENT)
	}
	if (hkey != expval) {
		t.Errorf("%s key has wrong value, expected: '%s', got: '%s'",
			USERAGENT,expval,hkey)
	}

	req,_ = http.NewRequest("POST","http://what_me_worry.com",nil)
	req.Header.Set("Content-Type","application/json")
	SetHTTPUserAgent(req,expval)
	hkey = req.Header.Get(USERAGENT)
	if (hkey == "") {
		t.Errorf("%s key not present!",USERAGENT)
	}
	if (hkey != expval) {
		t.Errorf("%s key has wrong value, expected: '%s', got: '%s'",
			USERAGENT,expval,hkey)
	}
}


