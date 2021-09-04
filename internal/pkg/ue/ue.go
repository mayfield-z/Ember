package ue

import (
	"context"
	"fmt"
	"github.com/free5gc/ngap/ngapType"
	"github.com/looplab/fsm"
	"github.com/mayfield-z/ember/internal/pkg/gnb"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/timer"
	"github.com/mayfield-z/ember/internal/pkg/utils"
)

const (
	stateRRCConnected   = "RRC-CONNECTED"
	stateRRCInactive    = "RRC-INACTIVE"
	stateRRCIdle        = "RRC-IDLE"
	stateRMDeregistered = "RM-DEREGISTERED"
	stateRMRegistered   = "RM-REGISTERED"
	stateCMIdle         = "CM-IDLE"
	stateCMConnceted    = "CM-CONNECTED"

	eventRRCSetup                   = "RRC-SETUP-EVENT"
	eventRRCConnectionRelease       = "RRC-CONNECTION-RELEASE-EVENT"
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

	ctx     context.Context
	cancel  context.CancelFunc
	running bool
	gnb     UEGNB
}

func NewUE(supi string, mcc, mnc string, key, op, opType, amf string, pduSessions []PDU) *UE {
	// TODO: check dup
	mqueue.NewQueue(supi)
	// TODO: change context
	ctx, cancelFunc := context.WithCancel(context.Background())
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
				{Name: eventRRCSetup, Src: []string{stateRRCIdle}, Dst: stateRRCConnected},
				{Name: eventRRCConnectionRelease, Src: []string{stateRRCConnected, stateRRCIdle}, Dst: stateRRCIdle},
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
		ctx:    ctx,
		cancel: cancelFunc,
	}
}

func (u *UE) NodeName() string {
	return fmt.Sprintf("ue-%v", u.supi)
}

func (u *UE) SUPI() string {
	return u.supi
}

func (u *UE) getMessageChan() chan interface{} {
	return mqueue.GetMessageChan(u.supi)
}

func (u *UE) Run() {
	go u.messageHandler()
	u.running = true
}

func (u *UE) RRCSetupRequest(gnb *gnb.GNB) {
	if !u.running {
		logger.UeLog.Errorf("UE %v not start but want rrc setup", u.supi)
	}
	if gnb.Running() {
		msg := mqueue.RRCSetupRequestMessage{
			EstablishmentCause: ngapType.RRCEstablishmentCausePresentMoSignalling,
			SendBy:             u.supi,
		}
		mqueue.SendMessage(msg, gnb.NodeName())
		u.gnb = UEGNB{Name: gnb.NodeName()}
	}
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

type UEGNB struct {
	Name string
}
