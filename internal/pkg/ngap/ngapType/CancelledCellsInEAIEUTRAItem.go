package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type CancelledCellsInEAIEUTRAItem struct {
	EUTRACGI           EUTRACGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCancelledCellsInEAIEUTRAItemExtIEs `aper:"optional"`
}
