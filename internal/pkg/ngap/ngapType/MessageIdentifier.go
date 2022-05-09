package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type MessageIdentifier struct {
	Value aper.BitString `aper:"sizeLB:16,sizeUB:16"`
}
