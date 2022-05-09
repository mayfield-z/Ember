package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type AMFSetID struct {
	Value aper.BitString `aper:"sizeLB:10,sizeUB:10"`
}
