package xname

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	base "github.com/Cray-HPE/hms-base"
	"github.com/hashicorp/go-multierror"
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
			"x1c2b0",
			base.ChassisBMC,
			ChassisBMC{
				Cabinet: 1,
				Chassis: 2,
				BMC:     0,
			},
		// }, { // TODO This causes a panic
		// 	"x1c2b3",
		// 	base.ChassisBMC,
		// 	ChassisBMC{
		// 		Cabinet: 1,
		// 		Chassis: 2,
		// 		BMC:     3,
		// 	},
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

//
//
// Tests to verify that Parent/Children functions behave as expected
//
//

func TestSystemChildren(t *testing.T) {
	system := System{}

	// Create a child CDU
	cdu := system.CDU(1)
	expectedCDU := CDU{
		CoolingGroup: 1,
	}
	if !reflect.DeepEqual(expectedCDU, cdu) {
		t.Errorf("TestSystemChildren FAIL: Expected cdu=%v but instead got cdu=%v", expectedCDU, cdu)
	}

	// Create a child cabinet
	cabinet := system.Cabinet(1)
	expectedCabinet := Cabinet{
		Cabinet: 1,
	}
	if !reflect.DeepEqual(expectedCabinet, cabinet) {
		t.Errorf("TestSystemChildren FAIL: Expected cabinet=%v but instead got cabinet=%v", expectedCabinet, cabinet)
	}
}

func TestSystemParent(t *testing.T) {
	// A System doesn't have a parent
}

func TestCDUChildren(t *testing.T) {
	cdu := CDU{
		CoolingGroup: 1,
	}

	// Create a child CDUMgmtSwitch
	cduMgmtSwitch := cdu.CDUMgmtSwitch(2)
	expectedCDUMgmtSwitch := CDUMgmtSwitch{
		CoolingGroup: 1,
		Slot: 2,
	}
	if !reflect.DeepEqual(expectedCDUMgmtSwitch, cduMgmtSwitch) {
		t.Errorf("TestCDUChildren FAIL: Expected cduMgmtSwitch=%v but instead got cduMgmtSwitch=%v", expectedCDUMgmtSwitch, cduMgmtSwitch)
	}
}

