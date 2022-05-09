package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceFailedToModifyItemModRes struct {
	PDUSessionID                                 PDUSessionID
	PDUSessionResourceModifyUnsuccessfulTransfer aper.OctetString
	IEExtensions                                 *ProtocolExtensionContainerPDUSessionResourceFailedToModifyItemModResExtIEs `aper:"optional"`
}
