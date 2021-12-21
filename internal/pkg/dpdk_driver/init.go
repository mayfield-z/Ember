package dpdk_driver

import (
	"fmt"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/mempool"
	"os"
)

const (
	numMbufs       = 8191
	mbufCacheSize  = 250
	defaultBufSize = 2048 + 128
)

func InitDpdk() {
	var devInfo ethdev.DevInfo
	n, err := eal.Init(os.Args)
	if err != nil {
		logger.DpdkLog.Panicf("Error with EAL initialization")
	}

	os.Args[n], os.Args = os.Args[0], os.Args[n:]

	nbPorts := ethdev.CountAvail()
	logger.DpdkLog.Infof("Avail eth dev number is: %v", nbPorts)

	mbufPool, err := mempool.CreateMbufPool("MBUF_POOL",
		defaultBufSize,
		mbufCacheSize,
		mempool.OptSocket(int(eal.SocketID())),
		mempool.OptCacheSize(512))
	if mbufPool == nil {
		logger.DpdkLog.Panicf("Cannot create mbuf pool")
	}

	port0 := ethdev.Port(0)
	port0.InfoGet(&devInfo)
	fmt.Printf("port 0: %v\n", devInfo.DriverName())
	if !port0.IsValid() {
		logger.DpdkLog.Errorf("Port 0 is not valid")
	}

	err = port0.DevConfigure(1, 1)
	if err != nil {
		logger.DpdkLog.Errorf("Port 0 configure failed")
	}

}
