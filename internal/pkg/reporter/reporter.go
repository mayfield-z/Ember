package reporter

import (
	"context"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/sirupsen/logrus"
	"os"
)

var reporter = Reporter{}

type Reporter struct {
	Notify         chan interface{}
	outputPath     string
	outputFileName string
	outputFile     *os.File

	logger *logrus.Entry

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func Self() *Reporter {
	return &reporter
}

func (r *Reporter) Init() error {
	logger.ControllerLog.Debugf("Start inital controller")
	r.ctx, r.cancelFunc = context.WithCancel(context.Background())
	r.Notify = make(chan interface{})
	r.logger = logger.ReporterLog
	return nil
}

func (r *Reporter) parseConfig() error {
	return nil
}

func (r *Reporter) Start() {

}
