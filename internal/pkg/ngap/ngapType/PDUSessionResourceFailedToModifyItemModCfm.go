package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceFailedToModifyItemModCfm struct {
	PDUSessionID                                           PDUSessionID
	PDUSessionResourceModifyIndicationUnsuccessfulTransfer aper.OctetString
	IEExtensions                                           *ProtocolExtensionContainerPDUSessionResourceFailedToModifyItemModCfmExtIEs `aper:"optional"`
}
