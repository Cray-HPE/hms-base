package xname

import (
	"fmt"

	base "github.com/Cray-HPE/hms-base"
	multierror "github.com/hashicorp/go-multierror"
)

// s0
type System struct{}

func (s System) String() string {
	return "s0"
}

func (s System) Validate() error {
	return nil
}

func (s System) ValidateEnhanced() error {
	return nil
}

func (s System) CDU(coolingGroup int) CDU {
	return CDU{
		CoolingGroup: coolingGroup,
	}
}

func (s System) Cabinet(cabinet int) Cabinet {
	return Cabinet{
		Cabinet: cabinet,
	}
}

// dD
type CDU struct {
	CoolingGroup int // D: 0-999
}

func (c CDU) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.CDU)
	return fmt.Sprintf(formatStr, c.CoolingGroup)
}

func (c CDU) Parent() System {
	return System{}
}

func (c CDU) Validate() error {
	xname := c.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid CDU xname: %s", xname)
	}

	return nil
}

func (c CDU) ValidateEnhanced() error {
	var result error

	// Perform normal validation
	if err := c.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)  
	}

	if err := c.Parent().ValidateEnhanced(); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)  
	}

	if !(0 <= c.CoolingGroup && c.CoolingGroup <= 999) {
		// Cooling group range
		err := fmt.Errorf("invalid cooling group ordinal (%v) expected value between 0 and 999", c.CoolingGroup)
		result = multierror.Append(result, err)  
	}

	return result
}

func (c CDU) CDUMgmtSwitch(slot int) CDUMgmtSwitch {
	return CDUMgmtSwitch{
		CoolingGroup: c.CoolingGroup,
		Slot:         slot,
	}
}

// dDwW
type CDUMgmtSwitch struct {
	CoolingGroup int // D: 0-999
	Slot         int // W: 0-31
}

func (cms CDUMgmtSwitch) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.CDUMgmtSwitch)
	return fmt.Sprintf(formatStr, cms.CoolingGroup, cms.Slot)
}

func (cms CDUMgmtSwitch) Parent() CDU {
	return CDU{
		CoolingGroup: cms.CoolingGroup,
	}
}

func (cms CDUMgmtSwitch) Validate() error {
	xname := cms.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid CDUMgmtSwitch xname: %s", xname)
	}

	return nil
}

func (cms CDUMgmtSwitch) ValidateEnhanced() error {
	var result error

	// Perform normal validation
	if err := cms.Validate(); err  != nil {
		// Xname is not valid
		result = multierror.Append(result, err)  
	}

	if err := cms.Parent().ValidateEnhanced(); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)  
	}

	if !(0 <= cms.Slot && cms.Slot <= 31) {
		// CDU Switch slot
		err := fmt.Errorf("invalid slot ordinal (%v) expected value between 0 and 31", cms.Slot)
		result = multierror.Append(result, err)  
	}

	return result
}

// xX
type Cabinet struct {
	Cabinet int // X: 0-999
}

func (c Cabinet) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.Cabinet)
	return fmt.Sprintf(formatStr, c.Cabinet)
}

func (c Cabinet) Parent() System {
	return System{}
}

func (c Cabinet) Validate() error {
	xname := c.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid Cabinet xname: %s", xname)
	}

	return nil
}

func (c Cabinet) ValidateEnhanced() error {
	var result error

	// Perform normal validation
	if err := c.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)  
	}

	if err := c.Parent().ValidateEnhanced(); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)  
	}

	if !(0 <= c.Cabinet && c.Cabinet <= 999) {
		// Cabinet number out of range
		err := fmt.Errorf("invalid cabinet ordinal (%v) expected value between 0 and 999", c.Cabinet)
		result = multierror.Append(result, err)  
	}

	return result
}

func (c Cabinet) Chassis(chassis int) Chassis {
	return Chassis{
		Cabinet: c.Cabinet,
		Chassis: chassis,
	}
}

