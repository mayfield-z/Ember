package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	HandoverFlagPresentHandoverPreparation aper.Enumerated = 0
)

type HandoverFlag struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
