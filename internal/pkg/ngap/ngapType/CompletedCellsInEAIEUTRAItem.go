package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type CompletedCellsInEAIEUTRAItem struct {
	EUTRACGI     EUTRACGI                                                      `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCompletedCellsInEAIEUTRAItemExtIEs `aper:"optional"`
}
