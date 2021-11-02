package main

import (
	"flag"
	"github.com/mayfield-z/ember/internal/pkg/controller"
	"github.com/mayfield-z/ember/internal/pkg/logger"
)

func main() {
	configPath := flag.String("c", "../config/example.toml", "config path")
	//configPath := "../config/example.toml"
	c := controller.ControllerSelf()
	err := c.Init(*configPath)
	if err != nil {
		logger.AppLog.Errorf("config init failed")
		return
	}
	c.Start()
	select {}
}
