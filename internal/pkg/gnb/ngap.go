package gnb

import (
	"encoding/binary"
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/free5gc/aper"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/pkg/errors"
	"net"
)

func Dial(ip net.IP, port int) (*sctp.SCTPConn, error) {
	laddr := &sctp.SCTPAddr{}
	raddr := &sctp.SCTPAddr{
		IPAddrs: []net.IPAddr{
			{
				IP:   ip,
				Zone: "",
			},
		},
		Port: port,
	}
	conn, err := sctp.DialSCTPExt("sctp", laddr, raddr, sctp.InitMsg{NumOstreams: sctp.SCTP_MAX_STREAM, MaxInitTimeout: 1000, MaxAttempts: 3}, nil, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "Dial failed")
	}

	logger.SctpLog.Infof("Dial LocalAddr: %v; RemoteAddr: %v", conn.LocalAddr(), conn.RemoteAddr())

	sndbuf, err := conn.GetWriteBuffer()
	if err != nil {
		logger.SctpLog.Fatalf("failed to get write buf: %+v", err)
	}
	rcvbuf, err := conn.GetReadBuffer()
	if err != nil {
		logger.SctpLog.Fatalf("failed to get read buf: %+v", err)
	}
	logger.SctpLog.Tracef("SndBufSize: %d, RcvBufSize: %d", sndbuf, rcvbuf)

	var info *sctp.SndRcvInfo
	if infoTmp, err := conn.GetDefaultSentParam(); err != nil {
		logger.SctpLog.Errorf("Get default sent param error: %+v", err)
		if err = conn.Close(); err != nil {
			logger.SctpLog.Errorf("Sctp closed error: %+v", err)
		}
	} else {
		info = infoTmp
		logger.NgapLog.Debugf("Get default sent param[value: %+v]", info)
	}

	info.PPID = ngap.PPID
	if err = conn.SetDefaultSentParam(info); err != nil {
		logger.SctpLog.Errorf("Set default sent param error: %+v", err)
		if err = conn.Close(); err != nil {
			logger.SctpLog.Errorf("Sctp closed error: %+v", err)
		}
		//
	} else {
		logger.NgapLog.Debugf("Set default sent param[value: %v]", info)
	}

	events := sctp.SCTP_EVENT_DATA_IO | sctp.SCTP_EVENT_SHUTDOWN | sctp.SCTP_EVENT_ASSOCIATION
	if err := conn.SubscribeEvents(events); err != nil {
		logger.NgapLog.Errorf("Failed to accept: %+v", err)
		if err = conn.Close(); err != nil {
			logger.NgapLog.Errorf("Close error: %+v", err)
		}
	} else {
		logger.NgapLog.Debugln("Subscribe SCTP event[DATA_IO, SHUTDOWN_EVENT, ASSOCIATION_CHANGE]")
	}

	if err := conn.SetReadBuffer(readBufSize); err != nil {
		logger.NgapLog.Errorf("Set read buffer error: %+v, accept failed", err)
		if err = conn.Close(); err != nil {
			logger.NgapLog.Errorf("Close error: %+v", err)
		}
	} else {
		logger.NgapLog.Debugf("Set read buffer to %d bytes", readBufSize)
	}

	return conn, nil

}

func (g *GNB) sendInitialUEMessage(id int64, nas []byte) error {
	initialUEMessage, err := g.buildInitialUEMessage(id, nas)
	if err != nil {
		return errors.WithMessagef(err, "InitialUEMessage PDU build failed.")
	}
	_, err = g.sctpConn.Write(initialUEMessage)
	if err != nil {
		return errors.WithMessagef(err, "InitialUEMessage Send failed.")
	}

	return nil
}

