package xname

import (
	"testing"

	base "github.com/Cray-HPE/hms-base"
)

func TestFoo(t *testing.T) {
	n := Node{
		Cabinet: 1000,
		Chassis: 1,
		Slot: 7,
		BMC: 1,
		Node: 0,
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
		Slot: 7,
		BMC: 1,
		Node: 0,
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
		Cabinet: 1, // X: 0-999
		Chassis: 0,       // C: 0-7
		Slot:    32,    // W: 1-48
	}
	t.Log("MgmtSwitch:", ms)

}	