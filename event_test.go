package gw

import (
	"github.com/oceanho/gw/libs/gwjsoner"
	"testing"
	"time"
)

type TesterEvent struct {
}

func (t TesterEvent) MetaInfo() EventMetaInfo {
	return EventMetaInfo{
		Name:     "tester",
		Category: "test",
		Data:     "this is tester event",
	}
}

func TestDefaultEventManager(t *testing.T) {
	var server = NewTesterServer()
	var state = NewServerState(server)
	var eventManager = DefaultEventManager(state)
	var subscriberId = eventManager.Subscribe("tester", func(event IEvent) error {
		var b, _ = gwjsoner.Marshal(event.MetaInfo())
		t.Logf("got event, %s", string(b))
		return nil
	})
	var event = TesterEvent{}
	_ = eventManager.Publish(event)
	eventManager.Unsubscribe(subscriberId)
	_ = eventManager.Publish(event)
	server.ShutDown()
	time.Sleep(time.Second * 5)
}
