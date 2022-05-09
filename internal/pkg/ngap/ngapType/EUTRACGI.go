package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type EUTRACGI struct {
	PLMNIdentity      PLMNIdentity
	EUTRACellIdentity EUTRACellIdentity
	IEExtensions      *ProtocolExtensionContainerEUTRACGIExtIEs `aper:"optional"`
}
