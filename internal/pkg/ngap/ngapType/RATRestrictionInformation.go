package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type RATRestrictionInformation struct {
	Value aper.BitString `aper:"sizeExt,sizeLB:8,sizeUB:8"`
}
