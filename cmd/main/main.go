package main

import (
	"flag"
	"github.com/mayfield-z/ember/internal/pkg/controller"
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Ctrl+C handler
	CloseHandler()

	configPath := flag.String("c", "../config/example.toml", "config path")
	//configPath := "../config/example.toml"
	flag.Parse()
	c := controller.Self()
	err := c.Init(*configPath)
	if err != nil {
		logger.AppLog.Errorf("config init failed")
	} else {
		c.Start()
		select {}
	}
}

func CloseHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-c
		logger.AppLog.Info("shutdown signal received, exiting...")
		controller.Self().Stop()
		os.Exit(0)
	}()
}
