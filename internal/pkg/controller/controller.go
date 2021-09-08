package controller

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/mayfield-z/ember/internal/pkg/gnb"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/ue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	controller = Controller{}
)

type Controller struct {
	gnbList                    []*gnb.GNB
	ueList                     []*ue.UE
	templateGnb                *gnb.GNB
	templateUE                 *ue.UE
	gnbName                    string
	n2IpFrom                   net.IP
	n2IpTo                     net.IP
	n2IpNum                    uint32
	supiFrom                   string
	supiTo                     string
	ueNum                      uint32
	uePerGnb                   uint32
	uePerSec                   float64
	initPDUWhenAllUERegistered bool
	amfIp                      net.IP
	amfPort                    int
	configPath                 string

	initialed              bool
	supiPointer            string
	n2IpPointer            net.IP
	globalRANNodeIDPointer uint32
	nrCellIdentityPointer  uint64

	ctx        context.Context
	cancelFunc context.CancelFunc
	mutex      sync.Mutex
}

func ControllerSelf() *Controller {
	return &controller
}

func (c *Controller) addGnb(gnb *gnb.GNB) {
	c.mutex.Lock()
	c.gnbList = append(c.gnbList, gnb)
	c.mutex.Unlock()
}

func (c *Controller) addUE(ue *ue.UE) {
	c.mutex.Lock()
	c.ueList = append(c.ueList, ue)
	c.mutex.Unlock()
}

func (c *Controller) Init(configPath string) error {
	// do not change exec order
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())
	c.globalRANNodeIDPointer = 1
	c.nrCellIdentityPointer = 1
	c.configPath = configPath

	err := c.parseConfig(configPath)
	if err != nil {
		logger.ControllerLog.Errorf("Controller init failed: %+v", err)
	}

	c.initialed = true
	return nil
}

func (c *Controller) Reset() {
	c.Init(c.configPath)
}

func (c *Controller) parseConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "Controller ReadConfig failed")
	}

	// change to another place?
	level, err := logrus.ParseLevel(viper.GetString("app.logLevel"))
	if err != nil {
		level, _ = logrus.ParseLevel("info")
	}
	logger.SetLogLevel(level)
	logger.AppLog.Infof("Set log level to: %v", level)

	// TODO: check if exist and validation
	c.supiFrom = viper.GetString("ue.supiFrom")
	c.supiPointer = c.supiFrom
	c.supiTo = viper.GetString("ue.supiTo")

	c.gnbName = viper.GetString("gnb.name")
	n2IpFrom := net.ParseIP(viper.GetString("controller.n2IpFrom"))
	n2IpTo := net.ParseIP(viper.GetString("controller.n2IpTo"))
	n2IpNum := binary.BigEndian.Uint32(n2IpTo[12:]) - binary.BigEndian.Uint32(n2IpFrom[12:]) + 1
	c.n2IpFrom = n2IpFrom
	c.n2IpPointer = n2IpFrom
	c.n2IpTo = n2IpTo
	c.n2IpNum = n2IpNum
	logger.ControllerLog.Infof("N2 IP From: %v, To: %v, total: %v", n2IpFrom, n2IpTo, n2IpNum)

	c.ueNum = viper.GetUint32("controller.ueNum")
	c.uePerSec = viper.GetFloat64("controller.uePerSec")
	c.uePerGnb = viper.GetUint32("controller.uePerGnb")
	c.initPDUWhenAllUERegistered = viper.GetBool("controller.initPDUWhenAllUERegistered")
	logger.ControllerLog.Infof("Total UE num is: %v, register %v ue per second", c.ueNum, c.uePerSec)
	logger.ControllerLog.Infof("%v UE per gNB, will use %v gnb", c.uePerGnb, math.Ceil(float64(c.ueNum)/float64(c.uePerGnb)))
	if c.initPDUWhenAllUERegistered {
		logger.ControllerLog.Infof("Will inital PDU when all UE registered in core")
	} else {
		logger.ControllerLog.Infof("Will inital PDU when every UE registered in core")
	}

	c.amfIp = net.ParseIP(viper.GetString("amf.ip"))
	c.amfPort = viper.GetInt("amf.port")

	c.templateGnb = gnb.NewGNB(
		"templateGnb",
		c.globalRANNodeIDPointer,
		viper.GetString("gnb.mcc"),
		viper.GetString("gnb.mnc"),
		c.nrCellIdentityPointer,
		viper.GetUint32("gnb.tac"),
		uint8(viper.GetUint32("gnb.idLength")),
		c.amfIp,
		c.amfPort,
		uint8(viper.GetUint32("gnb.sst")),
		viper.GetUint32("gnb.sd"),
		c.ctx,
	)

	var pduSessions []ue.PDU
	for i := 0; ; i++ {
		pduConfig := viper.Sub(fmt.Sprintf("ue.sessions.%v", i))
		if pduConfig == nil {
			break
		}
		ipType, err := ue.ParseIpVersion(pduConfig.GetString("type"))
		if err != nil {
			return errors.WithMessagef(err, "parseConfig failed")
		}
		pdu := ue.PDU{
			IpType: ipType,
			Apn:    pduConfig.GetString("apn"),
			Nssai: utils.SNSSAI{
				Sst: uint8(pduConfig.GetUint32("sst")),
				Sd:  pduConfig.GetUint32("sd"),
			},
		}
		pduSessions = append(pduSessions, pdu)
	}

	c.templateUE = ue.NewUE(
		"templateUE",
		viper.GetString("ue.mcc"),
		viper.GetString("ue.mnc"),
		viper.GetString("ue.key"),
		viper.GetString("ue.op"),
		viper.GetString("ue.opType"),
		viper.GetString("ue.amf"),
		viper.GetString("ue.dataRate.uplink"),
		viper.GetString("ue.dataRate.downlink"),
		pduSessions,
		c.ctx,
	)

	return nil
}