func (c Cabinet) CabinetPDUController(pduController int) CabinetPDUController {
	return CabinetPDUController{
		Cabinet:       c.Cabinet,
		PDUController: pduController,
	}
}

// xXmM
type CabinetPDUController struct {
	Cabinet       int // X: 0-999
	PDUController int // M: 0-3
}

func (p CabinetPDUController) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.CabinetPDUController)
	return fmt.Sprintf(formatStr, p.Cabinet, p.Cabinet)
}

func (p CabinetPDUController) Parent() Cabinet {
	return Cabinet{
		Cabinet: p.Cabinet,
	}
}

func (p CabinetPDUController) Validate() error {
	xname := p.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid CabinetPDUController xname: %s", xname)
	}

	return nil
}

func (p CabinetPDUController) ValidateEnhanced() error {
	var result error

	// Perform normal validation
	if err := p.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)  
	}

	if err := p.Parent().ValidateEnhanced(); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)  

	}

	if !(0 <= p.PDUController && p.PDUController <= 3) {
		// Cabinet number out of range
		err := fmt.Errorf("invalid pdu controller ordinal (%v) expected value between 0 and 3", p.PDUController)
		result = multierror.Append(result, err)  
	}

	return result
}

// xXcC
// Mountain Have 8 c0-c8
// Hill have 2 c1 and c3
// River always have c0
type Chassis struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
}

func (c Chassis) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.Chassis)
	return fmt.Sprintf(formatStr, c.Cabinet, c.Chassis)
}

func (c Chassis) Parent() Cabinet {
	return Cabinet{
		Cabinet: c.Cabinet,
	}
}

func (c Chassis) Validate() error {
	xname := c.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid Chassis xname: %s", xname)
	}

	return nil
}

func (c Chassis) ValidateEnhanced(class base.HMSClass) error {
	var result error

	// Perform normal validation
	if err := c.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := c.Parent().ValidateEnhanced(); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// Chassis Validation
	switch class {
	case base.ClassRiver:
		if c.Chassis != 0 {
			// River chassis must be equal to 0
			err := fmt.Errorf("invalid river chassis ordinal (%v) expected 0", c.Chassis)
			result = multierror.Append(result, err)
		}
	case base.ClassHill:
		if !(c.Chassis == 1 || c.Chassis == 3) {
			// Hill has Chassis 1 or 3
			err := fmt.Errorf("invalid hill chassis ordinal (%v) expected 1 or 3", c.Chassis)
			result = multierror.Append(result, err)
		}
	case base.ClassMountain:
		if !(0 <= c.Chassis && c.Chassis <= 7) {
			// Mountain must chassis between 0 and 7
			err := fmt.Errorf("invalid hill chassis ordinal (%v) expected value between 0 and 7", c.Chassis)
			result = multierror.Append(result, err)
		}
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)
	}

	return result
}

func (c Chassis) ChassisBMC(bmc int) ChassisBMC {
	return ChassisBMC{
		Cabinet: c.Cabinet,
		Chassis: c.Chassis,
		BMC: bmc,
	}
}

func (c Chassis) MgmtHLSwitchEnclosure(slot int) MgmtHLSwitchEnclosure {
	return MgmtHLSwitchEnclosure{
		Cabinet: c.Cabinet,
		Chassis: c.Chassis,
		Slot:    slot,
	}
}

func (c Chassis) MgmtSwitch(slot int) MgmtSwitch {
	return MgmtSwitch{
		Cabinet: c.Cabinet,
		Chassis: c.Chassis,
		Slot:    slot,
	}
}

// This is a convience function, as we normally do not work with MgmtHLSwitchEnclosures directly
func (c Chassis) MgmtHLSwitch(slot, space int) MgmtHLSwitch {
	return c.MgmtHLSwitchEnclosure(slot).MgmtHLSwitch(space)
}

func (c Chassis) RouterModule(slot int) RouterModule {
	return RouterModule{
		Cabinet: c.Cabinet,
		Chassis: c.Chassis,
		Slot:    slot,
	}
}

