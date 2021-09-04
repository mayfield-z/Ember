package packet_driver

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

type Sender interface {
	SendTo()
}

type GoSender struct {
}

func Send(src, dst net.IP) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}

	// ethernet?
	ip := &layers.IPv4{
		Protocol: layers.IPProtocolSCTP,
		SrcIP:    src,
		DstIP:    dst,
	}
	sctp := &layers.SCTP{
		DstPort:         38412,
		VerificationTag: 0,
		Checksum:        0,
	}
}
