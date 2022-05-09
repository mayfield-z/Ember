package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type QosFlowNotifyItem struct {
	QosFlowIdentifier QosFlowIdentifier
	NotificationCause NotificationCause
	IEExtensions      *ProtocolExtensionContainerQosFlowNotifyItemExtIEs `aper:"optional"`
}
