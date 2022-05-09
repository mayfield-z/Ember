package reporter

import (
	"context"
	"fmt"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/message"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

var reporter = Reporter{}

type Status struct {
	EmulatedUeNum                       int
	RegistrationSuccessUeNum            int
	RegistrationRejectUeNum             int
	PDUSessionEstablishmentSuccessUeNum int
	PDUSessionEstablishmentRejectUeNum  int

	EmulatedGnbNum     int
	SetupSuccessGnbNum int
	SetupRejectGnbNum  int
}

type Reporter struct {
	outputFolder      string
	outputFileName    string
	outputFile        *os.File
	recordGranularity time.Duration
	exportRawData     bool

	logger             *logrus.Entry
	statusChannel      chan interface{}
	statusBuffer       []message.StatusReport
	statusBufferMutex  sync.Mutex
	rawStatus          []message.StatusReport
	lastStatus         Status
	lastStatusMutex    sync.Mutex
	processedStatus    []Status
	initialed          bool
	recordStartChannel chan struct{}
	memStats           runtime.MemStats

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func Self() *Reporter {
	return &reporter
}

func (r *Reporter) Init() error {
	logger.ControllerLog.Debugf("Start inital controller")
	r.ctx, r.cancelFunc = context.WithCancel(context.Background())
	r.logger = logger.ReporterLog
	mqueue.NewQueue("reporter")
	r.statusChannel = mqueue.GetMessageChannel("reporter")
	err := r.parseConfig()
	r.recordStartChannel = make(chan struct{})
	if err != nil {
		return err
	}
	r.outputFile.WriteString(fmt.Sprintf("emulatedGNBNum,connectedGNBNum,emulatedUENum,registrationSuccessedUENum,PDUSessionEstablishedUENum,goroutineNum,alloc,sys,cpuUsage\n"))
	r.initialed = true
	return nil
}

func (r *Reporter) parseConfig() error {
	var err error
	r.outputFolder = viper.GetString("reporter.outputFolder")
	err = os.MkdirAll(r.outputFolder, os.ModePerm)
	if err != nil {
		logger.ControllerLog.Errorf("Failed to create dir %s", err)
		return err
	}
	r.outputFileName = viper.GetString("reporter.outputFileName")
	r.recordGranularity = viper.GetDuration("reporter.recordGranularity")
	r.exportRawData = viper.GetBool("reporter.exportRawData")
	r.outputFile, err = os.OpenFile(path.Join(r.outputFolder, r.outputFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.ControllerLog.Errorf("Failed to open file %s", err)
		return err
	}
	r.logger.Printf("Output folder: %v, output file name: %v, record granularity: %v", r.outputFolder, r.outputFileName, r.recordGranularity)
	return nil
}

func (r *Reporter) Start() {
	if !r.initialed {
		logger.ControllerLog.Errorf("Start reporter before initial it")
		return
	}
	go r.start()
}

func (r *Reporter) start() {
	go r.receiveStatusReport()
	<-r.recordStartChannel
	ticker := time.NewTicker(r.recordGranularity)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.processStatus()
		case <-r.ctx.Done():
			r.logger.Printf("Stop reporter")
			return
		}
	}
}

func (r *Reporter) Done() <-chan struct{} {
	return r.ctx.Done()
}

func (r *Reporter) processStatus() {
	statusBuffer := make([]message.StatusReport, 8)
	r.statusBufferMutex.Lock()
	r.lastStatusMutex.Lock()
	status := r.lastStatus
	r.lastStatusMutex.Unlock()
	statusBuffer = r.statusBuffer
	r.statusBuffer = nil
	r.statusBufferMutex.Unlock()

	for _, report := range statusBuffer {
		switch report.Event {
		case message.EmulateUE:
			status.EmulatedUeNum++
		case message.EmulateGNB:
			status.EmulatedGnbNum++
		case message.UERegistrationSuccess:
			status.RegistrationSuccessUeNum++
		case message.UERegistrationReject:
			status.RegistrationRejectUeNum++
		case message.UEPDUSessionEstablishmentAccept:
			status.PDUSessionEstablishmentSuccessUeNum++
		case message.UEPDUSessionEstablishmentReject:
			status.PDUSessionEstablishmentRejectUeNum++
		case message.GNBSetupSuccess:
			status.SetupSuccessGnbNum++
		case message.GNBSetupReject:
			status.SetupRejectGnbNum++
		}
	}
	r.processedStatus = append(r.processedStatus, status)
	r.lastStatusMutex.Lock()
	r.lastStatus = status
	r.lastStatusMutex.Unlock()
	runtime.ReadMemStats(&r.memStats)
	cmd := fmt.Sprintf("top -b -n 2 -d 0.4 -p %v | tail -1 | awk '{print $9}'", os.Getpid())
	cpuUsageCmd := exec.Command("bash", "-c", cmd)
	cpuUsage, err := cpuUsageCmd.Output()
	if err != nil {
		logger.ControllerLog.Errorf("Failed to get pid: %v cpu usage %s. command is: %v", os.Getpid(), err, cpuUsageCmd.String())
	}
	_, err = r.outputFile.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
		status.EmulatedGnbNum,
		status.SetupSuccessGnbNum,
		status.EmulatedUeNum,
		status.RegistrationSuccessUeNum,
		status.PDUSessionEstablishmentSuccessUeNum,
		runtime.NumGoroutine(),
		r.memStats.Alloc/1024/1024,
		r.memStats.Sys/1024/1024,
		strings.TrimSpace(string(cpuUsage))),
	)
	if err != nil {
		r.logger.Errorf("Failed to write to file %s", err)
	}
}

func (r *Reporter) receiveStatusReport() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case status := <-r.statusChannel:
			if status.(message.StatusReport).Event == message.ControllerStart {
				r.recordStartChannel <- struct{}{}
			}
			if status.(message.StatusReport).Event == message.ControllerStop {
				r.Stop()
			}
			r.statusBufferMutex.Lock()
			r.statusBuffer = append(r.statusBuffer, status.(message.StatusReport))
			r.statusBufferMutex.Unlock()
			if r.exportRawData {
				r.rawStatus = append(r.rawStatus, status.(message.StatusReport))
			}
		}
	}
}

func (r *Reporter) Stop() {
	r.processStatus() // process remaining status
	r.cancelFunc()
}
