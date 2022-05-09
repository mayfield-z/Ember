package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type AreaOfInterestCellItem struct {
	NGRANCGI     NGRANCGI                                                `aper:"valueLB:0,valueUB:2"`
	IEExtensions *ProtocolExtensionContainerAreaOfInterestCellItemExtIEs `aper:"optional"`
}
