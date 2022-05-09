package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceToBeSwitchedDLItem struct {
	PDUSessionID              PDUSessionID
	PathSwitchRequestTransfer aper.OctetString
	IEExtensions              *ProtocolExtensionContainerPDUSessionResourceToBeSwitchedDLItemExtIEs `aper:"optional"`
}
