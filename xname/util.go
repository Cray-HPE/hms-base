package xname

import (
	"errors"
	"strconv"

	base "github.com/Cray-HPE/hms-base"
)

var ErrUnknownStruct = errors.New("unable to determine HMS Type from struct")

func GetHMSType(obj interface{}) (base.HMSType, error) {
	// Handy bash fragment to generate the type switch below
	// for hms_type in $(cat ./xname/types.go | grep '^type' | awk '{print $2}'); do
	// echo "	case $hms_type, *$hms_type:"
	// echo "		return base.$hms_type, nil"
	// done
	switch obj.(type) {
	case System, *System:
		return base.System, nil
	case CDU, *CDU:
		return base.CDU, nil
	case CDUMgmtSwitch, *CDUMgmtSwitch:
		return base.CDUMgmtSwitch, nil
	case Cabinet, *Cabinet:
		return base.Cabinet, nil
	case CabinetPDUController, *CabinetPDUController:
		return base.CabinetPDUController, nil
	case Chassis, *Chassis:
		return base.Chassis, nil
	case MgmtSwitch, *MgmtSwitch:
		return base.MgmtSwitch, nil
	case MgmtSwitchConnector, *MgmtSwitchConnector:
		return base.MgmtSwitchConnector, nil
	case MgmtHLSwitchEnclosure, *MgmtHLSwitchEnclosure:
		return base.MgmtHLSwitchEnclosure, nil
	case MgmtHLSwitch, *MgmtHLSwitch:
		return base.MgmtHLSwitch, nil
	case RouterModule, *RouterModule:
		return base.RouterModule, nil
	case RouterBMC, *RouterBMC:
		return base.RouterBMC, nil
	case ComputeModule, *ComputeModule:
		return base.ComputeModule, nil
	case NodeBMC, *NodeBMC:
		return base.NodeBMC, nil
	case Node, *Node:
		return base.Node, nil
	}

	return base.HMSTypeInvalid, ErrUnknownStruct
}

func FromString(xname string) (interface{}, base.HMSType) {
	hmsType := base.GetHMSType(xname)
	if hmsType == base.HMSTypeInvalid {
		return nil, hmsType
	}

	re, err := base.GetHMSTypeRegex(hmsType)
	if err != nil {
		return nil, base.HMSTypeInvalid
	}

	_, argCount, err := base.GetHMSTypeFormatString(hmsType)
	if err != nil {
		return nil, base.HMSTypeInvalid
	}

	matchesRaw := re.FindStringSubmatch(xname)
	if (argCount + 1) != len(matchesRaw) {
		return nil, base.HMSTypeInvalid
	}

	// If we have gotten to this point these matches should be integers, so we can safely convert them
	// to integers from strings.
	matches := []int{}
	for _, matchRaw := range matchesRaw[1:] {
		match, err := strconv.Atoi(matchRaw)
		if err != nil {
			return nil, base.HMSTypeInvalid
		}

		matches = append(matches, match)
	}

	var component interface{}

	switch hmsType {
	case base.System:
		component = System{}
	case base.CDU:
		component = CDU{
			CoolingGroup: matches[0],
		}
	case base.CDUMgmtSwitch:
		component = CDUMgmtSwitch{
			CoolingGroup: matches[0],
			Slot:         matches[1],
		}
	case base.Cabinet:
		component = Cabinet{
			Cabinet: matches[0],
		}
	case base.CabinetPDUController:
		component = CabinetPDUController{
			Cabinet:       matches[0],
			PDUController: matches[1],
		}
	case base.Chassis:
		component = Chassis{
			Cabinet: matches[0],
			Chassis: matches[1],
		}
	case base.MgmtSwitch:
		component = MgmtSwitch{
			Cabinet: matches[0],
			Chassis: matches[1],
			Slot:    matches[2],
		}
	case base.MgmtSwitchConnector:
		component = MgmtSwitchConnector{
			Cabinet:    matches[0],
			Chassis:    matches[1],
			Slot:       matches[2],
			SwitchPort: matches[3],
		}
	case base.MgmtHLSwitchEnclosure:
		component = MgmtHLSwitchEnclosure{
			Cabinet: matches[0],
			Chassis: matches[1],
			Slot:    matches[2],
		}
	case base.MgmtHLSwitch:
		component = MgmtHLSwitch{
			Cabinet: matches[0],
			Chassis: matches[1],
			Slot:    matches[2],
			Space:   matches[3],
		}
	case base.RouterModule:
		component = RouterModule{
			Cabinet: matches[0],
			Chassis: matches[1],
			Slot:    matches[2],
		}
	case base.RouterBMC:
		component = RouterBMC{
			Cabinet: matches[0],
			Chassis: matches[1],
			Slot:    matches[2],
			BMC:     matches[3],
		}
	case base.ComputeModule:
		component = ComputeModule{
			Cabinet: matches[0],
			Chassis: matches[1],
			Slot:    matches[2],
		}
	case base.NodeBMC:
		component = NodeBMC{
			Cabinet: matches[0],
			Chassis: matches[1],
			Slot:    matches[2],
			BMC:     matches[3],
		}
	case base.Node:
		component = Node{
			Cabinet: matches[0],
			Chassis: matches[1],
			Slot:    matches[2],
			BMC:     matches[3],
			Node:    matches[4],
		}
	default:
		// TODO should this be a generic error? This means that this is not apart of the struct family of xnames.
		return nil, base.HMSTypeInvalid
	}
	return component, hmsType
}
