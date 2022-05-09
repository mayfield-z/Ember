package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceSetupItemHOReq struct {
	PDUSessionID            PDUSessionID
	SNSSAI                  SNSSAI `aper:"valueExt"`
	HandoverRequestTransfer aper.OctetString
	IEExtensions            *ProtocolExtensionContainerPDUSessionResourceSetupItemHOReqExtIEs `aper:"optional"`
}
