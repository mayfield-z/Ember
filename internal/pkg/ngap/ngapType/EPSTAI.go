package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type EPSTAI struct {
	PLMNIdentity PLMNIdentity
	EPSTAC       EPSTAC
	IEExtensions *ProtocolExtensionContainerEPSTAIExtIEs `aper:"optional"`
}
