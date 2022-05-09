package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type QosFlowSetupResponseItemSURes struct {
	QosFlowIdentifier QosFlowIdentifier
	IEExtensions      *ProtocolExtensionContainerQosFlowSetupResponseItemSUResExtIEs `aper:"optional"`
}
