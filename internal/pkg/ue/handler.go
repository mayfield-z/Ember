package ue

import (
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"github.com/mayfield-z/ember/internal/pkg/mqueue"
)

func (u *UE) messageHandler() {
	for {
		select {
		case <-u.ctx.Done():
			return
		default:
			c := u.getMessageChan()
			select {
			case msg := <-c:
				switch msg.(type) {
				case mqueue.RRCSetupMessage:
					u.handleRRCSetupMessage(msg)
				case mqueue.RRCRejectMessage:
					u.handleRRCRejectMessage(msg)
				default:
					logger.UeLog.Errorf("Type %T message not supported by ue", msg)
				}
			}
		}
	}
}

func (u *UE) handleRRCSetupMessage(msg interface{}) {
	u.rrcFSM.Event(eventRRCSetup)
	u.sendRRCSetupComplete(u.gnb.Name)
}

func (u *UE) handleRRCRejectMessage(msg interface{}) {
	// implement when needed.
	panic("not implement")
}
func (u *UE) sendRRCSetupComplete(name string) {
	nas, err := u.buildRegistrationRequest()
	if err != nil {
		logger.UeLog.Errorf("Build RegistrationRequest failed: %+v", err)
		return
	}
	msg := mqueue.RRCSetupCompleteMessage{
		NASRegistrationRequest: nas,
		SendBy:                 u.supi,
		PLMN:                   u.plmn,
	}
	mqueue.SendMessage(msg, name)
}
