package gnb

import (
	"encoding/hex"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"io"
	"syscall"
)

func (g *GNB) messageHandler() {
	c := g.getMessageChan()
	for {
		select {
		case <-g.ctx.Done():
			return
		default:
			select {
			case msg := <-c:
				g.logger.Tracef("recieve message %T, %v", msg, msg)
				switch msg.(type) {
				case message.RRCSetupRequest:
					g.handleRRCSetupRequestMessage(msg.(message.RRCSetupRequest))
				case message.RRCSetupComplete:
					g.handleRRCSetupCompleteMessage(msg.(message.RRCSetupComplete))
				case message.NASUplinkPdu:
					g.handleNASUplinkPdu(msg.(message.NASUplinkPdu))
				default:
					g.logger.Errorf("Type %T message not supported by gnb", msg)
				}
			}
		}
	}
}

func (g *GNB) handleRRCSetupRequestMessage(msg message.RRCSetupRequest) {
	g.logger.Debugf("handle RRCSetupRequestMessage")
	rRCSetupRequestMessage := msg
	supi := rRCSetupRequestMessage.SendBy
	if u := g.FindUEBySUPI(supi); u != nil {
		g.sendRRCRejectMessage(supi)
		g.logger.Errorf("UE %v has already in gnb", supi)
	} else {
		ue := &utils.GnbUe{
			SUPI: supi,
		}
		g.ueMapBySupi.Store(supi, ue)
		g.sendRRCSetupMessage(supi)
	}
	return
}

func (g *GNB) handleRRCSetupCompleteMessage(msg message.RRCSetupComplete) {
	g.logger.Debugf("handle RRCSetupCompleteMessage")
	rRCSetupCompleteMessage := msg
	ue := g.FindUEBySUPI(rRCSetupCompleteMessage.SendBy)
	if ue == nil {
		g.logger.Errorf("UE %v not in gnb but got RRCSetupCompelte", rRCSetupCompleteMessage.SendBy)
		return
	}

	ue.PLMN = rRCSetupCompleteMessage.PLMN
	id := g.allocateRANUENGAPID()
	if id == int64(-1) {
		g.logger.Errorf("RANUENGAPID full")
		return
	}
	ue.RANUENGAPID = id
	g.ueMapByRANUENGAPID.Store(id, ue)
	err := g.sendInitialUEMessage(id, rRCSetupCompleteMessage.NASRegistrationRequest)
	if err != nil {
		g.logger.Errorf("Initial UE Message send failed %+v", err)
		return
	}
}

func (g *GNB) handleNASUplinkPdu(msg message.NASUplinkPdu) {
	g.logger.Debugf("handle NASUplinkPDU")
	ue := g.FindUEBySUPI(msg.SendBy)
	if ue == nil {
		g.logger.Errorf("UE %v not in gnb but got NASUplinkPdu", msg.SendBy)
		return
	}
	g.logger.Tracef("received uplink pdu is from: %v\n %+v", msg.SendBy, hex.Dump(msg.PDU))
	pdu, err := g.buildUplinkNASTransport(ue, msg.PDU)
	if err != nil {
		g.logger.Errorf("handle NAS Uplink PDU failed %v", err)
		return
	}
	_, err = g.sctpConn.Write(pdu)
	if err != nil {
		g.logger.Errorf("NGSetupRequestPDU send failed, %v", err)
	}
}

func (g *GNB) sendRRCSetupMessage(supi string) {
	g.logger.Debugf("send RRCSetupMessage")
	mqueue.SendMessage(message.RRCSetup{}, supi)
}

func (g *GNB) sendRRCRejectMessage(supi string) {
	mqueue.SendMessage(message.RRCReject{}, supi)
}

func (g *GNB) sendNASDownlinkPdu(supi string, pdu []byte) {
	g.logger.Debugf("send NASDownlinkPdu")
	msg := message.NASDownlinkPdu{
		PDU: pdu,
	}
	mqueue.SendMessage(msg, supi)
}

// handle packages

func (g *GNB) sctpHandler(bufsize uint32) {
	conn := g.sctpConn
	for {
		select {
		case <-g.ctx.Done():
			return
		default:
			buf := make([]byte, bufsize)

			n, info, notification, err := conn.SCTPRead(buf)
			if err != nil {
				logger.NgapLog.Tracef("Read %d bytes", n)
				logger.NgapLog.Tracef("Packet content:\n%+v", hex.Dump(buf[:n]))
				switch err {
				case io.EOF, io.ErrUnexpectedEOF:
					logger.SctpLog.Debugln("Read EOF from client")
					return
				case syscall.EAGAIN:
					logger.SctpLog.Debugln("SCTP read timeout")
					continue
				case syscall.EINTR:
					logger.SctpLog.Debugf("SCTPRead: %+v", err)
					continue
				default:
					logger.SctpLog.Errorf("Handle connection[addr: %+v] error: %+v", conn.RemoteAddr(), err)
					return
				}
			}

			if notification != nil {
				// TODO: handle notification
				logger.NgapLog.Warnf("Received sctp notification[type 0x%x] but not handled", notification.Type())
			} else {
				if info == nil || info.PPID != ngap.PPID {
					logger.NgapLog.Warnln("Received SCTP PPID != 60, discard this packet")
					continue
				}

			}

			// TODO: concurrent on per-UE message
			g.ngapHandler(buf[:n])
			logger.NgapLog.Traceln("NGAP packet handled successful")
		}
	}
}

func (g *GNB) ngapHandler(buf []byte) {
	pdu, err := ngap.Decoder(buf)
	if err != nil {
		logger.NgapLog.Errorf("NGAP decode error: %+v", err)
		return
	}

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			logger.NgapLog.Errorln("InitiatingMessage is nil")
			return
		}
		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeDownlinkNASTransport:
			g.handleDownlinkNASTransport(pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			g.handleInitialContextSetupRequest(pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			g.handlePDUSessionResourceSetupRequest(pdu)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			logger.NgapLog.Errorln("Successful Outcome is nil")
			return
		}
		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			g.handleNGSetupResponse(pdu)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			logger.NgapLog.Errorln("Unsuccessful outcome is nil")
			return
		}
		switch unsuccessfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			g.handleNGSetupFailure(pdu)
		}
	}
}
