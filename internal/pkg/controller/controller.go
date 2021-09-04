package controller

import (
	"github.com/mayfield-z/ember/internal/pkg/context"
	"sync"
)

var (
	controller = Controller{}
)

type Controller struct {
	gnbList []*context.GNB
	mutex   sync.Mutex
}

func ControllerSelf() *Controller {
	return &controller
}

func (c *Controller) AddGnb(gnb *context.GNB) {
	c.mutex.Lock()
	c.gnbList = append(c.gnbList, gnb)
	c.mutex.Unlock()
}

func (receiver) name() {

}