// This is a convince function, as we normally do not work with RouterModules directly.
func (c Chassis) RouterBMC(slot, bmc int) RouterBMC {
	return c.RouterModule(slot).RouterBMC(bmc)
}

func (c Chassis) ComputeModule(slot int) ComputeModule {
	return ComputeModule{
		Cabinet: c.Cabinet,
		Chassis: c.Chassis,
		Slot:    slot,
	}
}

func (c Chassis) NodeBMC(slot, bmc int) NodeBMC {
	return c.ComputeModule(slot).NodeBMC(bmc)
}

// xXcCbB
// Mountain and Hill have only b0
// River does not have ChassisBMCs
type ChassisBMC struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
	BMC int // B: 0
}

func (c ChassisBMC) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.ChassisBMC)
	return fmt.Sprintf(formatStr, c.Cabinet, c.Chassis, c.BMC)
}

func (c ChassisBMC) Parent() Chassis {
	return Chassis{
		Cabinet: c.Cabinet,
		Chassis: c.Chassis,
	}
}

func (c ChassisBMC) Validate() error {
	xname := c.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid Chassis xname: %s", xname)
	}

	return nil
}

func (c ChassisBMC) ValidateEnhanced(class base.HMSClass) error {
	var result error

	// Perform normal validation
	if err := c.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := c.Parent().ValidateEnhanced(class); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// Chassis Validation
	switch class {
	case base.ClassRiver:
		// River does not have ChassisBMCs 
		err := fmt.Errorf("invalid - chassis bmcs do not exist for river")
		result = multierror.Append(result, err)
	case base.ClassHill:
		fallthrough
	case base.ClassMountain:
		if c.BMC != 0 {
			// Mountain and Hill must have b0 for there ChassisBMC
			err := fmt.Errorf("invalid chassis bmc ordinal (%v) expected value is 0", c.BMC)
			result = multierror.Append(result, err)
		}
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)
	}

	return result
}


// xXcCwW
type MgmtSwitch struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
	Slot    int // W: 1-48
}

func (ms MgmtSwitch) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.MgmtSwitch)
	return fmt.Sprintf(formatStr, ms.Cabinet, ms.Chassis, ms.Slot)
}

func (ms MgmtSwitch) Parent() Chassis {
	return Chassis{
		Cabinet: ms.Cabinet,
		Chassis: ms.Chassis,
	}
}

func (ms MgmtSwitch) Validate() error {
	xname := ms.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid MgmtSwitch xname: %s", xname)
	}

	return nil
}

func (ms MgmtSwitch) ValidateEnhanced(class base.HMSClass) error {
	var result error

	// Perform normal validation
	if err := ms.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := ms.Parent().ValidateEnhanced(class); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// Chassis Validation
	switch class {
	case base.ClassRiver:
		// Expected to be river only
		if !(1 <= ms.Slot && ms.Slot <= 48) {
			// Verify that the U is within a standard rack slot
			err := fmt.Errorf("invalid rack slot ordinal (%v) expected value is between 1 and 48", ms.Slot)			
			result = multierror.Append(result, err)
		}
	case base.ClassHill:
		fallthrough
	case base.ClassMountain:
		// MgmtSwitches are only for river
		err := fmt.Errorf("invalid - mgmt switches do not exist for %s", class)
		result = multierror.Append(result, err)
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)
	}

	return result
}

func (ms MgmtSwitch) MgmtSwitchConnector(switchPort int) MgmtSwitchConnector {
	return MgmtSwitchConnector{
		Cabinet:    ms.Cabinet,
		Chassis:    ms.Chassis,
		Slot:       ms.Slot,
		SwitchPort: switchPort,
	}
}

// xXcCwWjJ
type MgmtSwitchConnector struct {
	Cabinet    int // X: 0-999
	Chassis    int // C: 0-7
	Slot       int // W: 1-48
	SwitchPort int // J: 1-32 // TODO the HSOS page, should allow upto at least 48
}

