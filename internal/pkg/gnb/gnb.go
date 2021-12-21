package gnb

import (
	"context"
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
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
	rANUENGAPIDPointer      int64
	rANUENGAPIDPointerMutex sync.Mutex
	running                 bool
	gnbAmf                  utils.GnbAmf
	sctpConn                *sctp.SCTPConn
	Notify                  chan interface{}
	logger                  *logrus.Entry
	ctx                     context.Context
	cancelFunc              context.CancelFunc
}

func NewGNB(name string, globalRANNodeID uint32, mcc, mnc string, nci uint64, tac uint32, idLength uint8, n2Address, n3Address, amfAddress net.IP, amfPort int, sst uint8, sd uint32, parent context.Context) *GNB {
	//TODO: check if same name gnb exists
	mqueue.NewQueue(name)
	ctx, cancelFunc := context.WithCancel(parent)
	return &GNB{
		name:            name,
		globalRANNodeID: globalRANNodeID,
		plmn: utils.PLMN{
			Mcc: mcc,
			Mnc: mnc,
		},
		nci:        nci,
		idLength:   idLength,
		tac:        tac,
		n2Address:  n2Address,
		n3Address:  n3Address,
		amfAddress: amfAddress,
		amfPort:    amfPort,
		snssai: utils.SNSSAI{
			Sst: sst,
			Sd:  sd,
		},
		running:    false,
		sctpConn:   nil,
		Notify:     make(chan interface{}, 1),
		logger:     logger.GnbLog.WithFields(logrus.Fields{"name": name}),
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
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
	g.amfAddress = ip
	return g
}

func (g *GNB) getMessageChan() chan interface{} {
	return mqueue.GetMessageChan(g.name)
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

func (g *GNB) Run() error {
	g.running = true
	g.logger.Debugf("SCTP dial %v:%v", g.amfAddress, g.amfPort)
	conn, err := Dial(g.n2Address, g.amfAddress, g.amfPort)
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
