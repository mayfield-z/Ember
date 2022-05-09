package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

type PDUSessionResourceSecondaryRATUsageItem struct {
	PDUSessionID                        PDUSessionID
	SecondaryRATDataUsageReportTransfer aper.OctetString
	IEExtensions                        *ProtocolExtensionContainerPDUSessionResourceSecondaryRATUsageItemExtIEs `aper:"optional"`
}
