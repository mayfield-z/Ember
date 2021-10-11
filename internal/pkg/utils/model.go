package utils

import "net"

type PLMN struct {
	Mcc string
	Mnc string
}

type GnbAmf struct {
	Connected bool
	AmfName   string
	GUAMI     GUAMI
	Capacity  int64
}

type GnbUe struct {
	SUPI                  string
	PLMN                  PLMN
	AMFUENGAPID           int64
	RANUENGAPID           int64
	PDUSessionID          int64
	QosFlowIdentifier     int64
	FiveQI                int64
	GTPTEID               []byte
	TransportLayerAddress net.IP
}

type SNSSAI struct {
	Sst uint8
	Sd  uint32
}

type GUAMI struct {
	Plmn  PLMN
	AmfId uint32
}

type IpVersion int

const (
	IPv4 IpVersion = iota
	IPv6
	IPv4_AND_IPv6
)

type PDU struct {
	Id     uint8
	IpType IpVersion
	Apn    string
	Nssai  SNSSAI
}

type UeGnb struct {
	Name string
}
