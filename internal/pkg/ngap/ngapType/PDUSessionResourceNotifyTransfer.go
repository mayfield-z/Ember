package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceNotifyTransfer struct {
	QosFlowNotifyList   *QosFlowNotifyList                                                `aper:"optional"`
	QosFlowReleasedList *QosFlowListWithCause                                             `aper:"optional"`
	IEExtensions        *ProtocolExtensionContainerPDUSessionResourceNotifyTransferExtIEs `aper:"optional"`
}
