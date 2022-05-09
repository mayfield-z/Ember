package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceModifyItemModReq struct {
	PDUSessionID                            PDUSessionID
	NASPDU                                  *NASPDU `aper:"optional"`
	PDUSessionResourceModifyRequestTransfer aper.OctetString
	IEExtensions                            *ProtocolExtensionContainerPDUSessionResourceModifyItemModReqExtIEs `aper:"optional"`
}
