package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	UEPresencePresentIn      aper.Enumerated = 0
	UEPresencePresentOut     aper.Enumerated = 1
	UEPresencePresentUnknown aper.Enumerated = 2
)

type UEPresence struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:2"`
}
