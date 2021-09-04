package main

import (
	"github.com/mayfield-z/ember/internal/pkg/gnb"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/ue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

func main() {
	gnb := gnb.NewGNB(
		"test-gnb",
		1,
		"208",
		"93",
		0x10,
		1,
		32,
		net.ParseIP("10.100.0.11"),
		38412,
		0x1,
		0x010203,
	)

	err := gnb.Run()
	logger.SetLogLevel(logrus.TraceLevel)
	if err != nil {
		logger.NgapLog.Errorf("error: %+v", err)
	}

	ue := ue.NewUE(
		"imsi-208930000000105",
		"208",
		"93",
		"8baf473f2f8fd09487cccbd7097c6862",
		"8e27b6af0e692e750f32667a3b14605d",
		"OPC",
		"8000",
		[]ue.PDU{
			{
				IpType: ue.IPv4,
				Apn:    "internet",
				Nssai: utils.SNSSAI{
					Sst: 0x01,
					Sd:  0x010203,
				},
			},
		},
	)

	ue.Run()
	time.Sleep(2 * time.Second)
	ue.RRCSetupRequest(gnb)
	select {}
}
