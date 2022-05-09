package controller

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/mayfield-z/ember/internal/pkg/gnb"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/mayfield-z/ember/internal/pkg/ue"
	"github.com/mayfield-z/ember/internal/pkg/utils"
	"github.com/mdlayher/arp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math"
	"net"
	"os"
	"path"
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
	userPlaneRuleList          []UserPlaneRule
	templateGnb                *gnb.GNB
	templateUE                 *ue.UE
	n2OriginalAddresses        []net.Addr
	n3OriginalAddresses        []net.Addr
	n2AddedIp                  []net.IP
	n3AddedIp                  []net.IP
	gnbName                    string
	n2Interface                *net.Interface
	n2IpFrom                   net.IP
	n2IpTo                     net.IP
	n2IpNum                    uint32
	n2IpPointer                net.IP
	n3Interface                *net.Interface
	n3IpFrom                   net.IP
	n3IpTo                     net.IP
	n3IpNum                    uint32
	n3IpPointer                net.IP
	dnIp                       net.IP
	dnInterface                *net.Interface
	supiFrom                   string
	supiTo                     string
	supiNum                    uint32
	ueNum                      uint32
	uePerGnb                   uint32
	uePerSec                   float64
	ueTimeout                  time.Duration
	realMaxUe                  uint32
	initPDUWhenAllUERegistered bool
	emulateUserPlane           bool
	userPlaneRuleOutputFile    *os.File
	userPlaneEmulateUEMaxNum   uint32
	amfIp                      net.IP
	amfPort                    int
	upfIp                      net.IP
	upfPort                    int
	upfMac                     net.HardwareAddr

	initialed              bool
	running                bool
	supiPointer            string
	ueIdPointer            uint64
	pduSessionIdPointer    uint8
	globalRANNodeIDPointer uint32
	nrCellIdentityPointer  uint64

	ueNumOfCurrentGnb uint32 // atomic
	emulatedUeNum     uint32 // atomic
	gnbCounter        uint32 // atomic

	logger                       *logrus.Entry
	ctx                          context.Context
	cancelFunc                   context.CancelFunc
	wg                           *sync.WaitGroup
	creatingGnbWg                *sync.WaitGroup
	statusReportChannel          chan message.StatusReport
	n2AddedIpMutex               sync.Mutex
	n3AddedIpMutex               sync.Mutex
	ueListMutex                  sync.Mutex
	ueInfoMutex                  sync.Mutex
	userPlaneRuleMutex           sync.Mutex
	userPlaneRuleOutputFileMutex sync.Mutex
	gnbListMutex                 sync.Mutex
}

type UserPlaneRule struct {
	UplinkTEID uint32
	GNBN3IP    net.IP
	UPFN3IP    net.IP
	UEIP       net.IP
	DNIP       net.IP
}

func Self() *Controller {
	return &controller
}

func (c *Controller) addGnb(gnb *gnb.GNB) {
	c.gnbListMutex.Lock()
	defer c.gnbListMutex.Unlock()
	c.logger.Debugf("Adding gNB %v", gnb.NodeName())
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
	c.logger.Debugf("Adding UE: %v", ue.NodeName())
	c.ueList = append(c.ueList, ue)
}

