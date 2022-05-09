package gnb

import (
	"context"
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

const (
	readBufSize    = 8192
	rANUENGAPIDMax = int64(1<<32 - 1)
)

// GNB only support one amf and one slice now
type GNB struct {
	// need to read from config file
	name            string
	globalRANNodeID uint32
	plmn            utils.PLMN
	nci             uint64 // encoded as 36-bit string
	idLength        uint8
	tac             uint32 // encoded as 2 or 3-octet string
	n2Address       net.IP
	n3Address       net.IP
	amfAddress      net.IP
	amfPort         int
	snssai          utils.SNSSAI

	// auto
	ueMapBySupi             sync.Map //map[string]*utils.GnbUe
	ueMapByRANUENGAPID      sync.Map //map[int64]*utils.GnbUe
	allUE                   []*utils.GnbUe
	allUEMutex              sync.Mutex
	rANUENGAPIDPointer      int64
	rANUENGAPIDPointerMutex sync.Mutex
	running                 bool
	gnbAmf                  utils.GnbAmf
	sctpConn                *sctp.SCTPConn
	statusReportChannel     chan message.StatusReport
	logger                  *logrus.Entry
	ctx                     context.Context
	cancelFunc              context.CancelFunc
}

func NewGNB(name string, globalRANNodeID uint32, mcc, mnc string, nci uint64, tac uint32, idLength uint8, n2Address, n3Address, amfAddress net.IP, amfPort int, sst uint8, sd uint32, parent context.Context) *GNB {
	//TODO: check if same name gnb exists
	mqueue.NewQueue(name)
	ctx, cancelFunc := context.WithCancel(parent)
	g := GNB{
		name:            name,
		globalRANNodeID: globalRANNodeID,
		plmn: utils.PLMN{
			Mcc: mcc,
			Mnc: mnc,
		},
		nci:        nci,
		idLength:   idLength,
		tac:        tac,
		amfPort:    amfPort,
		n2Address:  make([]byte, 16),
		n3Address:  make([]byte, 16),
		amfAddress: make([]byte, 16),
		snssai: utils.SNSSAI{
			Sst: sst,
			Sd:  sd,
		},
		running:             false,
		sctpConn:            nil,
		statusReportChannel: make(chan message.StatusReport, 1),
		logger:              logger.GnbLog.WithFields(logrus.Fields{"name": name}),
		ctx:                 ctx,
		cancelFunc:          cancelFunc,
	}
	copy(g.n2Address[:], n2Address)
	copy(g.n3Address[:], n3Address)
	copy(g.amfAddress[:], amfAddress)
	return &g
}

func (g *GNB) NodeName() string {
	return g.name
}

func (g *GNB) Running() bool {
	return g.running
}

func (g *GNB) Connected() bool {
	return g.gnbAmf.Connected
}

func (g *GNB) SetName(name string) *GNB {
	mqueue.DelQueue(g.name)
	g.name = name
	g.logger = logger.GnbLog.WithFields(logrus.Fields{"name": name})
	mqueue.NewQueue(g.name)
	return g
}

func (g *GNB) Copy(name string) *GNB {
	gnb := *g
	gnb.name = name
	gnb.logger = logger.GnbLog.WithFields(logrus.Fields{"name": name})
	gnb.statusReportChannel = make(chan message.StatusReport, 1)
	gnb.ueMapBySupi = sync.Map{}
	gnb.ueMapByRANUENGAPID = sync.Map{}
	gnb.rANUENGAPIDPointerMutex = sync.Mutex{}
	gnb.sctpConn = &sctp.SCTPConn{}
	gnb.running = false
	mqueue.NewQueue(name)
	return &gnb
}

func (g *GNB) SetGlobalRANNodeID(id uint32) *GNB {
	g.globalRANNodeID = id
	return g
}

func (g *GNB) SetNRCellIdentity(nci uint64) *GNB {
	g.nci = nci
	return g
}

func (g *GNB) SetAMFAddress(ip net.IP) *GNB {
	copy(g.amfAddress[:], ip)
	return g
}

func (g *GNB) SetN2Addresses(ip net.IP) *GNB {
	copy(g.n2Address[:], ip)
	return g
}

func (g *GNB) GetN2Addresses() net.IP {
	return g.n2Address
}

func (g *GNB) SetN3Addresses(ip net.IP) *GNB {
	copy(g.n3Address[:], ip)
	return g
}

func (g *GNB) GetN3Addresses() net.IP {
	return g.n3Address
}

func (g *GNB) getMessageChan() chan interface{} {
	return mqueue.GetMessageChannel(g.name)
}

func (g *GNB) FindUEBySUPI(supi string) *utils.GnbUe {
	if ue, ok := g.ueMapBySupi.Load(supi); ok {
		return ue.(*utils.GnbUe)
	}
	return nil
}

func (g *GNB) FindUEByRANUENGAPID(id int64) *utils.GnbUe {
	if ue, ok := g.ueMapByRANUENGAPID.Load(id); ok {
		return ue.(*utils.GnbUe)
	}
	return nil
}

func (g *GNB) GetAllUE() []*utils.GnbUe {
	var f func(key interface{}, value interface{}) bool
	f = func(key interface{}, value interface{}) bool {
		g.allUE = append(g.allUE, value.(*utils.GnbUe))
		return true
	}
	g.allUEMutex.Lock()
	defer g.allUEMutex.Unlock()
	g.ueMapBySupi.Range(f)
	return g.allUE
}

func (g *GNB) Run() error {
	g.running = true
	g.logger.Debugf("GNB %s is running, N2: %v, N3: %v", g.name, g.n2Address.String(), g.n3Address.String())
	g.logger.Debugf("SCTP dial %v:%v", g.amfAddress, g.amfPort)
	conn, err := Dial("sctp", g.n2Address, g.amfAddress, g.amfPort)
	if err != nil {
		g.Stop()
		return errors.Wrapf(err, "Failed to dial sctp address %v:%v", g.amfAddress, g.amfPort)
	}
	g.sctpConn = conn
	go g.sctpHandler(readBufSize)
	go g.messageHandler()
	err = g.connectToAmf()
	if err != nil {
		g.Stop()
		return errors.Wrapf(err, "Send NGSetupRequestPDU error.")
	}
	return nil
}

func (g *GNB) Stop() {
	g.cancelFunc()
	g.sctpConn.Close()
	g.running = false
	g.gnbAmf.Connected = false
	g.logger.Debugf("GNB stop")
}

func (g *GNB) connectToAmf() error {
	err := g.sendNGSetupRequestPDU()
	if err != nil {
		return errors.Wrap(err, "Connect to AMF failed")
	}
	return nil
}

func (g *GNB) sendNGSetupRequestPDU() error {
	nGSetupRequest, err := g.buildNGSetupRequest()
	if err != nil {
		return errors.Wrapf(err, "NGSetupRequestPDU build failed")
	}
	_, err = g.sctpConn.Write(nGSetupRequest)
	if err != nil {
		return errors.Wrapf(err, "NGSetupRequestPDU send error")
	}
	return nil
}

func (g *GNB) allocateRANUENGAPID() int64 {
	rANUENGAPID := int64(-1)
	g.rANUENGAPIDPointerMutex.Lock()
	defer g.rANUENGAPIDPointerMutex.Unlock()
	for i := g.rANUENGAPIDPointer; i < rANUENGAPIDMax; i++ {
		if g.FindUEByRANUENGAPID(i) == nil {
			rANUENGAPID = i
			g.rANUENGAPIDPointer += 1
			break
		}
	}
	if rANUENGAPID == int64(-1) {
		for i := int64(0); i < g.rANUENGAPIDPointer; i++ {
			if g.FindUEByRANUENGAPID(i) == nil {
				rANUENGAPID = i
				g.rANUENGAPIDPointer += 1
				break
			}
		}
	}
	return rANUENGAPID
}

func (g *GNB) SendStatusReport(event message.Event) {
	statusReport := message.StatusReport{
		NodeName: g.name,
		NodeType: message.GNB,
		Event:    event,
		Time:     time.Now(),
	}
	g.statusReportChannel <- statusReport
	if r := mqueue.GetQueue("reporter"); r != nil {
		mqueue.SendMessage(statusReport, "reporter")
	}
}

func (g *GNB) StatusReport() <-chan message.StatusReport {
	return g.statusReportChannel
}
