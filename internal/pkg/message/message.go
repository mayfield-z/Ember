package message

import (
	"github.com/free5gc/aper"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"time"
)

type RRCSetupRequest struct {
	// ue-Identity
	EstablishmentCause aper.Enumerated

	// not in spec, for convenience
	SendBy string
}

type RRCSetup struct {
	// implement when needed
}

type RRCReject struct{}

type RRCSetupComplete struct {
	NASRegistrationRequest []byte
	PLMN                   utils.PLMN
	SendBy                 string
}

type NASDownlinkPdu struct {
	PDU []byte
}

type NASUplinkPdu struct {
	PDU    []byte
	SendBy string
}

type GNBSetupSuccess struct{}

type GNBSetupReject struct{}

type UERRCSetupSuccess struct{}

type UERRCSetupReject struct{}

type UERegistrationSuccess struct{}

type UERegistrationReject struct{}

type UEPDUSessionEstablishmentAccept struct {
}

type UEPDUSessionEstablishmentReject struct {
}

type NodeType int

const (
	UE NodeType = iota
	GNB
)

type Event string

const ()

type StatusReport struct {
	NodeName string
	NodeType NodeType
	Event    interface{}
	Time     time.Time
}
