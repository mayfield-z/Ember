package ue

import (
	"context"
	"fmt"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/ngap/ngapType"
	"github.com/looplab/fsm"
	"github.com/mayfield-z/ember/internal/pkg/gnb"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/timer"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	sqn         string
	pduSessions []utils.PDU
	ulDataRate  string
	dlDataRate  string
	//sm
	rrcFSM *fsm.FSM
	rmFSM  *fsm.FSM
	cmFSM  *fsm.FSM
	//timers
	t3346 *timer.Timer
	t3396 *timer.Timer
	t3444 *timer.Timer
	t3445 *timer.Timer
	t3502 *timer.Timer
	t3510 *timer.Timer
	t3511 *timer.Timer
	t3512 *timer.Timer
	t3516 *timer.Timer
	t3517 *timer.Timer
	t3519 *timer.Timer
	t3520 *timer.Timer
	t3521 *timer.Timer
	t3525 *timer.Timer
	t3540 *timer.Timer
	t3584 *timer.Timer
	t3585 *timer.Timer

	kamf         []uint8
	cipheringAlg uint8
	integrityAlg uint8
	knasEnc      [16]uint8
	knasInt      [16]uint8
	ULCount      security.Count
	DLCount      security.Count

	aMFRegionID uint8
	aMFPointer  uint8
	aMFSetID    uint16
	gGuti       [4]uint8

	id        uint8
	snn       string
	ctx       context.Context
	cancel    context.CancelFunc
	running   bool
	gnb       utils.UeGnb
	Notify    chan interface{}
	logger    *logrus.Entry
	nasLogger *logrus.Entry
}

func NewUE(supi string, mcc, mnc, key, op, opType, amf, ulDataRate, dlDataRate string, pduSessions []utils.PDU, id uint8, parent context.Context) *UE {
	// TODO: check dup
	mqueue.NewQueue(supi)
	ctx, cancelFunc := context.WithCancel(parent)
	return &UE{
		supi: supi,
		plmn: utils.PLMN{
			Mcc: mcc,
			Mnc: mnc,
		},
		key:          key,
		op:           op,
		opType:       opType,
		amf:          amf,
		sqn:          "0000000",
		cipheringAlg: security.AlgCiphering128NEA0,
		integrityAlg: security.AlgCiphering128NEA2,
		snn:          deriveSNN(mnc, mcc),
		ulDataRate:   ulDataRate,
		dlDataRate:   dlDataRate,
		pduSessions:  pduSessions,
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
		id:        id,
		ctx:       ctx,
		cancel:    cancelFunc,
		Notify:    make(chan interface{}, 1),
		logger:    logger.UeLog.WithFields(logrus.Fields{"name": supi}),
		nasLogger: logger.UeLog.WithFields(logrus.Fields{"name": supi, "part": "NAS"}),
	}
}

func (u *UE) NodeName() string {
	return fmt.Sprintf("ue-%v", u.supi)
}

func (u *UE) SUPI() string {
	return u.supi
}

func (u *UE) SetSUPI(supi string) {
	mqueue.DelQueue(u.supi)
	u.supi = supi
	u.logger = logger.UeLog.WithFields(logrus.Fields{"name": supi})
	u.nasLogger = logger.UeLog.WithFields(logrus.Fields{"name": supi, "part": "NAS"})
	mqueue.NewQueue(supi)
}

func (u *UE) Copy(supi string, id uint8) *UE {
	// TODO: check source ue state
	ue := *u
	ue.supi = supi
	u.id = id
	ue.logger = logger.UeLog.WithFields(logrus.Fields{"name": supi})
	ue.nasLogger = logger.UeLog.WithFields(logrus.Fields{"name": supi, "part": "NAS"})
	ue.Notify = make(chan interface{})
	mqueue.NewQueue(supi)
	return &ue
}

func (u *UE) getMessageChan() chan interface{} {
	return mqueue.GetMessageChan(u.supi)
}

func (u *UE) Run() {
	u.logger.Debugf("UE run")
	go u.messageHandler()
	u.running = true
}

func (u *UE) RRCSetupRequest(gnb *gnb.GNB) {
	if !u.running {
		u.logger.Errorf("UE %v not start but want rrc setup", u.supi)
	}
	if gnb.Running() {
		// for concurrent safety, don't change exec order
		u.gnb = utils.UeGnb{Name: gnb.NodeName()}
		msg := message.RRCSetupRequest{
			EstablishmentCause: ngapType.RRCEstablishmentCausePresentMoSignalling,
			SendBy:             u.supi,
		}
		mqueue.SendMessage(msg, gnb.NodeName())
	}
}

func deriveSNN(mnc, mcc string) string {
	// 5G:mnc093.mcc208.3gppnetwork.org
	var resu string
	if len(mnc) == 2 {
		resu = "5G:mnc0" + mnc + ".mcc" + mcc + ".3gppnetwork.org"
	} else {
		resu = "5G:mnc" + mnc + ".mcc" + mcc + ".3gppnetwork.org"
	}

	return resu
}
func ParseIpVersion(str string) (utils.IpVersion, error) {
	if str == "IPv4" {
		return utils.IPv4, nil
	} else if str == "IPv6" {
		return utils.IPv6, nil
	} else if str == "IPv4AndIPv6" {
		return utils.IPv4_AND_IPv6, nil
	}
	return utils.IPv4, errors.New(fmt.Sprintf("IpVersion \"%v\" can not parsed", str))
}
