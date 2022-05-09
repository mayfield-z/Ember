package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type AdditionalDLUPTNLInformationForHOItem struct {
	AdditionalDLNGUUPTNLInformation        UPTransportLayerInformation `aper:"valueLB:0,valueUB:1"`
	AdditionalQosFlowSetupResponseList     QosFlowListWithDataForwarding
	AdditionalDLForwardingUPTNLInformation *UPTransportLayerInformation                                           `aper:"valueLB:0,valueUB:1,optional"`
	IEExtensions                           *ProtocolExtensionContainerAdditionalDLUPTNLInformationForHOItemExtIEs `aper:"optional"`
}
