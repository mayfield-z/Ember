package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type GTPTEID struct {
	Value aper.OctetString `aper:"sizeLB:4,sizeUB:4"`
}
