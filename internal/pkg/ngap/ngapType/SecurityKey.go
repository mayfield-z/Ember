package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type SecurityKey struct {
	Value aper.BitString `aper:"sizeLB:256,sizeUB:256"`
}
