package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceHandoverItem struct {
	PDUSessionID            PDUSessionID
	HandoverCommandTransfer aper.OctetString
	IEExtensions            *ProtocolExtensionContainerPDUSessionResourceHandoverItemExtIEs `aper:"optional"`
}
