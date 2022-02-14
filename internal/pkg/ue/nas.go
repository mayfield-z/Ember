package ue

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/pkg/errors"
	"net"
)

func (u *UE) buildRegistrationRequest(capability bool) ([]byte, error) {
	u.nasLogger.Debug("build RegistrationRequest")
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmMessage.SetMessageType(nas.MsgTypeRegistrationRequest)

	// wtf is iei??
	registrationRequest := nasMessage.NewRegistrationRequest(0)
	registrationRequest.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	registrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	registrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	registrationRequest.RegistrationRequestMessageIdentity.SetMessageType(nas.MsgTypeRegistrationRequest)

	// NAS Key set identifier and 5GS registration type
	registrationRequest.NgksiAndRegistrationType5GS.Octet = 0x79
	registrationRequest.NgksiAndRegistrationType5GS.SetRegistrationType5GS(nasMessage.RegistrationType5GSInitialRegistration)

	// 5GS mobile identity
	mobileIdentity5GS := append([]uint8{0x01}, utils.EncodePLMNToNgap(u.plmn).Value...)
	msin := utils.EncodeMsin(u.supi[len(u.supi)-10 : len(u.supi)])
	msin = append([]uint8{0x00, 0x00, 0x00, 0x00}, msin...)
	mobileIdentity5GS = append(mobileIdentity5GS, msin...)
	registrationRequest.MobileIdentity5GS.SetLen(uint16(len(mobileIdentity5GS)))
	registrationRequest.MobileIdentity5GS.SetMobileIdentity5GSContents(mobileIdentity5GS)
	registrationRequest.UESecurityCapability = new(nasType.UESecurityCapability)

	// UE Security Capability
	uESecurityCapability := registrationRequest.UESecurityCapability
	uESecurityCapability.SetIei(nasMessage.RegistrationRequestUESecurityCapabilityType)
	uESecurityCapability.SetLen(4)
	uESecurityCapability.SetEA0_5G(1)
	uESecurityCapability.SetEA1_128_5G(1)
	uESecurityCapability.SetEA2_128_5G(1)
	uESecurityCapability.SetEA3_128_5G(1)
	uESecurityCapability.SetEA4_5G(0)
	uESecurityCapability.SetEA5_5G(0)
	uESecurityCapability.SetEA6_5G(0)
	uESecurityCapability.SetEA7_5G(0)
	uESecurityCapability.SetIA0_5G(1)
	uESecurityCapability.SetIA1_128_5G(1)
	uESecurityCapability.SetIA2_128_5G(1)
	uESecurityCapability.SetIA3_128_5G(1)
	uESecurityCapability.SetIA4_5G(0)
	uESecurityCapability.SetIA5_5G(0)
	uESecurityCapability.SetIA6_5G(0)
	uESecurityCapability.SetIA7_5G(0)
	uESecurityCapability.SetEEA0(1)
	uESecurityCapability.SetEEA1_128(1)
	uESecurityCapability.SetEEA2_128(1)
	uESecurityCapability.SetEEA3_128(1)
	uESecurityCapability.SetEEA4(0)
	uESecurityCapability.SetEEA5(0)
	uESecurityCapability.SetEEA6(0)
	uESecurityCapability.SetEEA7(0)
	uESecurityCapability.SetEIA0(1)
	uESecurityCapability.SetEIA1_128(1)
	uESecurityCapability.SetEIA2_128(1)
	uESecurityCapability.SetEIA3_128(1)
	uESecurityCapability.SetEIA4(0)
	uESecurityCapability.SetEIA5(0)
	uESecurityCapability.SetEIA6(0)
	uESecurityCapability.SetEIA7(0)

	if capability {
		registrationRequest.Capability5GMM = &nasType.Capability5GMM{
			Iei:   nasMessage.RegistrationRequestCapability5GMMType,
			Len:   1,
			Octet: [13]uint8{0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		}
	}

	m.GmmMessage.RegistrationRequest = registrationRequest

	return m.PlainNasEncode()
}

func (u *UE) handleAuthenticationRequest(msg *nas.Message) {
	u.nasLogger.Debug("handle AuthenticationRequest")
	var authenticationResponse message.NASUplinkPdu
	var err error

	authenticationResponse.SendBy = u.supi
	rand := msg.AuthenticationRequest.GetRANDValue()
	autn := msg.AuthenticationRequest.GetAUTN()

	paramAutn, check := u.DeriveRESStarAndSetKey(rand[:], autn[:], u.snn)

	switch check {
	case mACFailure:
		u.nasLogger.Warnf("Authentication request message failed with MAC failure")
		// TODO: change state of UE, send response
	case sQNFailure:
		u.nasLogger.Warnf("Authentication request message failed with SQN failure")
		// TODO: change state of UE, send response
	case successful:
		u.nasLogger.Debugf("Send authentication response")
		authenticationResponse.PDU, err = u.buildAuthenticationResponse(paramAutn, "")
		u.nasLogger.Tracef("authentication response is:\n %+v", hex.Dump(authenticationResponse.PDU))
		if err != nil {
			u.nasLogger.Errorf("Build Authentication Response failed")
			return
		}
		// TODO: change state of UE
		mqueue.SendMessage(authenticationResponse, u.gnb.Name)
	}

}

func (u *UE) buildAuthenticationResponse(authenticationResponseParam []uint8, eapMsg string) ([]byte, error) {
	u.nasLogger.Debug("build AuthenticationResponse")
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmMessage.SetMessageType(nas.MsgTypeAuthenticationResponse)

	authenticationResponse := nasMessage.NewAuthenticationResponse(0)
	authenticationResponse.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)

	authenticationResponse.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationResponse.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationResponse.AuthenticationResponseMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationResponse)

	if len(authenticationResponseParam) > 0 {
		authenticationResponse.AuthenticationResponseParameter = nasType.NewAuthenticationResponseParameter(nasMessage.AuthenticationResponseAuthenticationResponseParameterType)
		authenticationResponse.AuthenticationResponseParameter.SetLen(uint8(len(authenticationResponseParam)))
		copy(authenticationResponse.AuthenticationResponseParameter.Octet[:], authenticationResponseParam[0:16])
	} else if eapMsg != "" {
		rawEapMsg, _ := base64.StdEncoding.DecodeString(eapMsg)
		authenticationResponse.EAPMessage = nasType.NewEAPMessage(nasMessage.AuthenticationResponseEAPMessageType)
		authenticationResponse.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		authenticationResponse.EAPMessage.SetEAPMessage(rawEapMsg)
	}

	m.GmmMessage.AuthenticationResponse = authenticationResponse

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	nasPdu := data.Bytes()
	return nasPdu, nil
}

