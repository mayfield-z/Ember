package main

import (
	"github.com/mayfield-z/ember/internal/pkg/context"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

func main() {
	gnb := context.NewGNB(
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

	ue := context.NewUE(
		"imsi-208930000000105",
		"208",
		"93",
		"8baf473f2f8fd09487cccbd7097c6862",
		"8e27b6af0e692e750f32667a3b14605d",
		"OPC",
		"8000",
		[]context.PDU{
			{
				IpType: context.IPv4,
				Apn:    "internet",
				Nssai: utils.SNSSAI{
					Sst: 0x01,
					Sd:  0x010203,
				},
			},
		},
	)

	time.Sleep(2 * time.Second)
	logger.AppLog.Infof("Add UE")
	err = gnb.AddUE(ue)
	if err != nil {
		logger.AppLog.Errorf("AddUE 也能错？: %+v", err)
	}

	_, err = gnb.InitialUE(ue)
	if err != nil {
		logger.AppLog.Errorf("InitialUE failed: %+v", err)
	}

	select {}
}
