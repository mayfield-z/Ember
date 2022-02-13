package main

import (
	"flag"
	"github.com/mayfield-z/ember/internal/pkg/controller"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/reporter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	reporter.Self().Start()
	controller.Self().Start()
	select {}
}

func exitHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-c
		logger.AppLog.Info("shutdown signal received, exiting...")
		controller.Self().Stop()
		os.Exit(0)
	}()
}

func init() {
	// Ctrl+C handler
	exitHandler()

	configPath := flag.String("c", "../config/example.toml", "config path")
	//configPath := "../config/example.toml"
	flag.Parse()
	viper.SetConfigFile(*configPath)
	err := viper.ReadInConfig()
	if err != nil {
		logger.AppLog.Errorf("config read failed, exiting...")
		os.Exit(1)
	}

	level, err := logrus.ParseLevel(viper.GetString("app.logLevel"))
	if err != nil {
		level, _ = logrus.ParseLevel("info")
	}
	logger.SetLogLevel(level)
	logger.AppLog.Infof("Set log level to: %v", level)

	c := controller.Self()
	err = c.Init()
	if err != nil {
		logger.AppLog.Errorf("controller init failed, exiting...")
		os.Exit(1)
	}

	r := reporter.Self()
	err = r.Init()
	if err != nil {
		logger.AppLog.Errorf("reporter init failed, exiting...")
		os.Exit(1)
	}
}
