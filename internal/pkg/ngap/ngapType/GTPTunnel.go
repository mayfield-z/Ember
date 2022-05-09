package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type GTPTunnel struct {
	TransportLayerAddress TransportLayerAddress
	GTPTEID               GTPTEID
	IEExtensions          *ProtocolExtensionContainerGTPTunnelExtIEs `aper:"optional"`
}
