package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceNotifyItem struct {
	PDUSessionID                     PDUSessionID
	PDUSessionResourceNotifyTransfer aper.OctetString
	IEExtensions                     *ProtocolExtensionContainerPDUSessionResourceNotifyItemExtIEs `aper:"optional"`
}
