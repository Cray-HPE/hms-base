// MIT License
//
// (C) Copyright [2018-2021] Hewlett Packard Enterprise Development LP
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
	"reflect"
	"testing"
)

type TypeDecodePair struct {
	xnameStr string
	hmsType  HMSType
}

var goodTests = []TypeDecodePair{
	{"d0", CDU},
	{"d0w0", CDUMgmtSwitch},
	{"x0d0", CabinetCDU},
	{"x0m1p0", CabinetPDU},
	{"x2000m3", CabinetPDUController},
	{"x0m0i1", CabinetPDUNic},
	{"x0m1p3j1", CabinetPDUOutlet},
	{"x0m1p3v1", CabinetPDUPowerConnector},
	{"x0", Cabinet},
	{"x0b0", CabinetBMC},
	{"x0c0", Chassis},
	{"x0c0b0", ChassisBMC},
	{"x0c0b0i1", ChassisBMCNic},
	{"x0c0t9", CMMRectifier},
	{"x0c0f0", CMMFpga},
	{"x0e0", CEC},
	{"x0c0s0", ComputeModule},
	{"x0c0r0", RouterModule},
	{"x0c0s0b0", NodeBMC},
	{"x0c0s0e0", NodeEnclosure},
	{"x0c0s0e0t0", NodeEnclosurePowerSupply},
	{"x0c0s0j1", NodePowerConnector},
	{"x0c0s0v1", NodePowerConnector},
	{"x0c0s0b0n0", Node},
	{"x0c0s0b0n0p0", Processor},
	{"x0c0s0b0n0g0k0", Drive},
	{"x0c0s0b0n0g0", StorageGroup},
	{"x0c0s0b0n0i1", NodeNic},
	{"x0c0s0b0n0h1", NodeHsnNic},
	{"x0c0s0b0n0d0", Memory},
	{"x0c0s0b0n0a0", NodeAccel},
	{"x0c0s0b0n0r0", NodeAccelRiser},
	{"x0c0s0b0f0", NodeFpga},
	{"x0c0r16e0", HSNBoard},
	{"x0c0r0a0", HSNAsic},
	{"x0c0r0f0", RouterFpga},
	{"x0c0r0t0f0", RouterTORFpga},
	{"x0c0r0b0", RouterBMC},
	{"x0c0r0b0i1", RouterBMCNic},
	{"x0c0r0v1", RouterPowerConnector},
	{"x0c0r0a0l0", HSNLink},
	{"x0c0r0j1", HSNConnector},
	{"x0c0r0j1p1", HSNConnectorPort},
	{"x0c0w1", MgmtSwitch},
	{"x0c0w1j1", MgmtSwitchConnector},
	{"x0c0h1s2", MgmtHLSwitch},
	{"sms1", SMSBox},
	{"p0.0", Partition},
	{"s0", System},
	{"all", HMSTypeAll},
	{"all_comp", HMSTypeAllComp},
	{"all_svc", HMSTypeAllSvc},
}

// Get the HMSType for a given xname, based on its pattern in the recognition
// table above.
// If no string matches, HMSTypeInvalid is returned.
func TestGetHMSType(t *testing.T) {
	if gotType := GetHMSType("asdfasdfa"); gotType != HMSTypeInvalid {
		t.Errorf("Testcase a: FAIL: Expected %s, but got %s",
			HMSTypeInvalid.String(), gotType)
	} else {
		t.Logf("Testcase a: PASS: Expected %s and got %s",
			HMSTypeInvalid.String(), gotType)
	}
	if gotType := GetHMSType(""); gotType != HMSTypeInvalid {
		t.Errorf("Testcase b: FAIL: Expected %s, but got %s",
			HMSTypeInvalid.String(), gotType.String())
	} else {
		t.Logf("Testcase b: PASS: Expected %s and got %s",
			HMSTypeInvalid.String(), gotType.String())
	}
	for i, pair := range goodTests {
		gotType := GetHMSType(pair.xnameStr)
		if gotType != pair.hmsType {
			t.Errorf("Testcase %d (%s): FAIL: Expected %s, but got %s",
				i, pair.xnameStr, pair.hmsType, gotType)
		} else {
			t.Logf("Testcase %d (%s): PASS: Expected %s and got %s",
				i, pair.xnameStr, pair.hmsType, gotType)
		}
	}
}

// Returns string value rather than HMSType for xname
func TestGetHMSTypeString(t *testing.T) {
	if gotStr := GetHMSTypeString("asdfasdfa"); gotStr != "" {
		t.Errorf("Testcase a: FAIL: Expected empty string, got %s", gotStr)
	} else {
		t.Logf("Testcase a: PASS: Got empty string")
	}
	if gotStr := GetHMSTypeString(""); gotStr != "" {
		t.Errorf("Testcase b: FAIL: Expected empty string, got %s", gotStr)
	} else {
		t.Logf("Testcase b: PASS: Got empty string")
	}
	for i, pair := range goodTests {
		gotStr := GetHMSTypeString(pair.xnameStr)
		if gotStr != pair.hmsType.String() {
			t.Errorf("Testcase %d (%s): FAIL: Expected %s, but got %s",
				i, pair.xnameStr, pair.hmsType, gotStr)
		} else {
			t.Logf("Testcase %d (%s): PASS: Expected %s and got %s",
				i, pair.xnameStr, pair.hmsType, gotStr)
		}
	}
}

