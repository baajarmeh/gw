package gw

import (
	"fmt"
	"github.com/oceanho/gw/logger"
	"sync"
	"time"
)

type IEvent interface {
	MetaInfo() EventMetaInfo
}

type EventMetaInfo struct {
	Name     string
	Category string
	Data     interface{}
}

type EventHandler func(event IEvent) error

type IEventManager interface {
	Publish(event IEvent) error
	Subscribe(eventName string, handler EventHandler) (subscriberId string)
	Unsubscribe(subscriberId string)
}

//
// Default impl
//
type EventSubscriber struct {
	SubscriberId    string
	Handler         EventHandler
	hasUnsubscribed bool
}

var (
	ErrorEventChannelHasNotReady = fmt.Errorf("event channel not ready, may be has closed")
)

type DefaultEventManagerImpl struct {
	locker      sync.Mutex
	isReady     bool
	eventChan   chan IEvent
	idGenerator IdentifierGenerator
	subscribers map[string][]*EventSubscriber
}

func (d *DefaultEventManagerImpl) Publish(event IEvent) error {
	if !d.isReady {
		return ErrorEventChannelHasNotReady
	}
	d.eventChan <- event
	return nil
}

func (d *DefaultEventManagerImpl) Subscribe(eventName string, handler EventHandler) (subscriberId string) {
	d.locker.Lock()
	defer d.locker.Unlock()
	if _, ok := d.subscribers[eventName]; !ok {
		d.subscribers[eventName] = make([]*EventSubscriber, 0, 12)
	}
	es := &EventSubscriber{
		SubscriberId: d.idGenerator.NewStrID(),
		Handler:      handler,
	}
	d.subscribers[eventName] = append(d.subscribers[eventName], es)
	return es.SubscriberId
}

func (d *DefaultEventManagerImpl) Unsubscribe(subscriberId string) {
	d.locker.Lock()
	defer d.locker.Unlock()
	for k, _d := range d.subscribers {
		for i := 0; i < len(_d); i++ {
			if d.subscribers[k][i].SubscriberId == subscriberId {
				d.subscribers[k][i].hasUnsubscribed = true
				break
			}
		}
	}
}

func DefaultEventManager(state *ServerState) IEventManager {
	var m = &DefaultEventManagerImpl{
		isReady:     true,
		eventChan:   make(chan IEvent),
		idGenerator: state.IDGenerator(),
		subscribers: make(map[string][]*EventSubscriber),
	}
	var quit = make(chan bool, 1)
	var cleanDuration = time.Minute * 5
	state.s.RegisterShutDownHandler(func(s *HostServer) error {
		quit <- true
		return nil
	})
	go func() {
		var isQuit bool
		defer func() {
			close(quit)
			m.isReady = false
			close(m.eventChan)
		}()
		for {
			select {
			case <-quit:
				isQuit = true
				break
			case <-time.After(cleanDuration):
				func() {
					m.locker.Lock()
					defer m.locker.Unlock()
					var subs = make(map[string][]*EventSubscriber)
					for k, d := range m.subscribers {
						var subscribers = make([]*EventSubscriber, 0, len(d))
						for _, sub := range d {
							if sub.hasUnsubscribed {
								continue
							}
							_sub := sub
							subscribers = append(subscribers, _sub)
						}
						subs[k] = subscribers
					}
					m.subscribers = subs
				}()
				break
			case event := <-m.eventChan:
				metaInfo := event.MetaInfo()
				if subscribers, ok := m.subscribers[metaInfo.Name]; ok {
					for _, sub := range subscribers {
						if !sub.hasUnsubscribed {
							_ = sub.Handler(event)
						}
					}
				}
			}
			if isQuit {
				break
			}
		}
		logger.Info("quit event manager.")
	}()
	return m
}
