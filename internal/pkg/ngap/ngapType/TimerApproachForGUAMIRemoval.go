package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	TimerApproachForGUAMIRemovalPresentApplyTimer aper.Enumerated = 0
)

type TimerApproachForGUAMIRemoval struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
