package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	DelayCriticalPresentDelayCritical    aper.Enumerated = 0
	DelayCriticalPresentNonDelayCritical aper.Enumerated = 1
)

type DelayCritical struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
