package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type TAIBroadcastNRItem struct {
	TAI                   TAI `aper:"valueExt"`
	CompletedCellsInTAINR CompletedCellsInTAINR
	IEExtensions          *ProtocolExtensionContainerTAIBroadcastNRItemExtIEs `aper:"optional"`
}
