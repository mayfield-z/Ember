package context

import (
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/free5gc/aper"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/mayfield-z/ember/internal/pkg/logger"
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

// BuildNGSetupRequest referring to TS 38.413 -> 9.2.6.1
func (g *GNB) BuildNGSetupRequest() ([]byte, error) {
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

func (g *GNB) BuildInitialUEMessage(id int64) ([]byte, error) {
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

	naspdu, err := ue.BuildRegistrationRequest()
	if err != nil {
		return nil, errors.WithMessage(err, "Nas RegistrationRequest build failed.")
	}
	ie.Value.NASPDU.Value = naspdu

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
	userLocationInformationNR.TAI.PLMNIdentity = utils.EncodePLMNToNgap(ue.plmn)
	// TODO:userLocationInformationNR.TimeStamp buhuisuan
	userLocationInformationNR.NRCGI.PLMNIdentity = utils.EncodePLMNToNgap(ue.plmn)
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