func (u *UE) handleSecurityModeCommand(msg *nas.Message) {
	u.nasLogger.Debug("handle SecurityModeCommand")
	switch msg.SecurityModeCommand.SelectedNASSecurityAlgorithms.GetTypeOfCipheringAlgorithm() {
	case 0:
		u.nasLogger.Debug("Type of ciphering algorithm is 5G-EA0")
	case 1:
		u.nasLogger.Debug("Type of ciphering algorithm is 128-5G-EA1")
	case 2:
		u.nasLogger.Debug("Type of ciphering algorithm is 128-5G-EA2")
	}

	switch msg.SecurityModeCommand.SelectedNASSecurityAlgorithms.GetTypeOfIntegrityProtectionAlgorithm() {
	case 0:
		u.nasLogger.Debug("Type of integrity protection algorithm is 5G-IA0")
	case 1:
		u.nasLogger.Debug("Type of integrity protection algorithm is 128-5G-IA1")
	case 2:
		u.nasLogger.Debug("Type of integrity protection algorithm is 128-5G-IA2")
	}

	// checking BIT RINMR that triggered registration request in security mode complete.
	rinmr := msg.SecurityModeCommand.Additional5GSecurityInformation.GetRINMR()

	// getting NAS Security Mode Complete.
	securityModeComplete, err := u.buildSecurityModeComplete(rinmr)
	if err != nil {
		u.nasLogger.Errorf("Error sending Security Mode Complete: %v", err)
		return
	}
	pdu, err := u.encodeNASPduWithSecurity(securityModeComplete, true, nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext)
	if err != nil {
		u.nasLogger.Errorf("Error sending Security Mode Complete: %v", err)
		return
	}

	// sending to GNB
	mqueue.SendMessage(message.NASUplinkPdu{PDU: pdu, SendBy: u.supi}, u.gnb.Name)
}

