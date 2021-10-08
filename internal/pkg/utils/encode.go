package utils

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/free5gc/ngap/logger"
	"github.com/free5gc/ngap/ngapType"
	"math"
	"strings"
)

func EncodeUint8(u uint8) []byte {
	b := make([]byte, 1)
	b[0] = u
	return b
}

func EncodeUint16(u uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, u)
	return b
}

func EncodeUint32(u uint32, length int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, u)
	return b[4-int(math.Ceil(float64(length)/8)):]
}

func EncodeUint64(u uint64, length int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	return b[8-int(math.Ceil(float64(length)/8)):]
}

func EncodePLMNToNgap(plmn PLMN) ngapType.PLMNIdentity {
	var hexString string
	mcc := strings.Split(plmn.Mcc, "")
	mnc := strings.Split(plmn.Mnc, "")
	if len(plmn.Mnc) == 2 {
		hexString = mcc[1] + mcc[0] + "f" + mcc[2] + mnc[1] + mnc[0]
	} else {
		hexString = mcc[1] + mcc[0] + mnc[0] + mcc[2] + mnc[2] + mnc[1]
	}

	var ngapPlmnId ngapType.PLMNIdentity
	if plmnId, err := hex.DecodeString(hexString); err != nil {
		logger.NgapLog.Warnf("Decode plmn failed: %+v", err)
	} else {
		ngapPlmnId.Value = plmnId
	}
	return ngapPlmnId
}

func DecodePLMNFromNgap(plmn ngapType.PLMNIdentity) PLMN {
	value := plmn.Value
	var p PLMN
	hexString := strings.Split(hex.EncodeToString(value), "")
	p.Mcc = hexString[1] + hexString[0] + hexString[3]
	if hexString[2] == "f" {
		p.Mnc = hexString[5] + hexString[4]
	} else {
		p.Mnc = hexString[2] + hexString[5] + hexString[4]
	}
	return p
}

func EncodeMsin(msin string) []byte {
	var ret []byte
	encoded, _ := hex.DecodeString(reverse(msin))
	for _, b := range encoded {
		ret = append([]byte{b}, ret...)
	}
	return ret
}

func reverse(s string) string {
	ret := ""
	for _, i := range s {
		ret = string(i) + ret
	}
	return ret
}
