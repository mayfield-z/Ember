package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type EmergencyAreaIDCancelledNRItem struct {
	EmergencyAreaID       EmergencyAreaID
	CancelledCellsInEAINR CancelledCellsInEAINR
	IEExtensions          *ProtocolExtensionContainerEmergencyAreaIDCancelledNRItemExtIEs `aper:"optional"`
}
