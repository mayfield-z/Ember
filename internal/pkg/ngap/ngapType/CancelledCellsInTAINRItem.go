package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type CancelledCellsInTAINRItem struct {
	NRCGI              NRCGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCancelledCellsInTAINRItemExtIEs `aper:"optional"`
}
