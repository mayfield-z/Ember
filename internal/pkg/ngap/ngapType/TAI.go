package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type TAI struct {
	PLMNIdentity PLMNIdentity
	TAC          TAC
	IEExtensions *ProtocolExtensionContainerTAIExtIEs `aper:"optional"`
}
