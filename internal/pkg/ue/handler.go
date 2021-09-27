package ue

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/free5gc/UeauCommon"
	"github.com/free5gc/milenage"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/security"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"reflect"
	"regexp"
)

func (u *UE) messageHandler() {
	for {
		select {
		case <-u.ctx.Done():
			return
		default:
			c := u.getMessageChan()
			select {
			case msg := <-c:
				switch msg.(type) {
				case message.RRCSetup:
					u.handleRRCSetupMessage(msg.(message.RRCSetup))
				case message.RRCReject:
					u.handleRRCRejectMessage(msg.(message.RRCReject))
				case message.NASDownlinkPdu:
					u.nasHandler(msg.(message.NASDownlinkPdu))
				default:
					u.logger.Errorf("Type %T message not supported by ue", msg)
				}
			}
		}
	}
}

func (u *UE) handleRRCSetupMessage(msg message.RRCSetup) {
	u.logger.Debugf("Handle message RRCSetup")
	u.rrcFSM.Event(eventRRCSetup)
	u.sendRRCSetupCompleteMessage(u.gnb.Name)
	u.Notify <- message.UERRCSetupSuccess{}
}

func (u *UE) handleRRCRejectMessage(msg message.RRCReject) {
	u.logger.Debugf("Handle message RRCReject")
	// implement when needed.
	panic("not implement")
}
func (u *UE) sendRRCSetupCompleteMessage(name string) {
	u.logger.Debugf("Send message RRCSetupComplete")
	nas, err := u.buildRegistrationRequest()
	if err != nil {
		u.logger.Errorf("Build RegistrationRequest failed: %+v", err)
		return
	}
	msg := message.RRCSetupComplete{
		NASRegistrationRequest: nas,
		SendBy:                 u.supi,
		PLMN:                   u.plmn,
	}
	mqueue.SendMessage(msg, name)
}

func (u *UE) nasHandler(msg message.NASDownlinkPdu) {
	u.logger.Debugf("Handle message NASPdu")
	payload := msg.PDU
	if payload == nil {
		u.nasLogger.Errorf("NAS payload is empty")
	}

	decodedMsg := new(nas.Message)
	decodedMsg.SecurityHeaderType = nas.GetSecurityHeaderType(payload) & 0x0f
	u.nasLogger.Tracef("securityHeaderType is %v", decodedMsg.SecurityHeaderType)
	if decodedMsg.SecurityHeaderType == nas.SecurityHeaderTypePlainNas {
		err := decodedMsg.PlainNasDecode(&payload)
		if err != nil {
			u.nasLogger.Errorf("decode NAS pdu failed %v", err)
		}
	} else {
		//securityHeader := payload[0:6]
		//sequenceNumber := payload[6]
		//
		//receivedMac32 := securityHeader[2:]
		//
		//payload = payload[6:]
		//
		//ciphered := false
		//switch msg.SecurityHeaderType {
		//case nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext:
		//	ciphered = true
		//}
	}

	if decodedMsg.GmmMessage == nil {
		u.nasLogger.Errorf("gmm message is nil")
	}
	switch decodedMsg.GmmMessage.GetMessageType() {
	case nas.MsgTypeAuthenticationRequest:
		u.handleAuthenticationRequest(decodedMsg)
	default:
		u.nasLogger.Errorf("unsupported message type %T", decodedMsg.GmmMessage.GetMessageType())
	}
}

type AuthenticationStatus int

const (
	mACFailure AuthenticationStatus = iota
	sQNFailure
	successful
)

