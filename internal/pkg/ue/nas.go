package ue

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
)

func (u *UE) buildRegistrationRequest() ([]byte, error) {
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

	// 5GS mobile identity
	mobileIdentity5GS := append([]uint8{0x01}, utils.EncodePLMNToNgap(u.plmn).Value...)
	mobileIdentity5GS = append(mobileIdentity5GS, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00)
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

	m.GmmMessage.RegistrationRequest = registrationRequest

	return m.PlainNasEncode()
}

func (u *UE) handleAuthenticationRequest(msg *nas.Message) {
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
		u.nasLogger.Infof("Send authentication response")
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
