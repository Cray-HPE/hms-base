// Copyright 2020 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.

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
