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
	"net"
	"sync"
	"time"
)

const (
	stateRRCConnected = "RRC_CONNECTED"
	stateRRCInactive  = "RRC_INACTIVE"
	stateRRCIdle      = "RRC_IDLE"

	eventRRCSetup             = "RRC_SETUP_EVENT"
	eventRRCConnectionRelease = "RRC_CONNECTION_RELEASE_EVENT"

	stateRMDeregistered = "RM_DEREGISTERED"
	stateRMRegistered   = "RM_REGISTERED"

	eventRMRegistrationAccept       = "RM_REGISTRATION_ACCEPT_EVENT"
	eventRMRegistrationReject       = "RM_REGISTRATION_REJECT_EVENT"
	eventRMDeregistration           = "RM_DEREGISTRATION_EVENT"
	eventRMRegistrationUpdateAccept = "RM_REGISTRATION_UPDATE_ACCEPT_EVENT"

	stateCMIdle      = "CM_IDLE"
	stateCMConnceted = "CM_CONNECTED"

	stateSMPDUSessionInactive            = "PDU_SESSION_INACTIVE"
	stateSMPDUSessionInactivePending     = "PDU_SESSION_INACTIVE_PENDING"
	stateSMPDUSessionActive              = "PDU_SESSION_ACTIVE"
	stateSMPDUSessionActivePending       = "PDU_SESSION_ACITVE_PENDING"
	stateSMPDUSessionModificationPending = "PDU_SESSION_MODIFICATION_PENDING"

	eventSMPDUSessionEstablishmentRequest = "PDU_SESSION_ESTABLISHMENT_REQUEST"
	eventSMPDUSessionEstablishmentAccept  = "PDU_SESSION_ESTABLISHMENT_ACCEPT"

	stateMMDeregistered            = "5GMM_DEREGISTERED"
	stateMMRegisteredInitiated     = "5GMM_REGISTERED_INITIATED"
	stateMMRegistered              = "5GMM_REGISTERED"
	stateMMDeregisteredInitated    = "5GMM_DEREGISTERED_INITIATED"
	stateMMServiceRequestInitiated = "5GMM_SERVICE_REQUEST_INITIATED"

	//eventMM
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
	smFSM  *fsm.FSM
	mmFSM  *fsm.FSM
	//cmFSM  *fsm.FSM
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

	ip net.IP

	id                  uint64
	snn                 string
	parentCtx           context.Context
	ctx                 context.Context
	cancel              context.CancelFunc
	running             bool
	gnb                 utils.UeGnb
	statusReportChannel chan message.StatusReport
	logger              *logrus.Entry
	nasLogger           *logrus.Entry
}

func NewUE(supi string, mcc, mnc, key, op, opType, amf, ulDataRate, dlDataRate string, pduSessions []utils.PDU, id uint64, parent context.Context) *UE {
	// TODO: check dup
	mqueue.NewQueue(supi)
	ctx, cancelFunc := context.WithCancel(parent)
	u := UE{
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
		pduSessions:  make([]utils.PDU, 1),
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
		smFSM: fsm.NewFSM(
			stateSMPDUSessionInactive,
			fsm.Events{
				{Name: eventSMPDUSessionEstablishmentRequest, Src: []string{stateSMPDUSessionInactive}, Dst: stateSMPDUSessionActivePending},
				{Name: eventSMPDUSessionEstablishmentAccept, Src: []string{stateSMPDUSessionActivePending}, Dst: stateSMPDUSessionActive},
			},
			nil,
		),
		rrcFSM: fsm.NewFSM(
			stateRRCIdle,
			fsm.Events{
				{Name: eventRRCSetup, Src: []string{stateRRCIdle}, Dst: stateRRCConnected},
				{Name: eventRRCConnectionRelease, Src: []string{stateRRCConnected, stateRRCIdle}, Dst: stateRRCIdle},
			},
			nil,
		),
		//cmFSM: fsm.NewFSM(
		//	stateCMIdle,
		//	nil,
		//	nil,
		//),
		id:                  id,
		parentCtx:           parent,
		ctx:                 ctx,
		cancel:              cancelFunc,
		statusReportChannel: make(chan message.StatusReport, 1),
		logger:              logger.UeLog.WithFields(logrus.Fields{"name": supi}),
		nasLogger:           logger.UeLog.WithFields(logrus.Fields{"name": supi, "part": "NAS"}),
	}
	copy(u.pduSessions[:], pduSessions)
	return &u
}