func (u *UE) handleDLNASTransport(msg *nas.Message) {
	u.nasLogger.Debug("handle DLNASTransport")
	payload := msg.GmmMessage.DLNASTransport.PayloadContainer.Buffer
	msg2 := new(nas.Message)
	err := msg2.PlainNasDecode(&payload)
	if err != nil {
		u.nasLogger.Errorf("DLNASTransport payload decode failed: %v", err)
	}
	if msg2.GsmHeader.GetMessageType() == nas.MsgTypePDUSessionEstablishmentAccept {
		u.nasLogger.Debug("PDU Session Establishment Accept")
		ueIp := msg2.PDUSessionEstablishmentAccept.GetPDUAddressInformation()
		u.ip = net.IPv4(ueIp[0], ueIp[1], ueIp[2], ueIp[3])
		u.smFSM.Event(eventSMPDUSessionEstablishmentAccept)
		u.SendStatusReport(message.UEPDUSessionEstablishmentAccept)
	} else if msg2.GsmHeader.GetMessageType() == nas.MsgTypePDUSessionEstablishmentReject {
		u.nasLogger.Debug("PDU Session Establishment Reject")
		u.SendStatusReport(message.UEPDUSessionEstablishmentReject)
	}
}

func (u *UE) buildSecurityModeComplete(rinmr uint8) ([]byte, error) {
	u.nasLogger.Debug("build SecurityModeComplete")
	registrationRequest, err := u.buildRegistrationRequest(true)
	if err != nil {
		return nil, errors.WithMessage(err, "build registration request in security mode complete failed.")
	}

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeSecurityModeComplete)

	m.GmmMessage.SecurityModeComplete = nasMessage.NewSecurityModeComplete(0)
	securityModeComplete := m.GmmMessage.SecurityModeComplete
	securityModeComplete.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	// TODO: modify security header type if need security protected
	securityModeComplete.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	securityModeComplete.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	securityModeComplete.SecurityModeCompleteMessageIdentity.SetMessageType(nas.MsgTypeSecurityModeComplete)

	securityModeComplete.IMEISV = nasType.NewIMEISV(nasMessage.SecurityModeCompleteIMEISVType)
	securityModeComplete.IMEISV.SetLen(9)
	securityModeComplete.SetOddEvenIdic(0)
	securityModeComplete.SetTypeOfIdentity(nasMessage.MobileIdentity5GSTypeImeisv)
	securityModeComplete.SetIdentityDigit1(1)
	securityModeComplete.SetIdentityDigitP_1(1)
	securityModeComplete.SetIdentityDigitP(1)

	if registrationRequest != nil {
		securityModeComplete.NASMessageContainer = nasType.NewNASMessageContainer(nasMessage.SecurityModeCompleteNASMessageContainerType)
		securityModeComplete.NASMessageContainer.SetLen(uint16(len(registrationRequest)))
		securityModeComplete.NASMessageContainer.SetNASMessageContainerContents(registrationRequest)
	}

	m.GmmMessage.SecurityModeComplete = securityModeComplete

	data := new(bytes.Buffer)
	err = m.GmmMessageEncode(data)
	if err != nil {
		return nil, errors.WithMessage(err, "gmm message encode in build security mode complete failed")
	}

	nasPdu := data.Bytes()
	return nasPdu, nil
}

func (u *UE) handleRegistrationAccept(msg *nas.Message) {
	u.nasLogger.Debug("handle RegistrationAccept")
	registrationAccept := msg.RegistrationAccept
	if registrationAccept == nil {
		u.nasLogger.Errorf("registration accept is nil")
		return
	}
	u.gGuti = registrationAccept.GetTMSI5G()
	u.aMFPointer = registrationAccept.GetAMFPointer()
	u.aMFRegionID = registrationAccept.GetAMFRegionID()
	u.aMFSetID = registrationAccept.GetAMFSetID()

	registrationComplete, err := u.buildRegistrationComplete()
	if err != nil {
		u.nasLogger.Errorf("handle registration accept failed")
		return
	}

	u.rmFSM.Event(eventRMRegistrationAccept)
	u.SendStatusReport(message.UERegistrationSuccess)
	mqueue.SendMessage(message.NASUplinkPdu{PDU: registrationComplete, SendBy: u.supi}, u.gnb.Name)
}

func (u *UE) buildRegistrationComplete() ([]byte, error) {
	u.nasLogger.Debug("build RegistrationComplete")
	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationComplete)

	registrationComplete := nasMessage.NewRegistrationComplete(0)
	registrationComplete.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	registrationComplete.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	registrationComplete.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	registrationComplete.RegistrationCompleteMessageIdentity.SetMessageType(nas.MsgTypeRegistrationComplete)

	m.GmmMessage.RegistrationComplete = registrationComplete

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		return nil, errors.WithMessage(err, "build registration complete failed")
	}

	pdu := data.Bytes()

	pdu, err = u.encodeNASPduWithSecurity(pdu, false, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered)
	if err != nil {
		return nil, errors.WithMessage(err, "build registration complete failed")
	}

	return pdu, nil
}

