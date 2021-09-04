package gnb

import (
	"context"
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/pkg/errors"
	"net"
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
	amfAddress      net.IP
	amfPort         int
	snssai          utils.SNSSAI

	// auto
	ueMapBySupi          map[string]*GNBUE
	supiMapByRANUENGAPID map[int64]string
	rANUENGAPIDPointer   int64
	running              bool
	gnbAmf               utils.GnbAmf
	sctpConn             *sctp.SCTPConn
	ctx                  context.Context
	cancelFunc           context.CancelFunc
}

func NewGNB(name string, globalRANNodeID uint32, mcc, mnc string, nci uint64, tac uint32, idLength uint8, amfAddress net.IP, amfPort int, sst uint8, sd uint32) *GNB {
	//TODO: check if same name gnb exists
	mqueue.NewQueue(name)
	//TODO: change context
	ctx, cancelFunc := context.WithCancel(context.Background())
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
		amfAddress: amfAddress,
		amfPort:    amfPort,
		snssai: utils.SNSSAI{
			Sst: sst,
			Sd:  sd,
		},
		ueMapBySupi:          make(map[string]*GNBUE),
		supiMapByRANUENGAPID: make(map[int64]string),
		running:              false,
		sctpConn:             nil,
		ctx:                  ctx,
		cancelFunc:           cancelFunc,
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

func (g *GNB) getMessageChan() chan interface{} {
	return mqueue.GetMessageChan(g.name)
}

func (g *GNB) FindUEBySUPI(supi string) *GNBUE {
	if ue, ok := g.ueMapBySupi[supi]; ok {
		return ue
	}
	return nil
}

func (g GNB) FindUEByRANUENGAPID(id int64) *GNBUE {
	if supi, ok := g.supiMapByRANUENGAPID[id]; ok {
		if ue, ok := g.ueMapBySupi[supi]; ok {
			return ue
		}
	}
	return nil
}

func (g *GNB) Run() error {
	g.running = true
	conn, err := Dial(g.amfAddress, g.amfPort)
	if err != nil {
		return errors.Wrapf(err, "Failed to dial sctp address %v:%v", g.amfAddress, g.amfPort)
	}
	g.sctpConn = conn
	err = g.connectToAmf()
	if err != nil {
		return errors.Wrapf(err, "Send NGSetupRequestPDU error.")
	}
	go g.connectionHandler(readBufSize)
	go g.messageHandler()
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
	nGSetupRequest, err := g.BuildNGSetupRequest()
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

func (g *GNB) sendInitialUEMessage(id int64, nas []byte) error {
	initialUEMessage, err := g.BuildInitialUEMessage(id, nas)
	if err != nil {
		return errors.WithMessagef(err, "InitialUEMessage PDU build failed.")
	}
	_, err = g.sctpConn.Write(initialUEMessage)
	if err != nil {
		return errors.WithMessagef(err, "InitialUEMessage Send failed.")
	}

	return nil
}

type GNBUE struct {
	supi    string
	plmn    utils.PLMN
	lastNas []byte
}