func (u *UE) NodeName() string {
	return fmt.Sprintf("ue-%v", u.supi)
}

func (u *UE) GetSUPI() string {
	return u.supi
}

func (u *UE) SetSUPI(supi string) {
	mqueue.DelQueue(u.supi)
	u.supi = supi
	u.logger = logger.UeLog.WithFields(logrus.Fields{"name": supi})
	u.nasLogger = logger.UeLog.WithFields(logrus.Fields{"name": supi, "part": "NAS"})
	mqueue.NewQueue(supi)
}

func (u *UE) GetIP() net.IP {
	return u.ip
}

func (u *UE) GetID() uint64 {
	return u.id
}

func (u *UE) Copy(supi string, ueId uint64, sessionId uint8) *UE {
	// TODO: check source ue state
	ue := *u
	ue.supi = supi
	ue.id = ueId
	for _, session := range ue.pduSessions {
		session.Id = sessionId
		sessionId += 1
	}
	ue.logger = logger.UeLog.WithFields(logrus.Fields{"name": supi})
	ue.nasLogger = logger.UeLog.WithFields(logrus.Fields{"name": supi, "part": "NAS"})
	ue.statusReportChannel = make(chan message.StatusReport, 1)
	ue.ctx, ue.cancel = context.WithCancel(ue.parentCtx)
	mqueue.NewQueue(supi)
	return &ue
}

func (u *UE) getMessageChan() chan interface{} {
	return mqueue.GetMessageChannel(u.supi)
}

func (u *UE) GetPDUSessionNum() uint8 {
	return uint8(len(u.pduSessions))
}

func (u *UE) GetRMState() string {
	return u.rmFSM.Current()
}

func (u *UE) GetMMState() string {
	return u.mmFSM.Current()
}

func (u *UE) GetSMState() string {
	return u.smFSM.Current()
}

func (u *UE) Run() {
	u.logger.Debugf("UE run")
	go u.messageHandler()
	u.running = true
}

func (u *UE) Running() bool {
	return u.running
}

func (u *UE) Stop(wg *sync.WaitGroup) {
	u.logger.Debugf("UE stop")
	u.cancel()
	u.running = false
	wg.Done()
}

func (u *UE) SendStatusReport(event message.Event) {
	statusReport := message.StatusReport{
		NodeName: u.supi,
		NodeType: message.GNB,
		Event:    event,
		Time:     time.Now(),
	}
	u.statusReportChannel <- statusReport
	if r := mqueue.GetQueue("reporter"); r != nil {
		mqueue.SendMessage(statusReport, "reporter")
	}
}

func (u *UE) StatusReport() <-chan message.StatusReport {
	return u.statusReportChannel
}

func (u *UE) EstablishPDUSession(pduSessionNumber uint8) {
	u.logger.Tracef("EstablishPDUSession(%v)", pduSessionNumber)
	if pduSessionNumber >= u.GetPDUSessionNum() {
		u.logger.Errorf("chose No.%v session, only have %v sessions", pduSessionNumber, u.GetPDUSessionNum())
		return
	}
	data, err := u.buildULNasTransportPDUSessionEstablishmentRequest(pduSessionNumber)
	if err != nil {
		u.logger.Errorf("Establish PDU Session failed: %v", err)
		return
	}
	mqueue.SendMessage(message.NASUplinkPdu{PDU: data, SendBy: u.supi}, u.gnb.Name)
	u.smFSM.Event(eventSMPDUSessionEstablishmentRequest)
}

func (u *UE) RRCSetupRequest(gnb *gnb.GNB) {
	u.logger.Tracef("RRCSetupRequest(%v)", gnb.NodeName())
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