func (g *GNB) handleNGSetupResponse(pdu *ngapType.NGAPPDU) {
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
			g.gnbAmf.AmfName = aMFName.Value
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
			g.gnbAmf.GUAMI.Plmn = utils.DecodePLMNFromNgap(servedGUAMIList.List[0].GUAMI.PLMNIdentity)
			gUAMI := &servedGUAMIList.List[0].GUAMI
			amfId := uint32(gUAMI.AMFRegionID.Value.Bytes[0])
			amfId <<= 16
			amfId += uint32(binary.BigEndian.Uint16(gUAMI.AMFSetID.Value.Bytes))
			// TODO: check if cafe01
			amfId += uint32(gUAMI.AMFPointer.Value.Bytes[0])
			g.gnbAmf.GUAMI.AmfId = amfId
		case ngapType.ProtocolIEIDRelativeAMFCapacity:
			relativeAMFCapacity = ie.Value.RelativeAMFCapacity
			logger.NgapLog.Traceln("Handle IE RelativeAMFCapacity")
			if relativeAMFCapacity == nil {
				logger.NgapLog.Errorln("RelativeAMFCapacity is nil")
			}
			g.gnbAmf.Capacity = relativeAMFCapacity.Value
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

	g.gnbAmf.Connected = true
}

func (g *GNB) handleNGSetupFailure(pdu *ngapType.NGAPPDU) {
	// TODO: implement
	g.gnbAmf.Connected = false
}

func (g *GNB) handleDownlinkNASTransport(pdu *ngapType.NGAPPDU) {
	var (
		aMFUENGAPID      *ngapType.AMFUENGAPID
		rANUENGAPID      *ngapType.RANUENGAPID
		nASPDU           *ngapType.NASPDU
		aMFUENGAPIDValue int64
		rANUENGAPIDValue int64
		gnbue            *utils.GnbUe
	)
	initiatingMessage := pdu.InitiatingMessage
	downlinkNasTransport := initiatingMessage.Value.DownlinkNASTransport
	if downlinkNasTransport == nil {
		logger.NgapLog.Errorln("downlinkNasTransport is nil")
		return
	}
	logger.NgapLog.Infof("Handle downlinkNas Transport")
	for _, ie := range downlinkNasTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDSourceAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Traceln("Handle IE AMFUENGAPID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Errorln("AMFUENGAPID is nil")
				return
			}
			aMFUENGAPIDValue = aMFUENGAPID.Value
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Traceln("Handle RANUENGAPID")
			if rANUENGAPID == nil {
				logger.NgapLog.Errorln("RANUENGAPID is nil")
				return
			}
			rANUENGAPIDValue = rANUENGAPID.Value
		case ngapType.ProtocolIEIDNASPDU:
			nASPDU = ie.Value.NASPDU
			logger.NgapLog.Traceln("Handle NASPDU")
			if nASPDU == nil {
				logger.NgapLog.Errorln("NASPDU is nil")
				return
			}
			g.sendNASDownlinkPdu(gnbue.SUPI, nASPDU.Value)
		}
		gnbue = g.FindUEByRANUENGAPID(rANUENGAPIDValue)
		gnbue.AMFUENGAPID = aMFUENGAPIDValue
	}
}

func (g *GNB) handleInitialContextSetupRequest(pdu *ngapType.NGAPPDU) {
	var (
		//aMFUENGAPIDValue int64
		rANUENGAPIDValue int64
		nasPdu           []byte
	)
	initialContextSetupRequest := pdu.InitiatingMessage.Value.InitialContextSetupRequest
	if initialContextSetupRequest == nil {
		logger.NgapLog.Errorln("initialContextSetupRequest is nil")
		return
	}
	logger.NgapLog.Info("Handle initial context setup request")
	for _, ie := range initialContextSetupRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		// TODO: more case
		//case ngapType.ProtocolIEIDAMFUENGAPID:
		//	aMFUENGAPIDValue = ie.Value.AMFUENGAPID.Value
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPIDValue = ie.Value.RANUENGAPID.Value
		case ngapType.ProtocolIEIDNASPDU:
			nasPdu = ie.Value.NASPDU.Value
		}

	}
	ue := g.FindUEByRANUENGAPID(rANUENGAPIDValue)
	if ue == nil {
		logger.NgapLog.Errorf("cannot find ue with RANUENGAPIDValue %v", rANUENGAPIDValue)
		return
	}
	initialContextSetupResponse, err := g.buildInitialContextSetupResponse(ue)
	if err != nil {
		logger.NgapLog.Errorln("initial context setup response build failed")
	}
	_, err = g.sctpConn.Write(initialContextSetupResponse)
	if err != nil {
		logger.SctpLog.Errorln("initial context setup response send failed")
	}
	mqueue.SendMessage(message.NASDownlinkPdu{PDU: nasPdu}, ue.SUPI)
}

