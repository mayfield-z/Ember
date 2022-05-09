package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceFailedToSetupItemCxtFail struct {
	PDUSessionID                                PDUSessionID
	PDUSessionResourceSetupUnsuccessfulTransfer aper.OctetString
	IEExtensions                                *ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemCxtFailExtIEs `aper:"optional"`
}
