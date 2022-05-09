package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type TAIBroadcastEUTRAItem struct {
	TAI                      TAI `aper:"valueExt"`
	CompletedCellsInTAIEUTRA CompletedCellsInTAIEUTRA
	IEExtensions             *ProtocolExtensionContainerTAIBroadcastEUTRAItemExtIEs `aper:"optional"`
}
