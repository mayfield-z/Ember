package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type QosFlowAddOrModifyResponseItem struct {
	QosFlowIdentifier QosFlowIdentifier
	IEExtensions      *ProtocolExtensionContainerQosFlowAddOrModifyResponseItemExtIEs `aper:"optional"`
}
