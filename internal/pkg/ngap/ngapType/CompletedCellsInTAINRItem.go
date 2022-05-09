package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type CompletedCellsInTAINRItem struct {
	NRCGI        NRCGI                                                      `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCompletedCellsInTAINRItemExtIEs `aper:"optional"`
}
