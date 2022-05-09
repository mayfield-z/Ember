package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type TNLAssociationItem struct {
	TNLAssociationAddress CPTransportLayerInformation                         `aper:"valueLB:0,valueUB:1"`
	Cause                 Cause                                               `aper:"valueLB:0,valueUB:5"`
	IEExtensions          *ProtocolExtensionContainerTNLAssociationItemExtIEs `aper:"optional"`
}
