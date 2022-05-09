package logger

import (
	"fmt"
	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"time"
)

var (
	log           *logrus.Logger
	AppLog        *logrus.Entry
	InitLog       *logrus.Entry
	CfgLog        *logrus.Entry
	ContextLog    *logrus.Entry
	NgapLog       *logrus.Entry
	SctpLog       *logrus.Entry
	GnbLog        *logrus.Entry
	UeLog         *logrus.Entry
	ControllerLog *logrus.Entry
	ReporterLog   *logrus.Entry
	QueueLog      *logrus.Entry
	DpdkLog       *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"category", "name", "part"},
	}

	AppLog = log.WithFields(logrus.Fields{"category": "App"})
	InitLog = log.WithFields(logrus.Fields{"category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"category": "Cfg"})
	ContextLog = log.WithFields(logrus.Fields{"category": "Context"})
	NgapLog = log.WithFields(logrus.Fields{"category": "Ngap"})
	SctpLog = log.WithFields(logrus.Fields{"category": "Sctp"})
	GnbLog = log.WithFields(logrus.Fields{"category": "Gnb"})
	UeLog = log.WithFields(logrus.Fields{"category": "Ue"})
	ControllerLog = log.WithFields(logrus.Fields{"category": "Controller"})
	ReporterLog = log.WithFields(logrus.Fields{"category": "Reporter"})
	QueueLog = log.WithFields(logrus.Fields{"category": "Queue"})
	DpdkLog = log.WithFields(logrus.Fields{"category": "DPDK"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetOutputFolder(folder string) {
	err := os.MkdirAll(folder, 0777)
	if err != nil {
		log.Errorf("Failed to create log folder: %s", err)
		return
	}
	logFile, err := os.OpenFile(path.Join(folder, fmt.Sprintf("%s.log", viper.GetString("startTime"))), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("Failed to open log file: %s", err)
		return
	}
	addOutput(logFile)
}

func addOutput(writer io.Writer) {
	mw := io.MultiWriter(writer, log.Out)
	log.SetOutput(mw)
}
