package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type SliceSupportItem struct {
	SNSSAI       SNSSAI                                            `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerSliceSupportItemExtIEs `aper:"optional"`
}
