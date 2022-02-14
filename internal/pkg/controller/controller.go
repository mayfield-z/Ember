package controller

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/mayfield-z/ember/internal/pkg/gnb"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/ue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	controller = Controller{name: "controller"}
)

type Controller struct {
	name                       string
	gnbList                    []*gnb.GNB
	ueList                     []*ue.UE
	templateGnb                *gnb.GNB
	templateUE                 *ue.UE
	n2Interface                *net.Interface
	n3Interface                *net.Interface
	n2OriginalAddresses        []net.Addr
	n3OriginalAddresses        []net.Addr
	n2AddedIp                  []net.IP
	n3AddedIp                  []net.IP
	gnbName                    string
	n2IpFrom                   net.IP
	n2IpTo                     net.IP
	n2IpNum                    uint32
	n3IpFrom                   net.IP
	n3IpTo                     net.IP
	n3IpNum                    uint32
	supiFrom                   string
	supiTo                     string
	supiNum                    uint32
	ueNum                      uint32
	uePerGnb                   uint32
	uePerSec                   float64
	realMaxUe                  uint32
	initPDUWhenAllUERegistered bool
	amfIp                      net.IP
	amfPort                    int

	initialed              bool
	running                bool
	supiPointer            string
	n2IpPointer            net.IP
	n3IpPointer            net.IP
	ueIdPointer            uint64
	pduSessionIdPointer    uint8
	globalRANNodeIDPointer uint32
	nrCellIdentityPointer  uint64

	ueNumOfCurrentGnb uint32 // atomic
	emulatedUeNum     uint32 // atomic
	gnbCounter        uint32 // atomic

	ctx                 context.Context
	cancelFunc          context.CancelFunc
	wg                  *sync.WaitGroup
	statusReportChannel chan message.StatusReport
	n2AddedIpMutex      sync.Mutex
	n3AddedIpMutex      sync.Mutex
	ueListMutex         sync.Mutex
	ueInfoMutex         sync.Mutex
	gnbListMutex        sync.Mutex
}

func Self() *Controller {
	return &controller
}

func (c *Controller) addGnb(gnb *gnb.GNB) {
	c.gnbListMutex.Lock()
	defer c.gnbListMutex.Unlock()
	logger.ControllerLog.Debugf("Adding gNB %v", gnb.NodeName())
	c.gnbList = append(c.gnbList, gnb)
}

func (c *Controller) getGnbByIndex(index int) *gnb.GNB {
	c.gnbListMutex.Lock()
	defer c.gnbListMutex.Unlock()
	if index < len(c.gnbList) {
		return c.gnbList[index]
	}
	return nil
}

func (c *Controller) addUE(ue *ue.UE) {
	c.ueListMutex.Lock()
	defer c.ueListMutex.Unlock()
	logger.ControllerLog.Debugf("Adding UE: %v", ue.NodeName())
	c.ueList = append(c.ueList, ue)
}

func (c *Controller) Init() error {
	// do not change exec order
	logger.ControllerLog.Debugf("Start inital controller")
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())
	c.globalRANNodeIDPointer = 1
	c.nrCellIdentityPointer = 1
	c.ueIdPointer = 0
	c.pduSessionIdPointer = 0
	c.wg = &sync.WaitGroup{}

	err := c.parseConfig()
	if err != nil {
		logger.ControllerLog.Errorf("Controller init failed: %+v", err)
	}

	c.initialed = true
	return nil
}

func (c *Controller) Reset() {
	c.Init()
}

