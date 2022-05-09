package ngapType

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceModifyResponseTransfer struct {
	DLNGUUPTNLInformation                *UPTransportLayerInformation                                              `aper:"valueLB:0,valueUB:1,optional"`
	ULNGUUPTNLInformation                *UPTransportLayerInformation                                              `aper:"valueLB:0,valueUB:1,optional"`
	QosFlowAddOrModifyResponseList       *QosFlowAddOrModifyResponseList                                           `aper:"optional"`
	AdditionalDLQosFlowPerTNLInformation *QosFlowPerTNLInformationList                                             `aper:"optional"`
	QosFlowFailedToAddOrModifyList       *QosFlowListWithCause                                                     `aper:"optional"`
	IEExtensions                         *ProtocolExtensionContainerPDUSessionResourceModifyResponseTransferExtIEs `aper:"optional"`
}
