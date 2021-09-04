package mqueue

import "sync"

var QueueMap sync.Map // map[string]*Queue

type Queue struct {
	Name    string
	Message chan interface{}
}

func NewQueue(name string) *Queue {
	q := &Queue{
		Name:    name,
		Message: make(chan interface{}),
	}
	QueueMap.Store(name, q)
	return q
}

func GetQueueByName(name string) *Queue {
	if q, ok := QueueMap.Load(name); ok {
		return q.(*Queue)
	} else {
		return nil
	}
}

func SendMessage(message interface{}, nodeName string) {
	queue := GetQueueByName(nodeName)
	queue.Message <- message
}

func GetMessageChan(nodeName string) chan interface{} {
	queue := GetQueueByName(nodeName)
	return queue.Message
}
