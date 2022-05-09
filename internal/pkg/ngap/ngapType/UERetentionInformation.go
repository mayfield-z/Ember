package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	UERetentionInformationPresentUesRetained aper.Enumerated = 0
)

type UERetentionInformation struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