func (c *Controller) parseConfig() error {
	var err error
	// TODO: check if exist and validation
	c.supiFrom = viper.GetString("ue.supiFrom")
	c.supiPointer = c.supiFrom
	c.supiTo = viper.GetString("ue.supiTo")
	c.supiNum = calcSUPINum(c.supiFrom, c.supiTo)
	logger.ControllerLog.Infof("GetSUPI from: %v, to: %v, total: %v", c.supiFrom, c.supiTo, c.supiNum)

	c.n2Interface, err = net.InterfaceByName(viper.GetString("controller.n2Interface"))
	if err != nil {
		return errors.Wrap(err, "Controller GetN2InterfaceByName failed")
	}

	c.n3Interface, err = net.InterfaceByName(viper.GetString("controller.n2Interface"))
	if err != nil {
		return errors.Wrap(err, "Controller GetN3InterfaceByName failed")
	}

	c.n2OriginalAddresses, err = c.n2Interface.Addrs()
	if err != nil {
		return errors.Wrap(err, "Controller GetN2InterfaceAddr failed")
	}

	c.n3OriginalAddresses, err = c.n3Interface.Addrs()
	if err != nil {
		return errors.Wrap(err, "Controller GetN3InterfaceAddr failed")
	}

	c.gnbName = viper.GetString("gnb.name")
	n2IpFrom := net.ParseIP(viper.GetString("controller.n2IpFrom"))
	n2IpTo := net.ParseIP(viper.GetString("controller.n2IpTo"))
	n2IpNum := binary.BigEndian.Uint32(n2IpTo[12:]) - binary.BigEndian.Uint32(n2IpFrom[12:]) + 1
	c.n2IpFrom = n2IpFrom
	c.n2IpPointer = n2IpFrom
	c.n2IpTo = n2IpTo
	c.n2IpNum = n2IpNum
	logger.ControllerLog.Infof("N2 IP From: %v, To: %v, total: %v", n2IpFrom, n2IpTo, n2IpNum)

	n3IpFrom := net.ParseIP(viper.GetString("controller.n3IpFrom"))
	n3IpTo := net.ParseIP(viper.GetString("controller.n3IpTo"))
	n3IpNum := binary.BigEndian.Uint32(n3IpTo[12:]) - binary.BigEndian.Uint32(n3IpFrom[12:]) + 1
	c.n3IpFrom = n3IpFrom
	c.n3IpPointer = n3IpFrom
	c.n3IpTo = n3IpTo
	c.n3IpNum = n3IpNum
	logger.ControllerLog.Infof("N2 IP From: %v, To: %v, total: %v", n3IpFrom, n3IpTo, n3IpNum)

	c.ueNum = viper.GetUint32("controller.ueNum")
	c.uePerSec = viper.GetFloat64("controller.uePerSec")
	c.uePerGnb = viper.GetUint32("controller.uePerGnb")
	c.realMaxUe = calcRealMaxUeNum(c.n2IpNum, c.n3IpNum, c.uePerGnb, c.supiNum, c.ueNum)
	c.initPDUWhenAllUERegistered = viper.GetBool("controller.initPDUWhenAllUERegistered")
	logger.ControllerLog.Infof("Real UE max is: %v, register %v ue per second", c.realMaxUe, c.uePerSec)
	logger.ControllerLog.Infof("%v UE per gNB, will use %v gnb", c.uePerGnb, math.Ceil(float64(c.realMaxUe)/float64(c.uePerGnb)))
	if c.initPDUWhenAllUERegistered {
		logger.ControllerLog.Infof("Will inital PDU after all UE registered in core")
	} else {
		logger.ControllerLog.Infof("Will inital PDU once UE registered in core")
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
		c.n2IpPointer,
		c.n3IpPointer,
		c.amfIp,
		c.amfPort,
		uint8(viper.GetUint32("gnb.sst")),
		viper.GetUint32("gnb.sd"),
		c.ctx,
	)

	var pduSessions []utils.PDU
	for i := 0; ; i++ {
		pduConfig := viper.Sub(fmt.Sprintf("ue.sessions.%v", i))
		if pduConfig == nil {
			break
		}
		ipType, err := ue.ParseIpVersion(pduConfig.GetString("type"))
		if err != nil {
			return errors.WithMessagef(err, "parseConfig failed")
		}
		pdu := utils.PDU{
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
		0,
		c.ctx,
	)

	mqueue.DelQueue(c.templateGnb.NodeName())
	mqueue.DelQueue(c.templateUE.NodeName())
	return nil
}

func (c *Controller) Start() {
	if !c.initialed {
		logger.ControllerLog.Errorf("Start controller before initial it")
		return
	}
	c.running = true
	go c.start()
}

func (c *Controller) start() {
	// TODO: UE init PDU
	logger.ControllerLog.Infof("Controller start")
	ueIntervalInMicrosecond := int64(math.Ceil(float64(1000000) / c.uePerSec))
	ticker := time.NewTicker(time.Duration(ueIntervalInMicrosecond) * time.Microsecond)
	defer ticker.Stop()
	c.SendStatusReport(message.ControllerStart)

	for c.emulatedUeNum < c.realMaxUe {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.emulatedUeNumAdd1()
			c.wg.Add(1)
			go c.emulateOneUEUserPlane(!c.initPDUWhenAllUERegistered)
		}
	}

	if c.initPDUWhenAllUERegistered {
		for _, u := range c.ueList {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				go c.establishOneUEPDUSession(u)
			}
		}
	}

	c.wg.Wait()
	c.Stop()
}

