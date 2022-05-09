package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type QoSFlowsUsageReportItem struct {
	QosFlowIdentifier       QosFlowIdentifier
	RATType                 aper.Enumerated
	QoSFlowsTimedReportList VolumeTimedReportList
	IEExtensions            *ProtocolExtensionContainerQoSFlowsUsageReportItemExtIEs `aper:"optional"`
}