func (u *UE) DeriveRESStarAndSetKey(RAND, AUTN []byte, snName string) ([]byte, AuthenticationStatus) {
	// parameters for authentication challenge.
	mac_a, mac_s := make([]byte, 8), make([]byte, 8)
	CK, IK := make([]byte, 16), make([]byte, 16)
	RES := make([]byte, 8)
	AK, AKstar := make([]byte, 6), make([]byte, 6)

	// Get OPC, K, SQN, AMF from USIM.
	OPC, _ := hex.DecodeString(u.op)
	K, _ := hex.DecodeString(u.key)

	//TODO: what is squence number
	sqnUe, _ := hex.DecodeString(u.sqn)
	AMF, _ := hex.DecodeString(u.amf)

	milenage.F2345(OPC, K, RAND, RES, CK, IK, AK, AKstar)

	sqnHn, _, mac_aHn := u.deriveAUTN(AUTN, AK)

	milenage.F1(OPC, K, RAND, sqnHn, AMF, mac_a, mac_s)

	if !reflect.DeepEqual(mac_a, mac_aHn) {
		return nil, mACFailure
	}

	if bytes.Compare(sqnUe, sqnHn) > 0 {
		milenage.F2345(OPC, K, RAND, RES, CK, IK, AK, AKstar)

		amfSynch, _ := hex.DecodeString("0000")

		milenage.F1(OPC, K, RAND, sqnUe, amfSynch, mac_a, mac_s)

		sqnUeXorAK := make([]byte, 6)
		for i := 0; i < len(sqnUe); i++ {
			sqnUeXorAK[i] = sqnUe[i] ^ AKstar[i]
		}

		failureParam := append(sqnUeXorAK, mac_s...)

		return failureParam, sQNFailure
	}

	u.sqn = fmt.Sprintf("%x", sqnHn)

	key := append(CK, IK...)
	FC := UeauCommon.FC_FOR_RES_STAR_XRES_STAR_DERIVATION
	P0 := []byte(snName)
	P1 := RAND
	P2 := RES

	u.DerivateKamf(key, snName, sqnHn, AK)
	u.DerivateAlgKey()
	kdfVal_for_resStar := UeauCommon.GetKDFValue(key, FC, P0, UeauCommon.KDFLen(P0), P1, UeauCommon.KDFLen(P1), P2, UeauCommon.KDFLen(P2))
	return kdfVal_for_resStar[len(kdfVal_for_resStar)/2:], successful
}

func (u *UE) deriveAUTN(autn []byte, ak []uint8) ([]byte, []byte, []byte) {

	sqn := make([]byte, 6)

	// get SQNxorAK
	SQNxorAK := autn[0:6]
	amf := autn[6:8]
	mac_a := autn[8:]

	// get SQN
	for i := 0; i < len(SQNxorAK); i++ {
		sqn[i] = SQNxorAK[i] ^ ak[i]
	}

	// return SQN, amf, mac_a
	return sqn, amf, mac_a
}

func (u *UE) DerivateKamf(key []byte, snName string, SQN, AK []byte) {

	FC := UeauCommon.FC_FOR_KAUSF_DERIVATION
	P0 := []byte(snName)
	SQNxorAK := make([]byte, 6)
	for i := 0; i < len(SQN); i++ {
		SQNxorAK[i] = SQN[i] ^ AK[i]
	}
	P1 := SQNxorAK
	Kausf := UeauCommon.GetKDFValue(key, FC, P0, UeauCommon.KDFLen(P0), P1, UeauCommon.KDFLen(P1))
	P0 = []byte(snName)
	Kseaf := UeauCommon.GetKDFValue(Kausf, UeauCommon.FC_FOR_KSEAF_DERIVATION, P0, UeauCommon.KDFLen(P0))

	supiRegexp, _ := regexp.Compile("(?:imsi|supi)-([0-9]{5,15})")
	groups := supiRegexp.FindStringSubmatch(u.supi)

	P0 = []byte(groups[1])
	L0 := UeauCommon.KDFLen(P0)
	P1 = []byte{0x00, 0x00}
	L1 := UeauCommon.KDFLen(P1)

	u.kamf = UeauCommon.GetKDFValue(Kseaf, UeauCommon.FC_FOR_KAMF_DERIVATION, P0, L0, P1, L1)
}

// Algorithm key Derivation function defined in TS 33.501 Annex A.9
func (u *UE) DerivateAlgKey() {
	// Security Key
	P0 := []byte{security.NNASEncAlg}
	L0 := UeauCommon.KDFLen(P0)
	P1 := []byte{u.cipheringAlg}
	L1 := UeauCommon.KDFLen(P1)

	kenc := UeauCommon.GetKDFValue(u.kamf, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	copy(u.knasEnc[:], kenc[16:32])

	// Integrity Key
	P0 = []byte{security.NNASIntAlg}
	L0 = UeauCommon.KDFLen(P0)
	P1 = []byte{u.integrityAlg}
	L1 = UeauCommon.KDFLen(P1)

	kint := UeauCommon.GetKDFValue(u.kamf, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	copy(u.knasInt[:], kint[16:32])
}