func (c *Controller) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Controller) Stop() {
	if c.running {
		c.cancelFunc()
		c.cleanUp()
		c.running = false
	}
}

func (c *Controller) cleanUp() {
	c.n2AddedIpMutex.Lock()
	defer c.n2AddedIpMutex.Unlock()
	for _, ip := range c.n2AddedIp {
		logger.ControllerLog.Debugf("deleting N2 IP addresses: %v", ip)
		err := utils.DelIpFromInterface(ip, c.n2Interface)
		if err != nil {
			logger.ControllerLog.Errorf("can not del N2 IP addresses: %v, %v", ip, err)
		}
	}

	c.n3AddedIpMutex.Lock()
	defer c.n3AddedIpMutex.Unlock()
	for _, ip := range c.n3AddedIp {
		logger.ControllerLog.Debugf("deleting N3 IP addresses: %v", ip)
		err := utils.DelIpFromInterface(ip, c.n3Interface)
		if err != nil {
			logger.ControllerLog.Errorf("can not del N3 IP addresses: %v, %v", ip, err)
		}
	}
}

func (c *Controller) createAndAddGnb() *gnb.GNB {
	logger.ControllerLog.Infof("Creating gnb %v-%v", c.gnbName, len(c.gnbList))
	g := c.templateGnb.Copy(fmt.Sprintf("%v-%v", c.gnbName, len(c.gnbList))).
		SetGlobalRANNodeID(c.globalRANNodeIDPointer).
		SetNRCellIdentity(c.nrCellIdentityPointer).
		SetN2Addresses(c.n2IpPointer).
		SetN3Addresses(c.n3IpPointer)
	c.globalRANNodeIDPointer += 1
	c.nrCellIdentityPointer += 1

	err := utils.AddIpToInterface(c.n2IpPointer, c.n2Interface)
	if err != nil {
		logger.ControllerLog.Fatalf("Add N2 Address to interface %s failed: %v", c.n2Interface.Name, err)
	} else {
		c.n2AddedIpMutex.Lock()
		c.n2AddedIp = append(c.n2AddedIp, c.n2IpPointer)
		c.n2AddedIpMutex.Unlock()
	}

	err = utils.AddIpToInterface(c.n3IpPointer, c.n3Interface)
	if err != nil {
		logger.ControllerLog.Fatalf("Add N3 Address to interface %s failed: %v", c.n3Interface.Name, err)
	} else {
		c.n3AddedIpMutex.Lock()
		c.n3AddedIp = append(c.n3AddedIp, c.n3IpPointer)
		c.n3AddedIpMutex.Unlock()
	}

	c.n2IpPointer = utils.Add1(c.n2IpPointer)
	c.n3IpPointer = utils.Add1(c.n3IpPointer)
	c.addGnb(g)
	return g
}

func (c *Controller) createAndAddUE() *ue.UE {
	c.ueInfoMutex.Lock()
	logger.ControllerLog.Infof("Creating UE: %v", c.supiPointer)
	//u := c.templateUE.Copy(c.supiPointer, c.ueIdPointer, c.pduSessionIdPointer)
	u := c.templateUE.Copy(c.supiPointer, c.ueIdPointer, 0)
	imsi, err := strconv.ParseUint(strings.Split(c.supiPointer, "-")[1], 10, 64)
	if err != nil {
		logger.ControllerLog.Errorf("create ue failed: %+v", err)
		return nil
	}
	c.supiPointer = fmt.Sprintf("imsi-%v", imsi+1)
	c.ueIdPointer += 1
	c.pduSessionIdPointer += u.GetPDUSessionNum()
	c.ueInfoMutex.Unlock()

	c.addUE(u)
	c.ueNumOfCurrentGnbAdd1()
	return u
}

