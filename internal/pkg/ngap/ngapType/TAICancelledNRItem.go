package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type TAICancelledNRItem struct {
	TAI                   TAI `aper:"valueExt"`
	CancelledCellsInTAINR CancelledCellsInTAINR
	IEExtensions          *ProtocolExtensionContainerTAICancelledNRItemExtIEs `aper:"optional"`
}