type ValidateTestData struct {
	input []string
	good  []string
	bad   []string
	dups  bool
}

var compIDValidateTD = []ValidateTestData{
	{
		// Test Case 0
		[]string{"x0", "x1", "x2", "x3", "x0", "x10", "x100", "x1000"},
		[]string{"x0", "x1", "x2", "x3", "x10", "x100", "x1000"},
		[]string{"x0"},
		false,
	},
	{
		// Test Case 1
		[]string{"x0", "x1", "x2", "x3", "x0", "x10", "x100", "x1000"},
		[]string{"x0", "x1", "x2", "x3", "x10", "x100", "x1000"},
		[]string{},
		true,
	},
	{
		// Test Case 2
		[]string{"x0", "xray", "d0", "down"},
		[]string{"x0", "d0"},
		[]string{"xray", "down"},
		false,
	},
	{
		// Test Case 3
		[]string{"x0d0", "x0b0", "x0c0", "x0e0", "x0m0"},
		[]string{"x0d0", "x0b0", "x0c0", "x0e0", "x0m0"},
		[]string{},
		false,
	},
	{
		// Test Case 4
		[]string{"x0m0", "x0m0p0", "x0m0i0", "x0m0p0j0", "x0m0p0j1", "x1c0"},
		[]string{"x0m0", "x0m0p0", "x0m0i0", "x0m0p0j1", "x1c0"},
		[]string{"x0m0p0j0"},
		false,
	},
	{
		// Test Case 5
		[]string{"x0c0b0", "x0c0t0", "x0c0", "x0c0f0", "x0c0s0", "x0c0r0", "x1c0w0", "x1c0w1", "x1c0h1s2", "x1c0h0s2", "x1c0h1s0"},
		[]string{"x0c0b0", "x0c0t0", "x0c0", "x0c0f0", "x0c0s0", "x0c0r0", "x1c0w1", "x1c0h1s2"},
		[]string{"x1c0w0", "x1c0h0s2", "x1c0h1s0"},
		false,
	},
	{
		// Test Case 6
		[]string{"x1c0r47j31", "x0c0r1a0", "x0c0r1a0l64", "x1c0r2e0", "x1c0r2b0", "x1c2r3b0i0"},
		[]string{"x1c0r47j31", "x0c0r1a0", "x0c0r1a0l64", "x1c0r2e0", "x1c0r2b0", "x1c2r3b0i0"},
		[]string{},
		false,
	},
	{
		// Test Case 7
		[]string{"x16c3r2f0", "x1234c0r1t0f0", "x16c3r4j7", "x16c3r4j7p0"},
		[]string{"x16c3r2f0", "x1234c0r1t0f0", "x16c3r4j7", "x16c3r4j7p0"},
		[]string{},
		false,
	},
	{
		// Test Case 8
		[]string{"x0c0s0", "x0c0s64", "x0c0s4j1", "x0c0s4j4"},
		[]string{"x0c0s0", "x0c0s64", "x0c0s4j1"},
		[]string{"x0c0s4j4"},
		false,
	},
	{
		// Test Case 9
		[]string{"x16c3s1b1", "x16c2s3b0n1", "x3c0s16e0", "x16c2s3b0n1p1", "x16c2s3b0n1g1", "x16c2s3b0n1g1k1", "x16c3s0b0n1h1", "x16c3s0b0n1i1", "x16c3s0b0n1d3", "x16c3s0b1n0a1", "x16c3s4b1f0"},
		[]string{"x16c3s1b1", "x16c2s3b0n1", "x3c0s16e0", "x16c2s3b0n1p1", "x16c2s3b0n1g1", "x16c2s3b0n1g1k1", "x16c3s0b0n1h1", "x16c3s0b0n1i1", "x16c3s0b0n1d3", "x16c3s0b1n0a1", "x16c3s4b1f0"},
		[]string{},
		false,
	},
	{
		// Test Case 10
		[]string{"x1m0p0j0", "x1m0p0j1", "x1m1p0j64", "x1m0p0j1", "x1m0p1j10"},
		[]string{"x1m0p0j1", "x1m1p0j64", "x1m0p1j10"},
		[]string{"x1m0p0j0"},
		true,
	},
	{
		// Test case 11 - CDU Mgmt Switches
		[]string{"d0w0", "d0w1", "d0w30", "d0wW"},
		[]string{"d0w0", "d0w1", "d0w30"},
		[]string{"d0wW"},
		false,
	},
}

// TestValidateCompIDs unit test function for ValidateCompIDs
func TestValidateCompIDs(t *testing.T) {

	for n, test := range compIDValidateTD {
		good, bad := ValidateCompIDs(test.input, test.dups)
		if !reflect.DeepEqual(good, test.good) {
			t.Errorf("TestValidateCompIDs Test Case %d: FAIL: Expected valid=%v and invalid=%v but instead got valid=%v and invalid=%v",
				n, test.good, test.bad, good, bad)
		}
		if !reflect.DeepEqual(bad, test.bad) {
			t.Errorf("TestValidateCompIDs Test Case %d: FAIL: Expected valid=%v and invalid=%v but instead got valid=%v and invalid=%v",
				n, test.good, test.bad, good, bad)
		}
	}
}
