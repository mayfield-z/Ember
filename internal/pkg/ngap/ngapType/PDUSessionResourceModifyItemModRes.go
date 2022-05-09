package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceModifyItemModRes struct {
	PDUSessionID                             PDUSessionID
	PDUSessionResourceModifyResponseTransfer aper.OctetString
	IEExtensions                             *ProtocolExtensionContainerPDUSessionResourceModifyItemModResExtIEs `aper:"optional"`
}
