package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type NRCGI struct {
	PLMNIdentity   PLMNIdentity
	NRCellIdentity NRCellIdentity
	IEExtensions   *ProtocolExtensionContainerNRCGIExtIEs `aper:"optional"`
}
