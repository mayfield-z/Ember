package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	GNBIDPresentNothing int = iota /* No components present */
	GNBIDPresentGNBID
	GNBIDPresentChoiceExtensions
)

type GNBID struct {
	Present          int
	GNBID            *aper.BitString `aper:"sizeLB:22,sizeUB:32"`
	ChoiceExtensions *ProtocolIESingleContainerGNBIDExtIEs
}