func (msc MgmtSwitchConnector) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.MgmtSwitchConnector)
	return fmt.Sprintf(formatStr, msc.Cabinet, msc.Chassis, msc.Slot, msc.SwitchPort)
}

func (msc MgmtSwitchConnector) Parent() MgmtSwitch {
	return MgmtSwitch{
		Cabinet: msc.Cabinet,
		Chassis: msc.Chassis,
		Slot:    msc.Slot,
	}
}

func (msc MgmtSwitchConnector) Validate() error {
	xname := msc.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid MgmtSwitchConnector xname: %s", xname)
	}

	return nil
}

func (msc MgmtSwitchConnector) ValidateEnhanced(class base.HMSClass) error {
	var result error

	// Perform normal validation
	if err := msc.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := msc.Parent().ValidateEnhanced(class); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// Chassis Validation
	switch class {
	case base.ClassRiver:
		// Expected to be river only
		if !(1 <= msc.SwitchPort) {
			// Verify that the switch port is valid
			err := fmt.Errorf("invalid switch port ordinal (%v) expected greater than 1", msc.SwitchPort)
			result = multierror.Append(result, err)
		}
	case base.ClassHill:
		fallthrough
	case base.ClassMountain:
		// MgmtSwitchConnectors are only for river
		err := fmt.Errorf("invalid - mgmt switch connectors do not exist for %s", class)
		result = multierror.Append(result, err)
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)
	}

	return result
}

// xXcChH
type MgmtHLSwitchEnclosure struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
	Slot    int // H: 1-48
}

func (enclosure MgmtHLSwitchEnclosure) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.MgmtHLSwitchEnclosure)
	return fmt.Sprintf(formatStr, enclosure.Cabinet, enclosure.Chassis, enclosure.Slot)
}

func (enclosure MgmtHLSwitchEnclosure) Parent() Chassis {
	return Chassis{
		Cabinet: enclosure.Cabinet,
		Chassis: enclosure.Chassis,
	}
}

func (enclosure MgmtHLSwitchEnclosure) Validate() error {
	xname := enclosure.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid MgmtHLSwitchEnclosure xname: %s", xname)
	}

	return nil
}

func (enclosure MgmtHLSwitchEnclosure) ValidateEnhanced(class base.HMSClass) error {
	var result error

	// Perform normal validation
	if err := enclosure.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := enclosure.Parent().ValidateEnhanced(class); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// Chassis Validation
	switch class {
	case base.ClassRiver:
		// Expected to be river only

		if !(1 <= enclosure.Slot && enclosure.Slot <= 48) {
			// Verify that the U is within a standard rack
			err := fmt.Errorf("invalid rack slot ordinal (%v) expected value is between 1 and 48", enclosure.Slot)
			result = multierror.Append(result, err)
		}
	case base.ClassHill:
		fallthrough
	case base.ClassMountain:
		// MgmtHLSwitchEnclosure are only for river
		err := fmt.Errorf("invalid - mgmt hl switch enclosures do not exist for %s", class)
		result = multierror.Append(result, err)
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)
	}

	return result
}

func (enclosure MgmtHLSwitchEnclosure) MgmtHLSwitch(space int) MgmtHLSwitch {
	return MgmtHLSwitch{
		Cabinet: enclosure.Cabinet,
		Chassis: enclosure.Chassis,
		Slot:    enclosure.Slot,
		Space:   space,
	}
}

//xXcChHsS
type MgmtHLSwitch struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
	Slot    int // H: 1-48
	Space   int // S: 1-4
}

func (mhls MgmtHLSwitch) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.MgmtHLSwitch)
	return fmt.Sprintf(formatStr, mhls.Cabinet, mhls.Chassis, mhls.Slot, mhls.Space)
}

func (mhls MgmtHLSwitch) Parent() MgmtHLSwitchEnclosure {
	return MgmtHLSwitchEnclosure{
		Cabinet: mhls.Cabinet,
		Chassis: mhls.Chassis,
		Slot:    mhls.Slot,
	}
}

