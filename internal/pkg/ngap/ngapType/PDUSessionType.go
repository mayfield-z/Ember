package ngapType

import "github.com/mayfield-z/ember/internal/pkg/aper"

// Need to import "github.com/mayfield-z/ember/internal/pkg/aper" if it uses "aper"

const (
	PDUSessionTypePresentIpv4         aper.Enumerated = 0
	PDUSessionTypePresentIpv6         aper.Enumerated = 1
	PDUSessionTypePresentIpv4v6       aper.Enumerated = 2
	PDUSessionTypePresentEthernet     aper.Enumerated = 3
	PDUSessionTypePresentUnstructured aper.Enumerated = 4
)

type PDUSessionType struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:4"`
}
