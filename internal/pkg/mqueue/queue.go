package mqueue

import (
	"github.com/mayfield-z/ember/internal/pkg/logger"
	"sync"
)

var queueMap sync.Map // map[string]*Queue

type Queue struct {
	Name    string
	Message chan interface{}
}

func NewQueue(name string) *Queue {
	q := &Queue{
		Name:    name,
		Message: make(chan interface{}),
	}
	queueMap.Store(name, q)
	logger.QueueLog.Debugf("Add new queue: %v", name)
	return q
}

func DelQueue(name string) {
	// TODO: make sure all message is handled.
	queueMap.Delete(name)
	logger.QueueLog.Debugf("Del queue: %v", name)
}

func GetQueueByName(name string) *Queue {
	if q, ok := queueMap.Load(name); ok {
		return q.(*Queue)
	} else {
		return nil
	}
}

func SendMessage(message interface{}, nodeName string) {
	logger.QueueLog.Debugf("sending message %T to node \"%v\"", message, nodeName)
	queue := GetQueueByName(nodeName)
	queue.Message <- message
}

func GetMessageChan(nodeName string) chan interface{} {
	queue := GetQueueByName(nodeName)
	return queue.Message
}