func (mhls MgmtHLSwitch) Validate() error {
	xname := mhls.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid MgmtHLSwitch xname: %s", xname)
	}

	return nil
}

func (mhls MgmtHLSwitch) ValidateEnhanced(class base.HMSClass) error {
	var result error

	// Perform normal validation
	if err := mhls.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := mhls.Parent().ValidateEnhanced(class); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// Chassis Validation
	switch class {
	case base.ClassRiver:
		// Expected to be river only
		if !(1 <= mhls.Space && mhls.Space <= 4) {
			// Verify a valid space value
			err := fmt.Errorf("invalid space ordinal (%v) expected value is between 1 and 4", mhls.Space)
			result = multierror.Append(result, err)
		}
	case base.ClassHill:
		fallthrough
	case base.ClassMountain:
		// MgmtHLSwitch are only for river
		err := fmt.Errorf("invalid - mgmt hl switches do not exist for %s", class)
		result = multierror.Append(result, err)
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)

	}

	return result
}

// xXcCrR
// Mountain/Hill: R: 0-8
// River: 1-48
type RouterModule struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
	Slot    int // R: 0-64
}

func (rm RouterModule) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.RouterModule)
	return fmt.Sprintf(formatStr, rm.Cabinet, rm.Chassis, rm.Slot)
}

func (rm RouterModule) Parent() Chassis {
	return Chassis{
		Cabinet: rm.Cabinet,
		Chassis: rm.Chassis,
	}
}

func (rm RouterModule) Validate() error {
	xname := rm.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid RouterModule xname: %s", xname)
	}

	return nil
}

func (rm RouterModule) ValidateEnhanced(class base.HMSClass) error {
	var result error

	// Perform normal validation
	if err := rm.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := rm.Parent().ValidateEnhanced(class); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// Router Module Validation
	switch class {
	case base.ClassRiver:
		if !(1 <= rm.Slot && rm.Slot <= 48) {
			// Standard Rack size			
			err := fmt.Errorf("invalid rack slot ordinal (%v) expected value is between 1 and 48", rm.Slot)
			result = multierror.Append(result, err)
		}
	case base.ClassHill:
		fallthrough
	case base.ClassMountain:
		if !(0 <= rm.Slot && rm.Slot <= 7) {
			err := fmt.Errorf("invalid chassis slot ordinal (%v) expected value is between 0 and 7", rm.Slot)
			result = multierror.Append(result, err)
		}
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)
	}

	return result
}

func (rm RouterModule) RouterBMC(bmc int) RouterBMC {
	return RouterBMC{
		Cabinet: rm.Cabinet,
		Chassis: rm.Chassis,
		Slot:    rm.Slot,
		BMC:     bmc,
	}
}

// xXcCrRbB
// B is always 0
type RouterBMC struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
	Slot    int // R: 0-64
	BMC     int // B: 0
}

func (bmc RouterBMC) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.RouterBMC)
	return fmt.Sprintf(formatStr, bmc.Cabinet, bmc.Chassis, bmc.Slot, bmc.BMC)
}

func (bmc RouterBMC) Parent() RouterModule {
	return RouterModule{
		Cabinet: bmc.Cabinet,
		Chassis: bmc.Chassis,
		Slot:    bmc.Slot,
	}
}

func (bmc RouterBMC) Validate() error {
	xname := bmc.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid RouterBMC xname: %s", xname)
	}

	return nil
}

func (bmc RouterBMC) ValidateEnhanced(class base.HMSClass) error {
	var result error

	// Perform normal validation
	if err := bmc.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := bmc.Parent().ValidateEnhanced(class); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// Router BMC Validation
	if bmc.BMC != 0 {
		// BMC should always be 0
		err := fmt.Errorf("invalid router bmc ordinal (%v) expected value is 0", bmc.BMC)
		result = multierror.Append(result, err)
	}

	return result
}

