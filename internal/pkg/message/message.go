package message

import (
	"github.com/mayfield-z/ember/internal/pkg/aper"
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

type NodeType int

const (
	UE NodeType = iota
	GNB
	Controller
)

type Event int

//TODO: more events
const (
	GNBSetupSuccess Event = iota
	GNBSetupReject
	UERRCSetupSuccess
	UERRCSetupReject
	UERegistrationSuccess
	UERegistrationReject
	UEPDUSessionEstablishmentAccept
	UEPDUSessionEstablishmentReject
	ControllerStart
	ControllerStop
	EmulateUE
	EmulateGNB
)

type StatusReport struct {
	NodeName string
	NodeType NodeType
	Event    Event
	Time     time.Time
}