func (u *UE) buildPDUSessionEstablishmentRequest(id uint8) ([]byte, error) {
	u.nasLogger.Debug("build PDUSessionEstablishmentRequest")
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionEstablishmentRequest)

	pduSessionEstablishmentRequest := nasMessage.NewPDUSessionEstablishmentRequest(0)
	pduSessionEstablishmentRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pduSessionEstablishmentRequest.SetMessageType(nas.MsgTypePDUSessionEstablishmentRequest)
	pduSessionEstablishmentRequest.PDUSessionID.SetPDUSessionID(id)
	pduSessionEstablishmentRequest.PTI.SetPTI(0x00)
	pduSessionEstablishmentRequest.IntegrityProtectionMaximumDataRate.SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink(0xff)
	pduSessionEstablishmentRequest.IntegrityProtectionMaximumDataRate.SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink(0xff)

	pduSessionEstablishmentRequest.PDUSessionType = nasType.NewPDUSessionType(nasMessage.PDUSessionEstablishmentRequestPDUSessionTypeType)
	pduSessionEstablishmentRequest.PDUSessionType.SetPDUSessionTypeValue(uint8(0x01)) //IPv4 type

	pduSessionEstablishmentRequest.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(nasMessage.PDUSessionEstablishmentRequestExtendedProtocolConfigurationOptionsType)
	protocolConfigurationOptions := nasConvert.NewProtocolConfigurationOptions()
	protocolConfigurationOptions.AddIPAddressAllocationViaNASSignallingUL()
	protocolConfigurationOptions.AddDNSServerIPv4AddressRequest()
	protocolConfigurationOptions.AddDNSServerIPv6AddressRequest()
	pcoContents := protocolConfigurationOptions.Marshal()
	pcoContentsLength := len(pcoContents)
	pduSessionEstablishmentRequest.ExtendedProtocolConfigurationOptions.SetLen(uint16(pcoContentsLength))
	pduSessionEstablishmentRequest.ExtendedProtocolConfigurationOptions.SetExtendedProtocolConfigurationOptionsContents(pcoContents)

	m.GsmMessage.PDUSessionEstablishmentRequest = pduSessionEstablishmentRequest

	data := new(bytes.Buffer)
	EncodePDUSessionEstablishmentRequest(m.GsmMessage.PDUSessionEstablishmentRequest, data)

	pdu := data.Bytes()
	return pdu, nil
}

func (u *UE) buildULNasTransportPDUSessionEstablishmentRequest(pduSessionNumber uint8) ([]byte, error) {
	u.nasLogger.Debug("build ULNASTransportPDUSessionEstablishmentRequest")
	pduSession := u.pduSessions[pduSessionNumber]

	pdu, err := u.buildPDUSessionEstablishmentRequest(pduSession.Id)
	if err != nil {
		return nil, errors.WithMessage(err, "build UL NAS transport PDU Session Establishment Request failed.")
	}

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeULNASTransport)

	ulNasTransport := nasMessage.NewULNASTransport(0)
	ulNasTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	ulNasTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
	ulNasTransport.PduSessionID2Value.SetIei(nasMessage.ULNASTransportPduSessionID2ValueType)
	ulNasTransport.PduSessionID2Value.SetPduSessionID2Value(pduSession.Id)
	ulNasTransport.RequestType = new(nasType.RequestType)
	ulNasTransport.RequestType.SetIei(nasMessage.ULNASTransportRequestTypeType)
	ulNasTransport.RequestType.SetRequestTypeValue(nasMessage.ULNASTransportRequestTypeInitialRequest)

	dnn := []byte(pduSession.Apn)
	ulNasTransport.DNN = new(nasType.DNN)
	ulNasTransport.DNN.SetIei(nasMessage.ULNASTransportDNNType)
	ulNasTransport.DNN.SetLen(uint8(len(dnn)))
	ulNasTransport.DNN.SetDNN(dnn)

	ulNasTransport.SNSSAI = nasType.NewSNSSAI(nasMessage.ULNASTransportSNSSAIType)
	ulNasTransport.SNSSAI.SetLen(4)
	var sdTemp [3]uint8
	sd := utils.EncodeUint32(pduSession.Nssai.Sd, 24)
	copy(sdTemp[:], sd)
	ulNasTransport.SNSSAI.SetSST(pduSession.Nssai.Sst)
	ulNasTransport.SNSSAI.SetSD(sdTemp)

	ulNasTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(nasMessage.PayloadContainerTypeN1SMInfo)
	ulNasTransport.PayloadContainer.SetLen(uint16(len(pdu)))
	ulNasTransport.PayloadContainer.SetPayloadContainerContents(pdu)

	m.GmmMessage.ULNASTransport = ulNasTransport

	data := new(bytes.Buffer)
	err = m.GmmMessageEncode(data)
	if err != nil {
		return nil, err
	}

	pdu = data.Bytes()

	pdu, err = u.encodeNASPduWithSecurity(pdu, false, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered)
	if err != nil {
		return nil, err
	}

	return pdu, nil
}