func TestCDUParent(t *testing.T) {
	cdu := CDU{
		CoolingGroup: 1,
	}

	parent := cdu.Parent()
	expectedParent := System{}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestCDUParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestCabinetChildren(t *testing.T) {
	cabinet := Cabinet{
		Cabinet: 1,
	}

	// Create a child CabinetPDUController
	cabinetPDUController := cabinet.CabinetPDUController(2)
	expectedCabinetPDUController := CabinetPDUController{
		Cabinet: 1,
		PDUController: 2,
	}
	if !reflect.DeepEqual(expectedCabinetPDUController, cabinetPDUController) {
		t.Errorf("TestCabinetChildren FAIL: Expected cabinetPDUController=%v but instead got cabinetPDUController=%v", expectedCabinetPDUController, cabinetPDUController)
	}

	// Create a child Chassis
	chassis := cabinet.Chassis(2)
	expectedChassis := Chassis{
		Cabinet: 1,
		Chassis: 2,
	}
	if !reflect.DeepEqual(expectedChassis, chassis) {
		t.Errorf("TestCabinetChildren FAIL: Expected chassis=%v but instead got chassis=%v", expectedChassis, chassis)
	}
}

func TestCabinetParent(t *testing.T) {
	cabinet := Cabinet{
		Cabinet: 1,
	}
	
	parent := cabinet.Parent()
	expectedParent := System{}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestCabinetParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestCabinetPDUControllerChildren(t *testing.T) {
	// TODO no children structures have bene defined yet, but child xname formats have been defined
}

func TestCabinetPDUControllerParent(t *testing.T) {
	cabinetPDUController := CabinetPDUController{
		Cabinet: 1,
		PDUController: 2,
	}
	
	parent := cabinetPDUController.Parent()
	expectedParent := Cabinet{
		Cabinet: 1,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestCabinetPDUControllerParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestChassisChildren(t *testing.T) {
	chassis := Chassis{
		Cabinet: 1,
		Chassis: 2,
	}

	// Create a child ComputeModule
	computeModule := chassis.ComputeModule(3)
	expectedComputeModule := ComputeModule{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	if !reflect.DeepEqual(expectedComputeModule, computeModule) {
		t.Errorf("TestChassisChildren FAIL: Expected computeModule=%v but instead got computeModule=%v", expectedComputeModule, computeModule)
	}

	// Create a child NodeBMC
	nodeBMC := chassis.NodeBMC(3, 4)
	expectedNodeBMC := NodeBMC{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
	}
	if !reflect.DeepEqual(expectedNodeBMC, nodeBMC) {
		t.Errorf("TestChassisChildren FAIL: Expected nodeBMC=%v but instead got nodeBMC=%v", expectedNodeBMC, nodeBMC)
	}

	// Create a child MgmtSwitch
	mgmtSwitch := chassis.MgmtSwitch(3)
	expectedMgmtSwitch := MgmtSwitch{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	if !reflect.DeepEqual(expectedMgmtSwitch, mgmtSwitch) {
		t.Errorf("TestChassisChildren FAIL: Expected mgmtSwitch=%v but instead got mgmtSwitch=%v", expectedMgmtSwitch, mgmtSwitch)
	}

	// Create a child MgmtHLSwitchEnclosure
	mgmtHLSwitchEnclosure := chassis.MgmtHLSwitchEnclosure(3)
	expectedMgmtHLSwitchEnclosure := MgmtHLSwitchEnclosure{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	if !reflect.DeepEqual(expectedMgmtHLSwitchEnclosure, mgmtHLSwitchEnclosure) {
		t.Errorf("TestChassisChildren FAIL: Expected mgmtHLSwitchEnclosure=%v but instead got mgmtHLSwitchEnclosure=%v", expectedMgmtHLSwitchEnclosure, mgmtHLSwitchEnclosure)
	}

	// Create a child MgmtHLSwitch
	mgmtHLSwitch := chassis.MgmtHLSwitch(3, 4)
	expectedMgmtHLSwitch := MgmtHLSwitch{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		Space: 4,
	}
	if !reflect.DeepEqual(expectedMgmtHLSwitch, mgmtHLSwitch) {
		t.Errorf("TestChassisChildren FAIL: Expected mgmtHLSwitch=%v but instead got mgmtHLSwitch=%v", expectedMgmtHLSwitch, mgmtHLSwitch)
	}

	// Create a child RouterModule
	routerModule := chassis.RouterModule(3)
	expetedRouterModule := RouterModule{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	if !reflect.DeepEqual(expetedRouterModule, routerModule) {
		t.Errorf("TestChassisChildren FAIL: Expected routerModule=%v but instead got routerModule=%v", expetedRouterModule, routerModule)
	}


	// Create a child RouterBMC
	routerBMC := chassis.RouterBMC(3, 4)
	expectedRouterBMC := RouterBMC{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
	}
	if !reflect.DeepEqual(expectedRouterBMC, routerBMC) {
		t.Errorf("TestChassisChildren FAIL: Expected routerBMC=%v but instead got routerBMC=%v", expectedRouterBMC, routerBMC)
	}

}

func TestChassisParent(t *testing.T) {
	chassis := Chassis{
		Cabinet: 1,
		Chassis: 2,
	}
	
	parent := chassis.Parent()
	expectedParent := Cabinet{
		Cabinet: 1,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestChassisParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestChassisBMCChildren(t *testing.T) {
	// TODO no children structures have bene defined yet, but child xname formats have been defined
}

func TestChassisBMCParent(t *testing.T) {
	chassisBMC := ChassisBMC{
		Cabinet: 1,
		Chassis: 2,
		BMC: 0,
	}
	
	parent := chassisBMC.Parent()
	expectedParent := Chassis{
		Cabinet: 1,
		Chassis: 2,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestChassisBMCParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestMgmtSwitchChildren(t *testing.T) {
	mgmtSwitch := MgmtSwitch{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}

	// Create a child MgmtSwitchConnector
	mgmtSwitchConnector := mgmtSwitch.MgmtSwitchConnector(4)
	expectedMgmtSwitchConnector := MgmtSwitchConnector{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		SwitchPort: 4,
	}
	if !reflect.DeepEqual(expectedMgmtSwitchConnector, mgmtSwitchConnector) {
		t.Errorf("TestMgmtSwitchChildren FAIL: Expected mgmtSwitchConnector=%v but instead got mgmtSwitchConnector=%v", expectedMgmtSwitchConnector, mgmtSwitchConnector)
	}
}

func TestMgmtSwitchParent(t *testing.T) {
	mgmtSwitch := MgmtSwitch{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	
	parent := mgmtSwitch.Parent()
	expectedParent := Chassis{
		Cabinet: 1,
		Chassis: 2,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestMgmtSwitchParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestMgmtSwitchConnectorChildren(t *testing.T) {
	// There are no childlen for a MgmtSwitchConnector
}

func TestMgmtSwitchConnectorParent(t *testing.T) {
	mgmtSwitchConnector := MgmtSwitchConnector{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		SwitchPort: 4,
	}
	
	parent := mgmtSwitchConnector.Parent()
	expectedParent := MgmtSwitch{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestMgmtSwitchConnectorParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestMgmtHLSwitchEnclosureChildren(t *testing.T) {
	mgmtHLSwitchEnclosure := MgmtHLSwitchEnclosure{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}

	// Create a child MgmtHLSwitch
	mgmtHLSwitch := mgmtHLSwitchEnclosure.MgmtHLSwitch(4)
	expectedMgmtHLSwitch := MgmtHLSwitch{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		Space: 4,
	}
	if !reflect.DeepEqual(expectedMgmtHLSwitch, mgmtHLSwitch) {
		t.Errorf("TestMgmtHLSwitchEnclosureChildren FAIL: Expected mgmtHLSwitch=%v but instead got mgmtHLSwitch=%v", expectedMgmtHLSwitch, mgmtHLSwitch)
	}
}

func TestMgmtHLSwitchEnclosureParent(t *testing.T) {
	mgmtHLSwitchEnclosure := MgmtHLSwitchEnclosure{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	
	parent := mgmtHLSwitchEnclosure.Parent()
	expectedParent := Chassis{
		Cabinet: 1,
		Chassis: 2,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestMgmtHLSwitchEnclosureParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestMgmtHLSwitchChildren(t *testing.T) {
	// TODO no children structures have bene defined yet, and currently no child xname formats have been defined
}

func TestMgmtHLSwitchParent(t *testing.T) {
	mgmtHLSwitch := MgmtHLSwitch{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		Space: 4,
	}
	
	parent := mgmtHLSwitch.Parent()
	expectedParent := MgmtHLSwitchEnclosure{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestMgmtHLSwitchParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestRouterModuleChildren(t *testing.T) {
	routerModule := RouterModule{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}

	// Create a child RouterBMC
	routerBMC := routerModule.RouterBMC(4)
	expectedRouterBMC := RouterBMC{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
	}
	if !reflect.DeepEqual(expectedRouterBMC, routerBMC) {
		t.Errorf("TestRouterModuleChildren FAIL: Expected routerBMC=%v but instead got routerBMC=%v", expectedRouterBMC, routerBMC)
	}
}

func TestRouterModuleParent(t *testing.T) {
	routerModule := RouterModule{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	
	parent := routerModule.Parent()
	expectedParent := Chassis{
		Cabinet: 1,
		Chassis: 2,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestRouterModuleParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestRouterBMCChildren(t *testing.T) {
	// TODO no children structures have bene defined yet, but child xname formats have been defined
}

func TestRouterBMCParent(t *testing.T) {
	routerModule := RouterBMC{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
	}
	
	parent := routerModule.Parent()
	expectedParent := RouterModule{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestRouterBMCParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestComputeModuleChildren(t *testing.T) {
	computeModule := ComputeModule{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}

	// Create a child NodeBMC
	nodeBMC := computeModule.NodeBMC(4)
	expectedNodeBMC := NodeBMC{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
	}

	if !reflect.DeepEqual(expectedNodeBMC, nodeBMC) {
		t.Errorf("TestComputeModuleChildren FAIL: Expected nodeBMC=%v but instead got nodeBMC=%v", expectedNodeBMC, nodeBMC)
	}
}

func TestComputeModuleParent(t *testing.T) {
	computeModule := ComputeModule{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	
	parent := computeModule.Parent()
	expectedParent := Chassis{
		Cabinet: 1,
		Chassis: 2,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestComputeModuleParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestNodeBMCChildren(t *testing.T) {
	nodeBMC := NodeBMC{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
	}

	// Create a child Node
	node := nodeBMC.Node(0)
	expectedNode := Node{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
		Node: 0,
	}
	if !reflect.DeepEqual(expectedNode, node) {
		t.Errorf("TestNodeBMCChildren FAIL: Expected node=%v but instead got node=%v", expectedNode, node)
	}
}

func TestNodeBMCParent(t *testing.T) {
	nodeBMC := NodeBMC{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
	}
	
	parent := nodeBMC.Parent()
	expectedParent := ComputeModule{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestNodeBMCParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

func TestNodeChildren(t *testing.T) {
	// TODO no children structures have bene defined yet, but child xname formats have been defined
}

func TestNodeParent(t *testing.T) {
	node := Node{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
		Node: 0,
	}
	
	parent := node.Parent()
	expectedParent := NodeBMC{
		Cabinet: 1,
		Chassis: 2,
		Slot: 3,
		BMC: 4,
	}
	if !reflect.DeepEqual(expectedParent, parent) {
		t.Errorf("TestNodeParent FAIL: Expected parent=%v but instead got parent=%v", expectedParent, parent)
	}
}

//
//
// Validation function testing
//
//

func TestValidate(t *testing.T) {
	// CDU
	// Negative tests
	// - cooling group number is negative
	// Positive Tests
	// - cooling group is 0
	// - cooling group is 999
	// - cooling group is 123

	// CDUMgmtSwitch
	// Negative tests
	// - cooling group number is negative
	// - slot is negative
	// Positive tests
	// - slot is 0
	// - slot is 15
	// - slot is 31

	// Cabinet
	// Negative tests
	// - cabinet number is negative
	// - cabinet number is greater then 999
	// Positive tests
	// - cabinet is 0
	// - cabinet is 10
	// - cabinet is 999
}


func TestCDUValidationEnhanced(t *testing.T) {
	// CDU
	// Negative tests
	// - cooling group number is negative
	// - cooling group number is greater then 999
	// Positive Tests
	// - cooling group is 0
	// - cooling group is 999
	// - cooling group is 123

	tests := []struct{
		cdu CDU
		expectedErrors []error
	}{{
		// Negative test - cooling group number is negative
		CDU{
			CoolingGroup: -1,
		},
		[]error{
			errors.New("invalid CDU xname: d-1"),
			errors.New("invalid cooling group ordinal (-1) expected value between 0 and 999"),
		},
	},{
		// Negative test - cooling group number is greater then 999
		CDU{
			CoolingGroup: 1000,
		},
		[]error{
			errors.New("invalid cooling group ordinal (1000) expected value between 0 and 999"),
		},
	}, {
		// Negative test - cooling group number is greater then 999
		CDU{
			CoolingGroup: 3000,
		},
		[]error{
			errors.New("invalid cooling group ordinal (3000) expected value between 0 and 999"),
		},
	}, {
		// Positive Tests - cooling group is 0
		CDU{
			CoolingGroup: 0,
		},
		nil,
	},{
		// Positive Tests - cooling group is 999
		CDU{
			CoolingGroup: 123,
		},
		nil,
	}, {
		// Positive Tests - cooling group is 999
		CDU{
			CoolingGroup: 999,
		},
		nil,
	}}

	for _, test := range tests {
		err := test.cdu.ValidateEnhanced()

		var errors []error
		if err != nil {
			errors = err.(*multierror.Error).Errors
		}
		if !compareErrorSlices(test.expectedErrors, errors){
			t.Errorf("Unexpected validation error for %s: Expected errors: %v, Actual errors: %v", test.cdu, test.expectedErrors, errors)
		}
	}
}

func TestCDUMgmtSwitchValidationEnhanced(t *testing.T) {
	// CDUMgmtSwitch
	// Negative tests
	// - cooling group number is negative
	// - cooling group number is greater then 999 (verify call to the parent ValidateEnhanced worked)
	// - slot is negative
	// - slot is greater then 31
	// Positive tests
	// - slot is 0
	// - slot is 15
	// - slot is 31

	tests := []struct{
		component CDUMgmtSwitch
		expectedErrors []error
	}{{
		// Negative test - cooling group number is negative
		CDUMgmtSwitch{
			CoolingGroup: -1,
			Slot: 2,
		},
		[]error{
			errors.New("invalid CDUMgmtSwitch xname: d-1w2"),
			errors.New("invalid CDU xname: d-1"),
			errors.New("invalid cooling group ordinal (-1) expected value between 0 and 999"),
		},
	},{
		// Negative test - cooling group number is greater then 999 (verify call to the parent ValidateEnhanced worked)
		CDUMgmtSwitch{
			CoolingGroup: 1000,
			Slot: 2,
		},
		[]error{
			errors.New("invalid cooling group ordinal (1000) expected value between 0 and 999"),
		},
	}, {
		// Negative test - cooling group number is 1000 and the slot is negative
		CDUMgmtSwitch{
			CoolingGroup: 1000,
			Slot: -1,
		},
		[]error{
			errors.New("invalid CDUMgmtSwitch xname: d1000w-1"),
			errors.New("invalid cooling group ordinal (1000) expected value between 0 and 999"),
			errors.New("invalid slot ordinal (-1) expected value between 0 and 31"),
		},
	}, {
		// Negative test - slot is negative
		CDUMgmtSwitch{
			CoolingGroup: 1,
			Slot: -1,
		},
		[]error{
			errors.New("invalid CDUMgmtSwitch xname: d1w-1"),
			errors.New("invalid slot ordinal (-1) expected value between 0 and 31"),
		},
	}, {
		// Negative test - slot is greater then 31
		CDUMgmtSwitch{
			CoolingGroup: 1,
			Slot: 33,
		},
		[]error{
			errors.New("invalid slot ordinal (33) expected value between 0 and 31"),
		},
	},{
		// Positive Tests - slot is 0
		CDUMgmtSwitch{
			CoolingGroup: 0,
			Slot: 0,
		},
		nil,
	},{
		// Positive Tests - slot is 15
		CDUMgmtSwitch{
			CoolingGroup: 123,
			Slot: 15,
		},
		nil,
	}, {
		// Positive Tests - slot is 31
		CDUMgmtSwitch{
			CoolingGroup: 999,
			Slot: 31,
		},
		nil,
	}}

	for _, test := range tests {
		err := test.component.ValidateEnhanced()

		var errors []error
		if err != nil {
			errors = err.(*multierror.Error).Errors
		}
		if !compareErrorSlices(test.expectedErrors, errors){
			t.Errorf("Unexpected validation error for %s: Expected errors: %v, Actual errors: %v", test.component, test.expectedErrors, errors)
		}
	}
}

func TestCabinetSwitchValidationEnhanced(t *testing.T) {
	// Cabinet
	// Negative tests
	// - cabinet number is negative
	// - cabinet number is greater then 999
	// Positive tests
	// - cabinet is 0
	// - cabinet is 10
	// - cabinet is 999
}

func TestChassisValidationEnhanced(t *testing.T) {
	// Chassis
	// - River
	// 	 	Negative tests
	// 		- cabinet number is negative
	// 		- cabinet number is greater then 999
	//		- Chassis is not 0 (1-7)
	//	 	Positive tests
	//		- Chassis is 0
	// - Hill
	// 	 	Negative tests
	// 		- cabinet number is negative
	// 		- cabinet number is greater then 999
	//		- Chassis is not 1 or 3 (0, 2, 4-7)
	//	 	Positive tests
	//		- Chassis is not 1 or 3
	// - Mountain
	// 	 	Negative tests
	// 		- cabinet number is negative
	// 		- cabinet number is greater then 999
	//		- Chassis is negative
	// 		- Chassis is greater then 7
	//	 	Positive tests
	//		- Chassis is 0
	//		- Chassis is 4
	// 		- Chassis is 7
}

func TestChassisBMCValidationEnhanced(t *testing.T) {

}

func TestMgmtSwitchValidationEnhanced(t *testing.T) {

}

func TestMgmtSwitchConnectorValidationEnhanced(t *testing.T) {

}

func TestMgmtHLSwitchEnclosureValidationEnhanced(t *testing.T) {

}

func TestMgmtHLSwitchValidationEnhanced(t *testing.T) {

}

func TestRouterModuleValidationEnhanced(t *testing.T) {

}

func TestRouterBMCValidationEnhanced(t *testing.T) {

}

func TestComputeModuleValidationEnhanced(t *testing.T) {

}

func TestNodeBMCValidationEnhanced(t *testing.T) {

}

func TestNodeValidationEnhanced(t *testing.T) {

}

//
//
// Test Helpers
//
//

func compareErrorSlices(x, y []error) bool {
	if len(x) != len(y) {
		return false
	}

	for i, errorX := range x {
		errorY := y[i]

		if errorX.Error() != errorY.Error() {
			return false
		}
	}

	return true
} 