// xXcCsS
// Mountain/Hill: 0-7
// River: 1-48
type ComputeModule struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
	Slot    int // S: 1-63
}

func (cm ComputeModule) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.ComputeModule)
	return fmt.Sprintf(formatStr, cm.Cabinet, cm.Chassis, cm.Slot)
}

func (cm ComputeModule) Parent() Chassis {
	return Chassis{
		Cabinet: cm.Cabinet,
		Chassis: cm.Chassis,
	}
}

func (cm ComputeModule) Validate() error {
	xname := cm.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid ComputeModule xname: %s", xname)
	}

	return nil
}

func (cm ComputeModule) ValidateEnhanced(class base.HMSClass) error {
	var result error

	// Perform normal validation
	if err := cm.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := cm.Parent().ValidateEnhanced(class); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// Compute Module Validation
	switch class {
	case base.ClassRiver:
		if !(1 <= cm.Slot && cm.Slot <= 48) {
			// Standard Rack size		
			err := fmt.Errorf("invalid rack slot ordinal (%v) expected value is between 1 and 48", cm.Slot)	
			result = multierror.Append(result, err)
		}
	case base.ClassHill:
		fallthrough
	case base.ClassMountain:
		if !(0 <= cm.Slot && cm.Slot <= 7) {
			// Mountain Chassis
			err := fmt.Errorf("invalid chassis slot ordinal (%v) expected value is between 0 and 7", cm.Slot)
			result = multierror.Append(result, err)			
		}
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)
	}

	return result
}

func (cm ComputeModule) NodeBMC(bmc int) NodeBMC {
	return NodeBMC{
		Cabinet: cm.Cabinet,
		Chassis: cm.Chassis,
		Slot:    cm.Slot,
		BMC:     bmc,
	}
}

// xXcCsSbB
// Node Card/Node BMC
// Mountain/Hill can be 0 or 1
// River
// - Single node chassis: always 0
// - Dual Node chassis: 1 or 2
// - Dense/Quad Node Chassis 1-4
type NodeBMC struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
	Slot    int // S: 1-63
	BMC     int // B: 0-1
}

func (bmc NodeBMC) Parent() ComputeModule {
	return ComputeModule{
		Cabinet: bmc.Cabinet,
		Chassis: bmc.Chassis,
		Slot:    bmc.Slot,
	}
}

func (bmc NodeBMC) Validate() error {
	xname := bmc.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid NodeBMC xname: %s", xname)
	}

	return nil
}

func (bmc NodeBMC) ValidateEnhanced(class base.HMSClass, nodeChassisType NodeBladeType) error {
	var result error

	// Perform normal validation
	if err :=bmc.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := bmc.Parent().ValidateEnhanced(class); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}

	// NodeBMC Validation
	switch class {
	case base.ClassRiver:
		switch nodeChassisType {
		case SingleNodeBlade:
			if bmc.BMC != 0 {
				// Single node chassis must have the BMC as 0
				err := fmt.Errorf("invalid bmc ordinal (%v) expected value for a single node chassis is 0", bmc.BMC)
				result = multierror.Append(result, err)
			}
		case DualNodeBlade:
			if !(bmc.BMC == 1 || bmc.BMC == 2) {
				// Dual node chassis must have BMC as 1 or 2
				err := fmt.Errorf("invalid bmc ordinal (%v) expected values for a dual node chassis are 1 or 2", bmc.BMC)
				result = multierror.Append(result, err)
			}
		case QuadNodeBlade:
			if !(1 <= bmc.BMC || bmc.BMC <= 4) {
				// Dense Quad node chassis must have BMC between 1 and 4
				err := fmt.Errorf("invalid bmc ordinal (%v) expected values for a quad node chassis are 1 to 4", bmc.BMC)
				result = multierror.Append(result, err)
			}
		default:
			err := fmt.Errorf("unknown node chassis type (%v)", nodeChassisType)
			result = multierror.Append(result, err)
		}
	case base.ClassHill:
		fallthrough
	case base.ClassMountain:
		if !(bmc.BMC == 0 || bmc.BMC == 1) {
			// Mountain blades have 2 BMCs
			err := fmt.Errorf("invalid bmc ordinal (%v) expected values for a mountain node bmc are 0 or 1", bmc.BMC)
			result = multierror.Append(result, err)
		}
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)
	}

	return result
}