func (u *UE) encodeNASPduWithSecurity(payload []byte, newSecurityContext bool, securityHeaderType uint8) ([]byte, error) {
	var sequenceNumber uint8
	msg := nas.NewMessage()
	err := msg.PlainNasDecode(&payload)
	if err != nil {
		return nil, errors.WithMessage(err, "encode NAS PDU with security failed")
	}
	msg.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    securityHeaderType,
	}
	if newSecurityContext {
		u.ULCount.Set(0, 0)
		u.DLCount.Set(0, 0)
	}

	sequenceNumber = u.ULCount.SQN()
	err = security.NASEncrypt(u.cipheringAlg, u.knasEnc, u.ULCount.Get(), security.Bearer3GPP, security.DirectionUplink, payload)
	if err != nil {
		return nil, errors.WithMessage(err, "NAS encrypt failed")
	}

	payload = append([]byte{sequenceNumber}, payload[:]...)
	mac32 := make([]byte, 4)

	mac32, err = security.NASMacCalculate(u.integrityAlg, u.knasInt, u.ULCount.Get(), security.Bearer3GPP, security.DirectionUplink, payload)
	if err != nil {
		return nil, errors.WithMessage(err, "NASMacCalculate failed")
	}

	payload = append(mac32, payload[:]...)
	msgSecurityHeader := []byte{msg.SecurityHeader.ProtocolDiscriminator, msg.SecurityHeader.SecurityHeaderType}
	payload = append(msgSecurityHeader, payload[:]...)

	u.ULCount.AddOne()

	return payload, nil
}

func EncodePDUSessionEstablishmentRequest(a *nasMessage.PDUSessionEstablishmentRequest, buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, &a.ExtendedProtocolDiscriminator.Octet)
	binary.Write(buffer, binary.BigEndian, &a.PDUSessionID.Octet)
	binary.Write(buffer, binary.BigEndian, &a.PTI.Octet)
	binary.Write(buffer, binary.BigEndian, &a.PDUSESSIONESTABLISHMENTREQUESTMessageIdentity.Octet)
	binary.Write(buffer, binary.BigEndian, &a.IntegrityProtectionMaximumDataRate.Octet)
	if a.PDUSessionType != nil {
		binary.Write(buffer, binary.BigEndian, &a.PDUSessionType.Octet)
	}
	if a.SSCMode != nil {
		binary.Write(buffer, binary.BigEndian, &a.SSCMode.Octet)
	}
	if a.Capability5GSM != nil {
		binary.Write(buffer, binary.BigEndian, a.Capability5GSM.GetIei())
		binary.Write(buffer, binary.BigEndian, a.Capability5GSM.GetLen())
		binary.Write(buffer, binary.BigEndian, a.Capability5GSM.Octet[:a.Capability5GSM.GetLen()])
	}
	if a.MaximumNumberOfSupportedPacketFilters != nil {
		binary.Write(buffer, binary.BigEndian, a.MaximumNumberOfSupportedPacketFilters.GetIei())
		binary.Write(buffer, binary.BigEndian, &a.MaximumNumberOfSupportedPacketFilters.Octet)
	}
	if a.AlwaysonPDUSessionRequested != nil {
		binary.Write(buffer, binary.BigEndian, &a.AlwaysonPDUSessionRequested.Octet)
	}
	if a.SMPDUDNRequestContainer != nil {
		binary.Write(buffer, binary.BigEndian, a.SMPDUDNRequestContainer.GetIei())
		binary.Write(buffer, binary.BigEndian, a.SMPDUDNRequestContainer.GetLen())
		binary.Write(buffer, binary.BigEndian, &a.SMPDUDNRequestContainer.Buffer)
	}
	if a.ExtendedProtocolConfigurationOptions != nil {
		binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.GetIei())
		binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.GetLen())
		binary.Write(buffer, binary.BigEndian, &a.ExtendedProtocolConfigurationOptions.Buffer)
	}
}
