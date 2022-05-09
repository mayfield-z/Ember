package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	WarningAreaListPresentNothing int = iota /* No components present */
	WarningAreaListPresentEUTRACGIListForWarning
	WarningAreaListPresentNRCGIListForWarning
	WarningAreaListPresentTAIListForWarning
	WarningAreaListPresentEmergencyAreaIDList
	WarningAreaListPresentChoiceExtensions
)

type WarningAreaList struct {
	Present                int
	EUTRACGIListForWarning *EUTRACGIListForWarning
	NRCGIListForWarning    *NRCGIListForWarning
	TAIListForWarning      *TAIListForWarning
	EmergencyAreaIDList    *EmergencyAreaIDList
	ChoiceExtensions       *ProtocolIESingleContainerWarningAreaListExtIEs
}
