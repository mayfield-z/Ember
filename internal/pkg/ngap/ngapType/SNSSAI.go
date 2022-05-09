package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type SNSSAI struct {
	SST          SST
	SD           *SD                                     `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerSNSSAIExtIEs `aper:"optional"`
}
