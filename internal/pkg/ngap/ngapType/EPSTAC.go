package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type EPSTAC struct {
	Value aper.OctetString `aper:"sizeLB:2,sizeUB:2"`
}
