package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type QosFlowItemWithDataForwarding struct {
	QosFlowIdentifier      QosFlowIdentifier
	DataForwardingAccepted *DataForwardingAccepted                                        `aper:"optional"`
	IEExtensions           *ProtocolExtensionContainerQosFlowItemWithDataForwardingExtIEs `aper:"optional"`
}
