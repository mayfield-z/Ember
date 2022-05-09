package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceToReleaseItemRelCmd struct {
	PDUSessionID                             PDUSessionID
	PDUSessionResourceReleaseCommandTransfer aper.OctetString
	IEExtensions                             *ProtocolExtensionContainerPDUSessionResourceToReleaseItemRelCmdExtIEs `aper:"optional"`
}
