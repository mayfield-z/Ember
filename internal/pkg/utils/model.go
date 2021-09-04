package utils

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
type SNSSAI struct {
	Sst uint8
	Sd  uint32
}

type GUAMI struct {
	Plmn  PLMN
	AmfId uint32
}
