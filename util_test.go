// Copyright (c) 2019 Cray Inc. All Rights Reserved.
package base

import (
	"testing"
)

func TestIsAlphaNum(t *testing.T) {
	var argsGood = []string{
		"x0c0b0",
		"X0c0s0b000",
		"0A0",
		"1",
		"P",
	}
	var argsBad = []string{
		"x0c0s0b0 b0",    // space
		"A0:K0",          // Punctuation
		" x0c0 ",         // Leading whitespace
		"\"hellothere\"", // Quotes
	}
	for i, arg := range argsGood {
		if IsAlphaNum(arg) == false {
			t.Errorf("Testcase %da: FAIL Got unexpected 'true' for '%s'",
				i, arg)
		} else {
			t.Logf("Testcase %da: PASS Got 'true' for '%s'", i, arg)
		}
	}
	for i, arg := range argsBad {
		if IsAlphaNum(arg) == false {
			t.Logf("Testcase %db: Pass Got expected 'false' for '%s'",
				i, arg)
		} else {
			t.Errorf("Testcase %db: Fail Got 'true' for '%s'", i, arg)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	var argsGood = []string{
		"1000",
		"010",
		"000",
		"1",
		"2",
		"0",
	}
	var argsBad = []string{
		"0000v",    // Letters
		"0 0",      // space
		"0:0",      // Punctuation
		" 0 ",      // Leading whitespace
		"\"1234\"", // Quotes
	}
	for i, arg := range argsGood {
		if IsNumeric(arg) == false {
			t.Errorf("Testcase %da: FAIL Got unexpected 'true' for '%s'",
				i, arg)
		} else {
			t.Logf("Testcase %da: PASS Got 'true' for '%s'", i, arg)
		}
	}
	for i, arg := range argsBad {
		if IsNumeric(arg) == false {
			t.Logf("Testcase %db: Pass Got expected 'false' for '%s'",
				i, arg)
		} else {
			t.Errorf("Testcase %db: Fail Got 'true' for '%s'", i, arg)
		}
	}
}

func TestRemoveLeadingZeros(t *testing.T) {
	var inputs = []string{
		"x0c0s0b0",
		"x0000c00s00b00",
		"x01000c01s010b01",
		"x10000c10s100b1000",
		"x0c0s0b01",
		"x0c0s0b10",
		"x0c0s0b1",
		"00",
		"0",
		"a",
		"1",
	}
	var outputs = []string{
		"x0c0s0b0",
		"x0c0s0b0",
		"x1000c1s10b1",
		"x10000c10s100b1000",
		"x0c0s0b1",
		"x0c0s0b10",
		"x0c0s0b1",
		"0",
		"0",
		"a",
		"1",
	}
	for i := 0; i < len(inputs); i++ {
		if RemoveLeadingZeros(inputs[i]) != outputs[i] {
			t.Errorf("Testcase %da: FAIL Got unexpected '%s' vs '%s' for: '%s'",
				i, RemoveLeadingZeros(inputs[i]), outputs[i], inputs[i])
		} else {
			t.Logf("Testcase %da: PASS Got '%s' for: '%s'",
				i, outputs[i], inputs[i])
		}
	}
}
