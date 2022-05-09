package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type CellIDCancelledNRItem struct {
	NRCGI              NRCGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCellIDCancelledNRItemExtIEs `aper:"optional"`
}
