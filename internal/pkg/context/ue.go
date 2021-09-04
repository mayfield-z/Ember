package context

import (
	"github.com/looplab/fsm"
	"github.com/mayfield-z/ember/internal/pkg/timer"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/pkg/errors"
)

const (
	stateRRCConnected   = "RRC-CONNECTED"
	stateRRCInactive    = "RRC-INACTIVE"
	stateRRCIdle        = "RRC-IDLE"
	stateRMDeregistered = "RM-DEREGISTERED"
	stateRMRegistered   = "RM-REGISTERED"
	stateCMIdle         = "CM-IDLE"
	stateCMConnceted    = "CM-CONNECTED"

	eventRRCConnect                 = "RRC-CONNECT-EVENT"
	eventRMRegistrationAccept       = "RM-REGISTRATION-ACCEPT-EVENT"
	eventRMRegistrationReject       = "RM-REGISTRATION-REJECT-EVENT"
	eventRMDeregistration           = "RM-DEREGISTRATION-EVENT"
	eventRMRegistrationUpdateAccept = "RM-REGISTRATION-UPDATE-ACCEPT-EVENT"
)

type UE struct {
	supi        string
	plmn        utils.PLMN
	key         string
	op          string
	opType      string
	amf         string
	pduSessions []PDU
	//sm
	rrcFSM *fsm.FSM
	rmFSM  *fsm.FSM
	cmFSM  *fsm.FSM
	//timers
	T3346 *timer.Timer
	T3396 *timer.Timer
	T3444 *timer.Timer
	T3445 *timer.Timer
	T3502 *timer.Timer
	T3510 *timer.Timer
	T3511 *timer.Timer
	T3512 *timer.Timer
	T3516 *timer.Timer
	T3517 *timer.Timer
	T3519 *timer.Timer
	T3520 *timer.Timer
	T3521 *timer.Timer
	T3525 *timer.Timer
	T3540 *timer.Timer
	T3584 *timer.Timer
	T3585 *timer.Timer

	connectedGnb *GNB
}

func NewUE(supi string, mcc, mnc string, key, op, opType, amf string, pduSessions []PDU) *UE {
	return &UE{
		supi: supi,
		plmn: utils.PLMN{
			Mcc: mcc,
			Mnc: mnc,
		},
		key:         key,
		op:          op,
		opType:      opType,
		amf:         amf,
		pduSessions: pduSessions,
		rrcFSM: fsm.NewFSM(
			stateRRCIdle,
			fsm.Events{
				{Name: eventRRCConnect, Src: []string{stateRRCIdle}, Dst: stateRRCConnected},
			},
			nil,
		),
		rmFSM: fsm.NewFSM(
			stateRMDeregistered,
			fsm.Events{
				{Name: eventRMRegistrationAccept, Src: []string{stateRMDeregistered}, Dst: stateRMRegistered},
				{Name: eventRMRegistrationReject, Src: []string{stateRMRegistered, stateRMDeregistered}, Dst: stateRMDeregistered},
				{Name: eventRMDeregistration, Src: []string{stateRMRegistered}, Dst: stateRMDeregistered},
				{Name: eventRMRegistrationUpdateAccept, Src: []string{stateRMRegistered}, Dst: stateRMRegistered},
			},
			nil,
		),
		cmFSM: fsm.NewFSM(
			stateCMIdle,
			nil,
			nil,
		),
	}
}

func (u *UE) SUPI() string {
	return u.supi
}

func (u *UE) rRCSetup(gnb *GNB) {
	if gnb.running {
		u.rrcFSM.Event(eventRRCConnect)
		u.connectedGnb = gnb
	}
}

func (u *UE) Register() error {
	gnb := u.connectedGnb
	if !gnb.Connected() {
		return errors.WithMessage(u.UeRegisterError(), "gnb not connected to amf")
	}
	_, err := gnb.InitialUE(u)
	if err != nil {
		return errors.WithMessage(err, "")
	}
	return nil
}

type IpVersion int

const (
	IPv4 IpVersion = iota
	IPv6
	IPv4_AND_IPv6
)

type PDU struct {
	IpType IpVersion
	Apn    string
	Nssai  utils.SNSSAI
}
