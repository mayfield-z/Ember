package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type QosFlowSetupRequestItem struct {
	QosFlowIdentifier         QosFlowIdentifier
	QosFlowLevelQosParameters QosFlowLevelQosParameters                                `aper:"valueExt"`
	ERABID                    *ERABID                                                  `aper:"optional"`
	IEExtensions              *ProtocolExtensionContainerQosFlowSetupRequestItemExtIEs `aper:"optional"`
}
