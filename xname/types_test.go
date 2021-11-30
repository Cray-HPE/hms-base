package xname

import (
	"fmt"
	"strconv"
	"testing"

	base "github.com/Cray-HPE/hms-base"
)

func TestFoo(t *testing.T) {
	n := Node{
		Cabinet: 1000,
		Chassis: 1,
		Slot:    7,
		BMC:     1,
		Node:    0,
	}

	t.Log("Node:", n)
	t.Log("NodeBMC:", n.Parent())
	t.Log("NodeModule:", n.Parent().Parent())
	t.Log("Chassis:", n.Parent().Parent().Parent())
	t.Log("Cabinet:", n.Parent().Parent().Parent().Parent())
	t.Log("System:", n.Parent().Parent().Parent().Parent().Parent())

	n = Cabinet{Cabinet: 1000}.Chassis(1).NodeBMC(7, 1).Node(0)
	t.Log("Node:", n)
	n = Cabinet{Cabinet: 1000}.Chassis(1).ComputeModule(7).NodeBMC(1).Node(0)
	t.Log("Node:", n)
	n = System{}.
		Cabinet(1000).
		Chassis(1).
		ComputeModule(7).
		NodeBMC(1).
		Node(0)
	t.Log("Node:", n)

	n = System{}.
		Cabinet(1000).
		Chassis(1).
		NodeBMC(7, 1).
		Node(0)
	t.Log("Node:", n)

	n = Node{
		Cabinet: 1000,
		Chassis: 1,
		Slot:    7,
		BMC:     1,
		Node:    0,
	}

	hmsType, err := GetHMSType(n)
	if err != nil {
		t.Log("GetHMSType error:", err)
		t.FailNow()
		return
	}
	t.Log("HMS Type:", hmsType)

	formatStr, numArgs, err := base.GetHMSTypeFormatString(hmsType)
	if err != nil {
		t.Log("GetHMSTypeFormatString error:", err)
		t.FailNow()
		return
	}
	t.Log("Format String args:", numArgs)
	t.Log("Format String:", formatStr)

	cduSwitch := System{}.CDU(0).CDUMgmtSwitch(1)
	t.Log("CDU Switch:", cduSwitch)

	ms := MgmtSwitch{
		Cabinet: 1,  // X: 0-999
		Chassis: 0,  // C: 0-7
		Slot:    32, // W: 1-48
	}
	t.Log("MgmtSwitch:", ms)

}

func TestRegex(t *testing.T) {
	xname := "x1c2s3b4n5"
	hmsType := base.GetHMSType(xname)
	t.Log("HMS Type:", hmsType)

	re, err := base.GetHMSTypeRegex(hmsType)
	if err != nil {
		t.Fatal("GetHMSTypeRegex error", err)
		return
	}

	_, argCount, err := base.GetHMSTypeFormatString(hmsType)
	if err != nil {
		t.Fatal("GetHMSTypeFormatString error", err)
		return
	}
	t.Log("Format String Args:", argCount)

	matchesRaw := re.FindStringSubmatch(xname)
	t.Log("Matches Raw", matchesRaw)

	if (argCount + 1) != len(matchesRaw) {
		t.Fatal("Unexpected number of matches found:", len(matchesRaw), "expected:", argCount)
		return
	}

	matches := []int{}
	for _, matchRaw := range matchesRaw[1:] {
		// If we have gotten to this point these matches should be integers
		match, err := strconv.Atoi(matchRaw)
		if err != nil {
			t.Fatal("unable to convert match to integer:", matchRaw, "error:", err)
			return
		}

		matches = append(matches, match)
	}

	t.Log("Matches", matches)

	node := Node{
		Cabinet: matches[0],
		Chassis: matches[1],
		Slot:    matches[2],
		BMC:     matches[3],
		Node:    matches[4],
	}

	t.Log("Node", node)

}

