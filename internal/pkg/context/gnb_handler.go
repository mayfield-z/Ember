package context

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"io"
	"syscall"
)

func (g *GNB) connectionHandler(bufsize uint32, ctx context.Context) {
	conn := g.sctpConn
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buf := make([]byte, bufsize)

			n, info, notification, err := conn.SCTPRead(buf)
			logger.NgapLog.Tracef("Read %d bytes", n)
			logger.NgapLog.Tracef("Packet content:\n%+v", hex.Dump(buf[:n]))
			if err != nil {
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
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			logger.NgapLog.Errorln("Successful Outcome is nil")
			return
		}
		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			handleNGSetupResponse(g, pdu)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			logger.NgapLog.Errorln("Unsuccessful outcome is nil")
			return
		}
		switch unsuccessfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			handleNGSetupFailure(g, pdu)
		}
	}
}

func handleNGSetupResponse(gnb *GNB, pdu *ngapType.NGAPPDU) {
	var (
		aMFName             *ngapType.AMFName
		servedGUAMIList     *ngapType.ServedGUAMIList
		relativeAMFCapacity *ngapType.RelativeAMFCapacity
		pLMNSupportList     *ngapType.PLMNSupportList
	)

	successfulOutcome := pdu.SuccessfulOutcome
	nGSetupResponse := successfulOutcome.Value.NGSetupResponse
	if nGSetupResponse == nil {
		logger.NgapLog.Errorln("NGSetupResponse is nil")
		return
	}
	logger.NgapLog.Infof("Handle NG Setup response")
	for _, ie := range nGSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFName:
			aMFName = ie.Value.AMFName
			logger.NgapLog.Traceln("Handle IE AMFName")
			if aMFName == nil {
				logger.NgapLog.Errorln("AMFName is nil")
				return
			}
			gnb.gnbAmf.AmfName = aMFName.Value
		case ngapType.ProtocolIEIDServedGUAMIList:
			servedGUAMIList = ie.Value.ServedGUAMIList
			logger.NgapLog.Traceln("Handle IE GUAMIList")
			if servedGUAMIList == nil {
				logger.NgapLog.Errorln("ServedGUAMIList is nil")
			}
			if len(servedGUAMIList.List) == 0 {
				logger.NgapLog.Errorln("ServedGUAMIList len is 0")
			}
			// TODO: support multiple GUAMIList
			gnb.gnbAmf.GUAMI.Plmn = utils.DecodePLMNFromNgap(servedGUAMIList.List[0].GUAMI.PLMNIdentity)
			gUAMI := &servedGUAMIList.List[0].GUAMI
			amfId := uint32(gUAMI.AMFRegionID.Value.Bytes[0])
			amfId <<= 16
			amfId += uint32(binary.BigEndian.Uint16(gUAMI.AMFSetID.Value.Bytes))
			// TODO: check if cafe01
			amfId += uint32(gUAMI.AMFPointer.Value.Bytes[0])
			gnb.gnbAmf.GUAMI.AmfId = amfId
		case ngapType.ProtocolIEIDRelativeAMFCapacity:
			relativeAMFCapacity = ie.Value.RelativeAMFCapacity
			logger.NgapLog.Traceln("Handle IE RelativeAMFCapacity")
			if relativeAMFCapacity == nil {
				logger.NgapLog.Errorln("RelativeAMFCapacity is nil")
			}
			gnb.gnbAmf.Capacity = relativeAMFCapacity.Value
		case ngapType.ProtocolIEIDPLMNSupportList:
			pLMNSupportList = ie.Value.PLMNSupportList
			// TODO: implement
			if pLMNSupportList == nil {
				logger.NgapLog.Errorln("pLMNSupportList is nil")
			}
			if len(pLMNSupportList.List) == 0 {
				logger.NgapLog.Errorln("pLMNSupportList len is 0")
			}
		}
	}

	gnb.gnbAmf.Connected = true
}

func handleNGSetupFailure(gnb *GNB, pdu *ngapType.NGAPPDU) {
	// TODO: implement
	gnb.gnbAmf.Connected = false
}