func (c *Controller) Start() {
	if !c.initialed {
		logger.ControllerLog.Errorf("Start controller before initial it")
		return
	}
	go c.start()
}

func (c *Controller) start() {
	// TODO: UE init PDU
	ueIntervalInMicrosecond := int64(math.Ceil(float64(1000000) / c.uePerSec))
	ticker := time.NewTicker(time.Duration(ueIntervalInMicrosecond) * time.Microsecond)
	defer ticker.Stop()
	ueNumOfCurrentGnb := uint32(0)
	gnbPointer := 0

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if len(c.ueList) == int(c.ueNum) {
				//	idle
				return
			}
			currentUE := c.createAndAddUE()
			currentUE.Run()
			if ueNumOfCurrentGnb == c.uePerGnb {
				gnbPointer += 1
				ueNumOfCurrentGnb = 0
			} else {
				ueNumOfCurrentGnb += 1
			}
			if len(c.gnbList) == gnbPointer {
				if len(c.gnbList) == int(c.n2IpNum) {
					// idle
					logger.ControllerLog.Errorf("can not create more gNB")
					return
				}
				err := c.createAndAddGnb().Run()
				if err != nil {
					logger.ControllerLog.Errorf("create gnb failed: %+v", err)
					return
				}
			}
			currentGNB := c.gnbList[len(c.gnbList)-1]
			currentUE.RRCSetupRequest(currentGNB)
		}
	}
}

func (c *Controller) createAndAddGnb() *gnb.GNB {
	gnb := c.templateGnb.Copy(fmt.Sprintf("%v-%v", c.gnbName, len(c.gnbList))).
		SetGlobalRANNodeID(c.globalRANNodeIDPointer).
		SetNRCellIdentity(c.nrCellIdentityPointer)
	c.globalRANNodeIDPointer += 1
	c.nrCellIdentityPointer += 1

	c.addGnb(gnb)
	return gnb
}

func (c *Controller) createAndAddUE() *ue.UE {
	ue := c.templateUE.Copy(c.supiPointer)
	imsi, err := strconv.ParseUint(strings.Split(c.supiPointer, "-")[1], 10, 64)
	if err != nil {
		logger.ControllerLog.Errorf("create ue failed: %+v", err)
		return nil
	}
	c.supiPointer = fmt.Sprintf("imsi-%v", imsi+1)

	c.addUE(ue)
	return ue
}

func (c *Controller) getN2Ip() net.IP {
	ip := c.n2IpPointer
	res, flag := add1(ip[15])
	if flag == 0 {
		ip[15] = res
		c.n2IpPointer = ip
		return ip
	}
	res, flag = add1(ip[14])
	if flag == 0 {
		ip[14] = res
		c.n2IpPointer = ip
		return ip
	}
	res, flag = add1(ip[13])
	if flag == 0 {
		ip[13] = res
		c.n2IpPointer = ip
		return ip
	}
	res, flag = add1(ip[12])
	if flag == 0 {
		ip[12] = res
		c.n2IpPointer = ip
		return ip
	}
	panic("AMF ip overflow")
}

func add1(i uint8) (res uint8, flag uint8) {
	if i == ^uint8(0) {
		res = 0
		flag = 1
	} else {
		res = i + 1
		flag = 0
	}
	return
}
