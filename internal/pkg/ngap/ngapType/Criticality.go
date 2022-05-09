package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	CriticalityPresentReject aper.Enumerated = 0
	CriticalityPresentIgnore aper.Enumerated = 1
	CriticalityPresentNotify aper.Enumerated = 2
)

type Criticality struct {
	Value aper.Enumerated `aper:"valueLB:0,valueUB:2"`
}
