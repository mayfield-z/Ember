package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type UEAssociatedLogicalNGConnectionItem struct {
	AMFUENGAPID  *AMFUENGAPID                                                         `aper:"optional"`
	RANUENGAPID  *RANUENGAPID                                                         `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerUEAssociatedLogicalNGConnectionItemExtIEs `aper:"optional"`
}
