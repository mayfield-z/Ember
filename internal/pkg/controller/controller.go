package controller

import (
	"github.com/mayfield-z/ember/internal/pkg/gnb"
	"sync"
)

var (
	controller = Controller{}
)

type Controller struct {
	gnbList []*gnb.GNB
	mutex   sync.Mutex
}

func ControllerSelf() *Controller {
	return &controller
}

func (c *Controller) AddGnb(gnb *gnb.GNB) {
	c.mutex.Lock()
	c.gnbList = append(c.gnbList, gnb)
	c.mutex.Unlock()
}
