package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type XnExtTLAItem struct {
	IPsecTLA     *TransportLayerAddress                        `aper:"optional"`
	GTPTLAs      *XnGTPTLAs                                    `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerXnExtTLAItemExtIEs `aper:"optional"`
}