func (c *Controller) emulateOneUEUserPlane(setupPDUSession bool) {
	// create UE
	c.SendStatusReport(message.EmulateUE)
	currentUE := c.createAndAddUE()
	currentUE.Run()
	defer currentUE.Stop(c.wg)

	if c.getEmulatedUeNum()%c.uePerGnb == 1 || c.uePerGnb == 1 {
		c.ueNumOfCurrentGnbClear()
		c.gnbCounterAdd1()
		// create gnb
		if len(c.gnbList) == int(c.n2IpNum) || len(c.gnbList) == int(c.n3IpNum) {
			logger.ControllerLog.Errorf("can not create more gNB")
			return
		}
		c.SendStatusReport(message.EmulateGNB)
		err := c.createAndAddGnb().Run()
		if err != nil {
			logger.ControllerLog.Errorf("create gnb failed: %+v", err)
			return
		}
	}

	logger.ControllerLog.Debugf("UE %v is using No.%v GNB", currentUE.GetSUPI(), int(currentUE.GetID()/uint64(c.uePerGnb)))
	currentGNB := c.getGnbByIndex(int(currentUE.GetID() / uint64(c.uePerGnb)))
	if currentGNB == nil {
		logger.ControllerLog.Logger.Panicf("GNB empty pointer")
	}
	for !currentGNB.Connected() {
	}

	// start RRCSetup
	currentUE.RRCSetupRequest(currentGNB)

	select {
	case msg := <-currentUE.StatusReport():
		switch msg.Event {
		case message.UERegistrationSuccess:
			logger.ControllerLog.Infof("UE %v Registration Success", currentUE.GetSUPI())
		case message.UERegistrationReject:
			// TODO: handle reject
			logger.ControllerLog.Infof("UE %v Registration Reject", currentUE.GetSUPI())
		}
	}

	if setupPDUSession {
		c.establishOneUEPDUSession(currentUE)
		logger.ControllerLog.Debugf("UE IP is: %v, TEID is :%v", currentUE.GetIP(), currentGNB.FindUEBySUPI(currentUE.GetSUPI()).GTPTEID)
	}
}

func (c *Controller) establishOneUEPDUSession(u *ue.UE) {
	if !u.Running() {
		u.Run()
		defer u.Stop(c.wg)
	}

	u.EstablishPDUSession(0)
	select {
	case msg := <-u.StatusReport():
		switch msg.Event {
		case message.UEPDUSessionEstablishmentAccept:
			logger.ControllerLog.Infof("UE %v PDU Session Establishment Accept", u.GetSUPI())
		case message.UEPDUSessionEstablishmentReject:
			logger.ControllerLog.Infof("UE %v PDU Session Establishment Reject", u.GetSUPI())
		}
	}
}

func (c *Controller) emulatedUeNumAdd1() {
	atomic.AddUint32(&c.emulatedUeNum, 1)
}

func (c *Controller) getEmulatedUeNum() uint32 {
	return atomic.LoadUint32(&c.emulatedUeNum)
}

func (c *Controller) ueNumOfCurrentGnbAdd1() {
	atomic.AddUint32(&c.ueNumOfCurrentGnb, 1)
}

func (c *Controller) ueNumOfCurrentGnbClear() {
	atomic.StoreUint32(&c.ueNumOfCurrentGnb, 0)
}

func (c *Controller) getUeNumOfCurrentGnb() uint32 {
	return atomic.LoadUint32(&c.ueNumOfCurrentGnb)
}

func (c *Controller) gnbCounterAdd1() {
	atomic.AddUint32(&c.gnbCounter, 1)
}

func (c *Controller) getGnbCounter() uint32 {
	return atomic.LoadUint32(&c.gnbCounter)
}

func (c *Controller) SendStatusReport(event message.Event) {
	statusReport := message.StatusReport{
		NodeName: c.name,
		NodeType: message.Controller,
		Event:    event,
		Time:     time.Now(),
	}
	// no one to send
	//c.statusReportChannel <- statusReport
	if r := mqueue.GetQueue("reporter"); r != nil {
		mqueue.SendMessage(statusReport, "reporter")
	}
}

func (c *Controller) StatusReport() <-chan message.StatusReport {
	return c.statusReportChannel
}

func calcRealMaxUeNum(n2IpNum, n3IpNum, uePerGnb, supiNum uint32, ueNum uint32) uint32 {
	t := u32Min(n2IpNum, n3IpNum) * uePerGnb
	return u32Min(u32Min(t, supiNum), ueNum)
}

func calcSUPINum(supiFrom string, supiTo string) uint32 {
	imsiFrom, err := strconv.ParseUint(strings.Split(supiFrom, "-")[1], 16, 64)
	if err != nil {
		logger.ControllerLog.Errorf("GetSUPI parse failed: %+v", err)
	}
	imsiTo, err := strconv.ParseUint(strings.Split(supiTo, "-")[1], 16, 64)
	if err != nil {
		logger.ControllerLog.Errorf("GetSUPI parse failed: %+v", err)
	}
	return uint32(imsiTo - imsiFrom)
}

func u32Min(a, b uint32) uint32 {
	if a < b {
		return a
	} else {
		return b
	}
}