func (bmc NodeBMC) Node(node int) Node {
	return Node{
		Cabinet: bmc.Cabinet,
		Chassis: bmc.Chassis,
		Slot:    bmc.Slot,
		BMC:     bmc.BMC,
		Node:    node,
	}
}

func (bmc NodeBMC) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.NodeBMC)
	return fmt.Sprintf(formatStr, bmc.Cabinet, bmc.Chassis, bmc.Slot, bmc.BMC)
}

// xCcCsSbBnN
// River - Always 0
// Mountain/Hill 0 or 1
type Node struct {
	Cabinet int // X: 0-999
	Chassis int // C: 0-7
	Slot    int // S: 1-63
	BMC     int // B: 0-1 - TODO the HSOS document is wrong here. as we do actually use greater than 1
	Node    int // N: 0-7
}

func (n Node) String() string {
	formatStr, _, _ := base.GetHMSTypeFormatString(base.Node)
	return fmt.Sprintf(formatStr, n.Cabinet, n.Chassis, n.Slot, n.BMC, n.Node)
}

func (n Node) Validate() error {
	xname := n.String()
	if !base.IsHMSCompIDValid(xname) {
		return fmt.Errorf("invalid node xname: %s", xname)
	}

	return nil
}

type NodeBladeType int // TODO Idk if this should be blade or chassis. This this could apply to both river and mountain

const (
	SingleNodeBlade NodeBladeType = iota
	DualNodeBlade
	QuadNodeBlade // TODO Should this have "dense"
)

func (n Node) ValidateEnhanced(class base.HMSClass, nodeChassisType NodeBladeType) error {
	var result error

	// Perform normal validation
	if err := n.Validate(); err != nil {
		// Xname is not valid
		result = multierror.Append(result, err)
	}

	if err := n.Parent().ValidateEnhanced(class, nodeChassisType); err != nil {
		// Verify all parents are valid 
		result = multierror.Append(result, err)
	}
	
	// Node Validation
	switch class {
	case base.ClassRiver:
		if n.Node != 0 {
			// River node value must be 0
			err := fmt.Errorf("invalid node ordinal (%v) expected value for a river node is 0", n.Node)
			result = multierror.Append(result, err)
		}
	case base.ClassHill:
		fallthrough
	case base.ClassMountain:
		switch nodeChassisType {
		case SingleNodeBlade:
			// We don't have this?
		case DualNodeBlade:
			if n.Node != 0 {
				// On a mountain dual node blade, each BMC controls 1 node.
				err := fmt.Errorf("invalid node ordinal (%v) expected value for a mountain dual node blade is 0", n.Node)
				result = multierror.Append(result, err)
			}
		case QuadNodeBlade:
			if !(n.Node == 0 || n.Node == 1) {
				// Dual node blade must have BMC as 1 or 2
				err := fmt.Errorf("invalid node ordinal (%v) expected values for a mountain quad node blade are 0 or 1", n.Node)
				result = multierror.Append(result, err)
			}
		default:
			err := fmt.Errorf("unknown node chassis type (%v)", nodeChassisType)
			result = multierror.Append(result, err)
		}
	default:
		err := fmt.Errorf("unknown HMSClass value (%v)", class)
		result = multierror.Append(result, err)
	}

	return result
}

func (n Node) Parent() NodeBMC {
	return NodeBMC{
		Cabinet: n.Cabinet,
		Chassis: n.Chassis,
		Slot:    n.Slot,
		BMC:     n.BMC,
	}
}
