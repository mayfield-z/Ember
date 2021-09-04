package context

import (
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/mayfield-z/ember/internal/pkg/utils"
)

func (u *UE) BuildRegistrationRequest() ([]byte, error) {
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