func (c *Controller) Init() error {
	// do not change exec order
	c.logger = logger.ControllerLog
	c.logger.Debugf("Start inital controller")
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())
	c.globalRANNodeIDPointer = 1
	c.nrCellIdentityPointer = 1
	c.ueIdPointer = 0
	c.pduSessionIdPointer = 0
	c.wg = &sync.WaitGroup{}
	c.creatingGnbWg = &sync.WaitGroup{}

	err := c.parseConfig()
	if err != nil {
		c.logger.Errorf("Controller init failed: %+v", err)
	}

	c.logger.Infof("getting upf mac address from %v: %v", c.n3Interface.Name, c.upfIp)
	arpClient, err := arp.Dial(c.n3Interface)
	if err != nil {
		c.logger.Errorf("Controller init failed: %+v", err)
		return err
	}
	c.upfMac, err = arpClient.Resolve(c.upfIp)
	if err != nil {
		c.logger.Errorf("Controller init failed: %+v", err)
		return err
	}
	c.logger.Infof("upf mac address: %v", c.upfMac)
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
	c.logger.Infof("GetSUPI from: %v, to: %v, total: %v", c.supiFrom, c.supiTo, c.supiNum)

	c.n2Interface, err = net.InterfaceByName(viper.GetString("controller.n2Interface"))
	if err != nil {
		return errors.Wrap(err, "Controller GetN2InterfaceByName failed")
	}

	c.n3Interface, err = net.InterfaceByName(viper.GetString("controller.n2Interface"))
	if err != nil {
		return errors.Wrap(err, "Controller GetN3InterfaceByName failed")
	}

	c.dnInterface, err = net.InterfaceByName(viper.GetString("controller.dnInterface"))
	if err != nil {
		return errors.Wrap(err, "Controller GetDNInterfaceByName failed")
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
	c.logger.Infof("N2 IP From: %v, To: %v, total: %v", n2IpFrom, n2IpTo, n2IpNum)

	n3IpFrom := net.ParseIP(viper.GetString("controller.n3IpFrom"))
	n3IpTo := net.ParseIP(viper.GetString("controller.n3IpTo"))
	n3IpNum := binary.BigEndian.Uint32(n3IpTo[12:]) - binary.BigEndian.Uint32(n3IpFrom[12:]) + 1
	c.n3IpFrom = n3IpFrom
	c.n3IpPointer = n3IpFrom
	c.n3IpTo = n3IpTo
	c.n3IpNum = n3IpNum
	c.logger.Infof("N2 IP From: %v, To: %v, total: %v", n3IpFrom, n3IpTo, n3IpNum)

	c.dnIp = net.ParseIP(viper.GetString("controller.dnIp"))
	c.ueNum = viper.GetUint32("controller.ueNum")
	c.uePerSec = viper.GetFloat64("controller.uePerSec")
	c.uePerGnb = viper.GetUint32("controller.uePerGnb")
	c.ueTimeout = viper.GetDuration("controller.ueTimeout")
	c.realMaxUe = calcRealMaxUeNum(c.n2IpNum, c.n3IpNum, c.uePerGnb, c.supiNum, c.ueNum)
	c.initPDUWhenAllUERegistered = viper.GetBool("controller.initPDUWhenAllUERegistered")
	c.emulateUserPlane = viper.GetBool("controller.emulateUserPlane")
	if c.emulateUserPlane {
		err = os.MkdirAll(viper.GetString("controller.userPlaneRuleOutputFolder"), 0777)
		if err != nil {
			c.userPlaneRuleOutputFile, _ = os.OpenFile(path.Join(viper.GetString("controller.userPlaneRuleOutputFolder"), "userplane.rule"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		} else {
			c.userPlaneRuleOutputFile, _ = os.OpenFile(path.Join("userplane.rule"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		}
	}

	c.userPlaneEmulateUEMaxNum = viper.GetUint32("controller.userPlaneEmulateUEMaxNum")
	c.logger.Infof("Real UE max is: %v, register %v ue per second", c.realMaxUe, c.uePerSec)
	c.logger.Infof("%v UE per gNB, will use %v gnb", c.uePerGnb, math.Ceil(float64(c.realMaxUe)/float64(c.uePerGnb)))
	if c.initPDUWhenAllUERegistered {
		c.logger.Infof("Will inital PDU after all UE registered in core")
	} else {
		c.logger.Infof("Will inital PDU once UE registered in core")
	}

	c.amfIp = net.ParseIP(viper.GetString("amf.ip"))
	c.amfPort = viper.GetInt("amf.port")

	c.upfIp = net.ParseIP(viper.GetString("upf.ip"))
	c.upfPort = viper.GetInt("amf.port")
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
		c.logger.Errorf("Start controller before initial it")
		return
	}
	c.running = true
	go c.start()
}

func (c *Controller) start() {
	// TODO: UE init PDU
	c.logger.Infof("Controller start")
	ueIntervalInMicrosecond := int64(math.Ceil(float64(1000000) / c.uePerSec))
	ticker := time.NewTicker(time.Duration(ueIntervalInMicrosecond) * time.Microsecond)
	defer ticker.Stop()
	c.SendStatusReport(message.ControllerStart)

	for i := uint32(0); i < c.realMaxUe/c.uePerGnb+1; i++ {
		c.creatingGnbWg.Add(1)
		c.ueNumOfCurrentGnbClear()
		c.gnbCounterAdd1()
		// create gnb
		if len(c.gnbList) == int(c.n2IpNum) || len(c.gnbList) == int(c.n3IpNum) {
			c.logger.Errorf("can not create more gNB")
			return
		}
		c.SendStatusReport(message.EmulateGNB)
		g := c.createAndAddGnb()
		err := g.Run()
		if err != nil {
			c.logger.Errorf("create gnb failed: %+v", err)
			return
		}
		for !g.Connected() {
		}
		c.creatingGnbWg.Done()
	}
	c.creatingGnbWg.Wait()
	c.logger.Infof("gnb creating finished, create %v gnb", len(c.gnbList))
	for c.emulatedUeNum < c.realMaxUe {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.emulatedUeNumAdd1()
			c.wg.Add(1)
			c.creatingGnbWg.Wait()
			go c.emulateOneUEControlPlane(!c.initPDUWhenAllUERegistered)
		}
	}

	c.wg.Wait()
	if c.initPDUWhenAllUERegistered {
		for _, u := range c.ueList {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				c.wg.Add(1)
				go c.establishOneUEPDUSession(u)
			}
		}
	}
	c.logger.Infof("hi")

	// for debug
	c.wg.Add(1)
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
	c.userPlaneRuleOutputFile.Close()

	c.n2AddedIpMutex.Lock()
	defer c.n2AddedIpMutex.Unlock()
	for _, ip := range c.n2AddedIp {
		c.logger.Debugf("deleting N2 IP addresses: %v", ip)
		err := utils.DelIpFromInterface(ip, c.n2Interface)
		if err != nil {
			c.logger.Errorf("can not del N2 IP addresses: %v, %v", ip, err)
		}
	}

	c.n3AddedIpMutex.Lock()
	defer c.n3AddedIpMutex.Unlock()
	for _, ip := range c.n3AddedIp {
		c.logger.Debugf("deleting N3 IP addresses: %v", ip)
		err := utils.DelIpFromInterface(ip, c.n3Interface)
		if err != nil {
			c.logger.Errorf("can not del N3 IP addresses: %v, %v", ip, err)
		}
	}
}

func (c *Controller) createAndAddGnb() *gnb.GNB {
	c.logger.Infof("Creating gnb %v-%v", c.gnbName, len(c.gnbList))
	g := c.templateGnb.Copy(fmt.Sprintf("%v-%v", c.gnbName, len(c.gnbList))).
		SetGlobalRANNodeID(c.globalRANNodeIDPointer).
		SetNRCellIdentity(c.nrCellIdentityPointer).
		SetN2Addresses(c.n2IpPointer).
		SetN3Addresses(c.n3IpPointer)
	c.globalRANNodeIDPointer += 1
	c.nrCellIdentityPointer += 1

	err := utils.AddIpToInterface(c.n2IpPointer, c.n2Interface)
	if err != nil {
		c.logger.Fatalf("Add N2 Address: %v to interface %s failed: %v", c.n2IpPointer, c.n2Interface.Name, err)
	} else {
		c.logger.Debugf("Add N2 Address: %v to interface %s success", c.n2IpPointer, c.n2Interface.Name)
		c.n2AddedIpMutex.Lock()
		c.n2AddedIp = append(c.n2AddedIp, c.n2IpPointer)
		c.n2AddedIpMutex.Unlock()
	}

	pinger, err := ping.NewPinger(string(c.amfIp))
	if err != nil {
		c.logger.Errorf("cannot create pinger")
	}
	pinger.Count = 5
	err = pinger.Run()
	if err != nil {
		c.logger.Errorf("cannot ping from %v to %v", c.n2IpPointer, c.amfIp)
	}

	//err = utils.AddIpToInterface(c.n3IpPointer, c.n3Interface)
	//if err != nil {
	//	c.logger.Fatalf("Add N3 Address: %v to interface %s failed: %v", c.n3IpPointer, c.n3Interface.Name, err)
	//} else {
	//	c.logger.Debugf("Add N3 Address: %v to interface %s success", c.n3IpPointer, c.n3Interface.Name)
	//	c.n3AddedIpMutex.Lock()
	//	c.n3AddedIp = append(c.n3AddedIp, c.n3IpPointer)
	//	c.n3AddedIpMutex.Unlock()
	//}

	c.n2IpPointer = utils.Add1(c.n2IpPointer)
	//c.n3IpPointer = utils.Add1(c.n3IpPointer)
	c.addGnb(g)
	return g
}

func (c *Controller) createAndAddUE() *ue.UE {
	c.ueInfoMutex.Lock()
	c.logger.Infof("Creating UE: %v", c.supiPointer)
	//u := c.templateUE.Copy(c.supiPointer, c.ueIdPointer, c.pduSessionIdPointer)
	u := c.templateUE.Copy(c.supiPointer, c.ueIdPointer, 0, c.ueTimeout)
	imsi, err := strconv.ParseUint(strings.Split(c.supiPointer, "-")[1], 10, 64)
	if err != nil {
		c.logger.Errorf("create ue failed: %+v", err)
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

func (c *Controller) emulateOneUEControlPlane(setupPDUSession bool) {
	defer c.wg.Done()
	// create UE
	c.SendStatusReport(message.EmulateUE)
	currentUE := c.createAndAddUE()
	currentUE.Run()
	defer currentUE.Stop()

	c.logger.Debugf("UE %v is using No.%v GNB", currentUE.GetSUPI(), int(currentUE.GetID()/uint64(c.uePerGnb)))
	currentGNB := c.getGnbByIndex(int(currentUE.GetID() % uint64(len(c.gnbList))))
	if currentGNB == nil {
		c.logger.Logger.Panicf("GNB empty pointer")
	}

	// start RRCSetup
	currentUE.RRCSetupRequest(currentGNB)

	select {
	case <-currentUE.Ctx().Done():
		return
	case msg := <-currentUE.StatusReport():
		switch msg.Event {
		case message.UERegistrationSuccess:
			c.logger.Infof("UE %v Registration Success", currentUE.GetSUPI())
		case message.UERegistrationReject:
			// TODO: handle reject
			c.logger.Infof("UE %v Registration Reject", currentUE.GetSUPI())
		}
	}

	if setupPDUSession {
		c.wg.Add(1)
		success := c.establishOneUEPDUSession(currentUE)
		if success {
			c.logger.Debugf("UE IP is: %v, TEID is :%v", currentUE.GetIP(), currentGNB.FindUEBySUPI(currentUE.GetSUPI()).RawUplinkTEID)
		}
	}
}

func (c *Controller) establishOneUEPDUSession(u *ue.UE) bool {
	defer c.wg.Done()
	if !u.Running() {
		u.Run()
		defer u.Stop()
	}

	u.EstablishPDUSession(0)
	select {
	case <-u.Ctx().Done():
		return false
	case msg := <-u.StatusReport():
		switch msg.Event {
		case message.UEPDUSessionEstablishmentAccept:
			c.logger.Infof("UE %v PDU Session Establishment Accept", u.GetSUPI())
			if c.emulateUserPlane {
				g := c.getGnbByIndex(int(u.GetID()) % len(c.gnbList))
				gu := g.FindUEBySUPI(u.GetSUPI())
				c.userPlaneRuleOutputFileMutex.Lock()
				c.userPlaneRuleOutputFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", gu.TransportLayerGNBAddress, gu.TransportLayerUPFAddress, gu.UplinkTEID, u.GetIP(), c.dnIp, c.upfMac))
				c.userPlaneRuleOutputFileMutex.Unlock()
			}
			return true
		case message.UEPDUSessionEstablishmentReject:
			c.logger.Infof("UE %v PDU Session Establishment Reject", u.GetSUPI())
			return false
		default:
			c.logger.Errorf("get wrong status report when establishing UE %v PDU Session: %+v", u.GetSUPI(), msg)
			return false
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
	imsiFrom, err := strconv.ParseUint(strings.Split(supiFrom, "-")[1], 10, 64)
	if err != nil {
		logger.ControllerLog.Errorf("GetSUPI parse failed: %+v", err)
	}
	imsiTo, err := strconv.ParseUint(strings.Split(supiTo, "-")[1], 10, 64)
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
