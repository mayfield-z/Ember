package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	UPTransportLayerInformationPresentNothing int = iota /* No components present */
	UPTransportLayerInformationPresentGTPTunnel
	UPTransportLayerInformationPresentChoiceExtensions
)

type UPTransportLayerInformation struct {
	Present          int
	GTPTunnel        *GTPTunnel `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerUPTransportLayerInformationExtIEs
}
