package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type EmergencyAreaIDBroadcastNRItem struct {
	EmergencyAreaID       EmergencyAreaID
	CompletedCellsInEAINR CompletedCellsInEAINR
	IEExtensions          *ProtocolExtensionContainerEmergencyAreaIDBroadcastNRItemExtIEs `aper:"optional"`
}