// buildNGSetupRequest referring to TS 38.413 -> 9.2.6.1
func (g *GNB) buildNGSetupRequest() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeNGSetup
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject
	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentNGSetupRequest
	initiatingMessage.Value.NGSetupRequest = new(ngapType.NGSetupRequest)

	nGSetupRequest := initiatingMessage.Value.NGSetupRequest
	nGSetupRequestIEs := &nGSetupRequest.ProtocolIEs

	// GlobalRANNodeID TS 38.413 -> 9.3.1.5
	ie := ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDGlobalRANNodeID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentGlobalRANNodeID
	ie.Value.GlobalRANNodeID = new(ngapType.GlobalRANNodeID)

	globalRANNodeID := ie.Value.GlobalRANNodeID
	globalRANNodeID.Present = ngapType.GlobalRANNodeIDPresentGlobalGNBID
	globalRANNodeID.GlobalGNBID = new(ngapType.GlobalGNBID)

	globalGNBID := globalRANNodeID.GlobalGNBID
	gNBID := &globalGNBID.GNBID
	gNBID.Present = ngapType.GNBIDPresentGNBID
	gNBID.GNBID = &aper.BitString{
		BitLength: uint64(g.idLength),
		Bytes:     utils.EncodeUint32(g.globalRANNodeID, 32),
	}
	globalGNBID.PLMNIdentity = utils.EncodePLMNToNgap(g.plmn)

	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	// RANNodeName
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANNodeName
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentRANNodeName
	ie.Value.RANNodeName = new(ngapType.RANNodeName)
	ie.Value.RANNodeName.Value = g.name

	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	// SupportedTAList
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSupportedTAList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentSupportedTAList
	ie.Value.SupportedTAList = new(ngapType.SupportedTAList)

	supportedTAList := ie.Value.SupportedTAList
	supportedTAItem := ngapType.SupportedTAItem{}
	supportedTAItem.TAC.Value = utils.EncodeUint32(g.tac, 24)

	broadcastPLMNList := &supportedTAItem.BroadcastPLMNList
	broadcastPLMNItem := ngapType.BroadcastPLMNItem{}
	broadcastPLMNItem.PLMNIdentity = utils.EncodePLMNToNgap(g.plmn)

	tAISliceSupportList := &broadcastPLMNItem.TAISliceSupportList
	sliceSupportItem := ngapType.SliceSupportItem{}
	sliceSupportItem.SNSSAI.SST.Value = utils.EncodeUint8(g.snssai.Sst)
	sliceSupportItem.SNSSAI.SD = new(ngapType.SD)
	sliceSupportItem.SNSSAI.SD.Value = utils.EncodeUint32(g.snssai.Sd, 24)

	tAISliceSupportList.List = append(tAISliceSupportList.List, sliceSupportItem)
	broadcastPLMNList.List = append(broadcastPLMNList.List, broadcastPLMNItem)
	supportedTAList.List = append(supportedTAList.List, supportedTAItem)

	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	// DefaultPagingDRX
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDDefaultPagingDRX
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentDefaultPagingDRX
	ie.Value.DefaultPagingDRX = new(ngapType.PagingDRX)
	ie.Value.DefaultPagingDRX.Value = ngapType.PagingDRXPresentV128

	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	return ngap.Encoder(pdu)
}

