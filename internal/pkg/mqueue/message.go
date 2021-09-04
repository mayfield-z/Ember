package mqueue

import (
	"github.com/free5gc/aper"
	"github.com/mayfield-z/ember/internal/pkg/utils"
)

type RRCSetupRequestMessage struct {
	// ue-Identity
	EstablishmentCause aper.Enumerated

	// not in spec, for convenience
	SendBy string
}

type RRCSetupMessage struct {
	// implement when needed
}

type RRCRejectMessage struct{}

type RRCSetupCompleteMessage struct {
	NASRegistrationRequest []byte
	PLMN                   utils.PLMN
	SendBy                 string
}
