package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	N3IWFIDPresentNothing int = iota /* No components present */
	N3IWFIDPresentN3IWFID
	N3IWFIDPresentChoiceExtensions
)

type N3IWFID struct {
	Present          int
	N3IWFID          *aper.BitString `aper:"sizeLB:16,sizeUB:16"`
	ChoiceExtensions *ProtocolIESingleContainerN3IWFIDExtIEs
}
