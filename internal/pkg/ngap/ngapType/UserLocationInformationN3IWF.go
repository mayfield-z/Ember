package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type UserLocationInformationN3IWF struct {
	IPAddress    TransportLayerAddress
	PortNumber   PortNumber
	IEExtensions *ProtocolExtensionContainerUserLocationInformationN3IWFExtIEs `aper:"optional"`
}
