package main

import (
	"github.com/mayfield-z/ember/internal/pkg/controller"
	"github.com/mayfield-z/ember/internal/pkg/logger"
)

func main() {
	configPath := "/home/cnic/src/Ember/config/example.toml"
	c := controller.ControllerSelf()
	err := c.Init(configPath)
	if err != nil {
		logger.AppLog.Errorf("fuck")
		return
	}
	c.Start()
	select {}
}
