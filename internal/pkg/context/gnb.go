package context

import (
	"context"
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/mayfield-z/ember/internal/pkg/logger"
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
	ueMapBySupi          map[string]*UE
	supiMapByRANUENGAPID map[int64]string
	rANUENGAPIDPointer   int64
	running              bool
	gnbAmf               utils.GnbAmf
	sctpConn             *sctp.SCTPConn
	ctx                  context.Context
	cancelFunc           context.CancelFunc
}

func NewGNB(name string, globalRANNodeID uint32, mcc, mnc string, nci uint64, tac uint32, idLength uint8, amfAddress net.IP, amfPort int, sst uint8, sd uint32) *GNB {
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
		ueMapBySupi:          make(map[string]*UE),
		supiMapByRANUENGAPID: make(map[int64]string),
		running:              false,
		sctpConn:             nil,
	}
}

func (g *GNB) Name() string {
	return g.name
}

func (g *GNB) Running() bool {
	return g.running
}

func (g *GNB) Connected() bool {
	return g.gnbAmf.Connected
}

func (g *GNB) AddUE(ue *UE) error {
	if _, ok := g.ueMapBySupi[ue.SUPI()]; ok {
		return errors.New("ue has already in gnb")
	} else {
		g.ueMapBySupi[ue.SUPI()] = ue
		ue.rRCSetup(g)
	}
	return nil
}

func (g *GNB) FindUEBySUPI(supi string) *UE {
	if ue, ok := g.ueMapBySupi[supi]; ok {
		return ue
	}
	return nil
}

func (g GNB) FindUEByRANUENGAPID(id int64) *UE {
	if supi, ok := g.supiMapByRANUENGAPID[id]; ok {
		if ue, ok := g.ueMapBySupi[supi]; ok {
			return ue
		}
	}
	return nil
}

func (g *GNB) Run() error {
	g.running = true
	g.ctx, g.cancelFunc = context.WithCancel(context.Background())
	conn, err := Dial(g.amfAddress, g.amfPort)
	if err != nil {
		return errors.Wrapf(err, "Failed to dial sctp address %v:%v", g.amfAddress, g.amfPort)
	}
	g.sctpConn = conn
	err = g.connectToAmf()
	if err != nil {
		return errors.Wrapf(err, "Send NGSetupRequestPDU error.")
	}
	go g.connectionHandler(readBufSize, g.ctx)
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

func (g *GNB) InitialUE(ue *UE) (int64, error) {
	logger.AppLog.Traceln("Start initial ue")
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

	g.supiMapByRANUENGAPID[rANUENGAPID] = ue.supi
	err := g.sendInitialUEMessage(rANUENGAPID)
	if err != nil {
		return -1, err
	}
	return rANUENGAPID, nil
}

func (g *GNB) sendInitialUEMessage(id int64) error {
	initialUEMessage, err := g.BuildInitialUEMessage(id)
	if err != nil {
		return errors.WithMessagef(err, "InitialUEMessage PDU build failed.")
	}
	_, err = g.sctpConn.Write(initialUEMessage)
	if err != nil {
		return errors.WithMessagef(err, "InitialUEMessage Send failed.")
	}

	//TODO:CM-STATE CHANGE
	return nil
}