func TestToFromXnames(t *testing.T) {
	// Note, not all of the xnames in the following tests are valid. Each ordinal is incremented by 1 to verify that each ordinal is being properly
	// handled and not getting lost or switched around.
	tests := []struct {
		xname             string
		hmsType           base.HMSType
		expectedComponent interface{}
	}{
		{
			"s0",
			base.System,
			System{},
		}, {
			"d0",
			base.CDU,
			CDU{
				CoolingGroup: 0,
			},
		}, {
			"d0w1", base.CDUMgmtSwitch,
			CDUMgmtSwitch{
				CoolingGroup: 0,
				Slot:         1,
			},
		}, {
			"x1",
			base.Cabinet,
			Cabinet{
				Cabinet: 1,
			},
		}, {
			"x1c2",
			base.Chassis,
			Chassis{
				Cabinet: 1,
				Chassis: 2,
			},
		}, {
			"x1c2b0", // TODO add a test to verify what happens when x1c2b3 is given.
			base.ChassisBMC,
			ChassisBMC{
				Cabinet: 1,
				Chassis: 2,
				BMC:     0,
			},
		}, {
			"x1c2h3",
			base.MgmtHLSwitchEnclosure,
			MgmtHLSwitchEnclosure{
				Cabinet: 1,
				Chassis: 2,
				Slot:    3,
			},
		}, {
			"x1c2h3s4",
			base.MgmtHLSwitch,
			MgmtHLSwitch{
				Cabinet: 1,
				Chassis: 2,
				Slot:    3,
				Space:   4,
			},
		}, {
			"x1c2w3",
			base.MgmtSwitch,
			MgmtSwitch{
				Cabinet: 1,
				Chassis: 2,
				Slot:    3,
			},
		}, {
			"x1c2w3j4",
			base.MgmtSwitchConnector,
			MgmtSwitchConnector{
				Cabinet:    1,
				Chassis:    2,
				Slot:       3,
				SwitchPort: 4,
			},
		}, {
			"x1c2r3",
			base.RouterModule,
			RouterModule{
				Cabinet: 1,
				Chassis: 2,
				Slot:    3,
			},
		}, {
			"x1c2r3b4",
			base.RouterBMC,
			RouterBMC{
				Cabinet: 1,
				Chassis: 2,
				Slot:    3,
				BMC:     4,
			},
		}, {
			"x1c2s3",
			base.ComputeModule,
			ComputeModule{
				Cabinet: 1,
				Chassis: 2,
				Slot:    3,
			},
		}, {
			"x1c2s3b4",
			base.NodeBMC,
			NodeBMC{
				Cabinet: 1,
				Chassis: 2,
				Slot:    3,
				BMC:     4,
			},
		}, {
			"x1c2s3b4n5",
			base.Node,
			Node{
				Cabinet: 1,
				Chassis: 2,
				Slot:    3,
				BMC:     4,
				Node:    5,
			},
		},
	}

	for _, test := range tests {
		xname := test.xname
		expectedHMSType := test.hmsType

		// Just a sanity check to verify that out test data is good
		if hmsType := base.GetHMSType(xname); hmsType != expectedHMSType {
			t.Errorf("unexpected HMS Type (%s) for xname (%s) in test data, expected (%s)", hmsType, xname, expectedHMSType)
		}

		// Verify FromString returns the HMS Type
		componentRaw, hmsType := FromString(xname)
		if expectedHMSType != hmsType {
			t.Error("Unexpected HMS Type:", hmsType, "expected:", expectedHMSType)
		}

		// Verify FromString returns the correct xname struct values
		if componentRaw != test.expectedComponent {
			t.Errorf("Unexpected xname struct (%v), expected (%v)", componentRaw, test.expectedComponent)
		}

		// Verify that GetHMSType works
		objXnameType, err := GetHMSType(componentRaw)
		if err != nil {
			t.Error("GetHMSType error:", err)
		}
		if expectedHMSType != objXnameType {
			t.Error("Unexpected HMS Type for xname struct:", objXnameType, "expected:", expectedHMSType)
		}

		// Verify the xname string built from the xname struct matches what was given to FromString
		generatedXname := componentRaw.(fmt.Stringer).String()
		if xname != generatedXname {
			t.Error("Unexpected generated xname:", generatedXname, "expected:", xname)
		}

		// Verify the HMS Type of the xname built FromString has the expected HMS Type
		generatedXnameType := base.GetHMSType(generatedXname)
		if expectedHMSType != generatedXnameType {
			t.Errorf("Unexpected generated xname %s (%s), expected (%s) %s", generatedXnameType, generatedXname, expectedHMSType, xname)
		}
	}
}