func (g *GNB) buildInitialUEMessage(id int64, nas []byte) ([]byte, error) {
	ue := g.FindUEByRANUENGAPID(id)
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initialingMessage := pdu.InitiatingMessage
	initialingMessage.ProcedureCode.Value = ngapType.ProcedureCodeInitialUEMessage
	initialingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore
	initialingMessage.Value.Present = ngapType.InitiatingMessagePresentInitialUEMessage
	initialingMessage.Value.InitialUEMessage = new(ngapType.InitialUEMessage)

	initialUEMessage := initialingMessage.Value.InitialUEMessage
	initialUEMessageIEs := &initialUEMessage.ProtocolIEs

	// RANUENGAPID
	ie := ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = id

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// NASPDU
	ie = ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNASPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentNASPDU
	ie.Value.NASPDU = new(ngapType.NASPDU)

	ie.Value.NASPDU.Value = nas

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// UserLocationInformation
	ie = ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)
	ie.Value.UserLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	ie.Value.UserLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := ie.Value.UserLocationInformation.UserLocationInformationNR
	userLocationInformationNR.TAI.TAC.Value = utils.EncodeUint32(g.tac, 24)
	userLocationInformationNR.TAI.PLMNIdentity = utils.EncodePLMNToNgap(ue.PLMN)
	// TODO:userLocationInformationNR.TimeStamp buhuisuan
	userLocationInformationNR.NRCGI.PLMNIdentity = utils.EncodePLMNToNgap(ue.PLMN)
	userLocationInformationNR.NRCGI.NRCellIdentity.Value.Bytes = utils.EncodeUint64(g.nci, 36)
	userLocationInformationNR.NRCGI.NRCellIdentity.Value.BitLength = 36

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// RRCEstablishmentCause
	ie = ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRRCEstablishmentCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentRRCEstablishmentCause
	ie.Value.RRCEstablishmentCause = new(ngapType.RRCEstablishmentCause)

	rRCEstablishmentCause := ie.Value.RRCEstablishmentCause
	rRCEstablishmentCause.Value = ngapType.RRCEstablishmentCausePresentMoSignalling

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	// UEContextRequest
	ie = ngapType.InitialUEMessageIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUEContextRequest
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialUEMessageIEsPresentUEContextRequest
	ie.Value.UEContextRequest = new(ngapType.UEContextRequest)

	uEContextRequest := ie.Value.UEContextRequest
	uEContextRequest.Value = ngapType.UEContextRequestPresentRequested

	initialUEMessageIEs.List = append(initialUEMessageIEs.List, ie)

	return ngap.Encoder(pdu)
}

func (g *GNB) buildUplinkNASTransport(ue *utils.GnbUe, nas []byte) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUplinkNASTransport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore
	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUplinkNASTransport
	initiatingMessage.Value.UplinkNASTransport = new(ngapType.UplinkNASTransport)

	uplinkNASTransport := initiatingMessage.Value.UplinkNASTransport
	uplinkNASTransportIEs := &uplinkNASTransport.ProtocolIEs

	// AMFUENGAPID
	ie := ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AMFUENGAPID

	uplinkNASTransportIEs.List = append(uplinkNASTransportIEs.List, ie)

	// RANUENGAPID
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RANUENGAPID

	uplinkNASTransportIEs.List = append(uplinkNASTransportIEs.List, ie)

	// NASPDU
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNASPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentNASPDU

	ie.Value.NASPDU = new(ngapType.NASPDU)
	ie.Value.NASPDU.Value = nas

	uplinkNASTransportIEs.List = append(uplinkNASTransportIEs.List, ie)

	// UserLocationInformation
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)
	ie.Value.UserLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationNR
	ie.Value.UserLocationInformation.UserLocationInformationNR = new(ngapType.UserLocationInformationNR)

	userLocationInformationNR := ie.Value.UserLocationInformation.UserLocationInformationNR
	userLocationInformationNR.TAI.TAC.Value = utils.EncodeUint32(g.tac, 24)
	userLocationInformationNR.TAI.PLMNIdentity = utils.EncodePLMNToNgap(ue.PLMN)
	// TODO:userLocationInformationNR.TimeStamp buhuisuan
	userLocationInformationNR.NRCGI.PLMNIdentity = utils.EncodePLMNToNgap(ue.PLMN)
	userLocationInformationNR.NRCGI.NRCellIdentity.Value.Bytes = utils.EncodeUint64(g.nci, 36)
	userLocationInformationNR.NRCGI.NRCellIdentity.Value.BitLength = 36

	uplinkNASTransportIEs.List = append(uplinkNASTransportIEs.List, ie)

	return ngap.Encoder(pdu)
}

func (g *GNB) buildInitialContextSetupResponse(ue *utils.GnbUe) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject
	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentInitialContextSetupResponse
	successfulOutcome.Value.InitialContextSetupResponse = new(ngapType.InitialContextSetupResponse)

	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	initialContextSetupResponseIEs := &initialContextSetupResponse.ProtocolIEs

	// AMFUENGAPID
	ie := ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AMFUENGAPID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// RANUENGAPID
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RANUENGAPID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	return ngap.Encoder(pdu)
}
