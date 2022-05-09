package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type AssociatedQosFlowItem struct {
	QosFlowIdentifier        QosFlowIdentifier
	QosFlowMappingIndication *aper.Enumerated                                       `aper:"optional,valueExt,valueLB:0,valueUB:1"`
	IEExtensions             *ProtocolExtensionContainerAssociatedQosFlowItemExtIEs `aper:"optional"`
}
