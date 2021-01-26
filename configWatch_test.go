// MIT License
//
// (C) Copyright [2020-2021] Hewlett Packard Enterprise Development LP
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
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadFile(t *testing.T) {
	tests := []struct {
		config       []byte
		roles        []string
		subroles     []string
	}{{ // Test 0: Normal case
		json.RawMessage(`{"HMSExtendedDefinitions":{"Role":["UAN","Foo"],"SubRole":["Data","Bar"]}}`),
		[]string{"UAN", "Foo"},
		[]string{"Data", "Bar"},
	}, { // Test 1: No extra SubRoles
		json.RawMessage(`{"HMSExtendedDefinitions":{"Role":["UAN","Foo"]}}`),
		[]string{"UAN", "Foo"},
		[]string{},
	}, { // Test 2: No extra Roles
		json.RawMessage(`{"HMSExtendedDefinitions":{"SubRole":["Data","Bar"]}}`),
		[]string{},
		[]string{"Data", "Bar"},
	}, { // Test 3: Extra parameters
		json.RawMessage(`{"HMSExtendedDefinitions":{"Role":["UAN","Foo"],"SubRole":["Data","Bar"],"Hello":"World"}}`),
		[]string{"UAN", "Foo"},
		[]string{"Data", "Bar"},
	}, { // Test 4: Bad json
		json.RawMessage(`"HMSExtendedDefinitions":{"Role":["UAN","Foo"],"SubRole":["Data","Bar"]}}`),
		[]string{},
		[]string{},
	}}
	fn := "test_config.json"
	for i, test := range tests {
		hmsRoleMap = defaultHMSRoleMap
		hmsSubRoleMap = defaultHMSSubRoleMap
		validRoles := GetHMSRoleList()
		if len(test.roles) > 0 {
			validRoles = append(validRoles, test.roles...)
		}
		validSubRoles := GetHMSSubRoleList()
		if len(test.subroles) > 0 {
			validSubRoles = append(validSubRoles, test.subroles...)
		}
		err := ioutil.WriteFile(fn, test.config, 0644)
		if err != nil {
			t.Errorf("Test %v Failed: Failed to write test file - %s", i, err)
		}
		loadFile(fn)
		err = os.Remove(fn)
		if err != nil {
			t.Errorf("Test %v Failed: Failed to remove test file - %s", i, err)
		}
		for _, role := range validRoles {
			if VerifyNormalizeRole(role) == "" {
				t.Errorf("Test %v Failed: Invalid role - %s", i, role)
			}
		}
		for _, subrole := range validSubRoles {
			if VerifyNormalizeSubRole(subrole) == "" {
				t.Errorf("Test %v Failed: Invalid subrole - %s", i, subrole)
			}
		}
	}
}
