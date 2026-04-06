package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// LPSetup is the Load Profile Setup.
type LPSetup struct {
	LogicalName    core.ObisCode
	CaptureObjects []string
	CapturePeriod  uint16
	ProfileEntries uint32
	Version        uint8
}

func (lp *LPSetup) ClassID() uint16               { return core.ClassIDLpSetup }
func (lp *LPSetup) GetLogicalName() core.ObisCode { return lp.LogicalName }
func (lp *LPSetup) GetVersion() uint8             { return lp.Version }
func (lp *LPSetup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessRead}
}
