package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	ConfidentialityProtectionIndicationPresentRequired  aper.Enumerated = 0
	ConfidentialityProtectionIndicationPresentPreferred aper.Enumerated = 1
	ConfidentialityProtectionIndicationPresentNotNeeded aper.Enumerated = 2
)

type ConfidentialityProtectionIndication struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:2"`
}
