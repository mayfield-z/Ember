package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type AMFRegionID struct {
	Value aper.BitString `aper:"sizeLB:8,sizeUB:8"`
}
