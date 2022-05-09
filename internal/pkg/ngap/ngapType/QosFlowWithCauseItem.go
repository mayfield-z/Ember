package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type QosFlowWithCauseItem struct {
	QosFlowIdentifier QosFlowIdentifier
	Cause             Cause                                                 `aper:"valueLB:0,valueUB:5"`
	IEExtensions      *ProtocolExtensionContainerQosFlowWithCauseItemExtIEs `aper:"optional"`
}
