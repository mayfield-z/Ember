package logger

import (
	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
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
	QueueLog      *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"category", "name"},
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
	QueueLog = log.WithFields(logrus.Fields{"category": "Queue"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}
