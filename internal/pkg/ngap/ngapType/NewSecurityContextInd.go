package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	NewSecurityContextIndPresentTrue aper.Enumerated = 0
)

type NewSecurityContextInd struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
