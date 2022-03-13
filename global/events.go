package global

import (
	"reflect"
	"sync"
)

type EventBus struct {
	subs map[reflect.Type][]chan interface{}
	rw   sync.RWMutex
}

type Subscription struct {
	eb    *EventBus
	typ   reflect.Type
	index int
	c     chan interface{}
}

func (s *Subscription) Chan() chan interface{} {
	return s.c
}

func (s *Subscription) Unsubscribe() {
	s.eb.unsubscribe(s.typ, s.index)
}

func NewEventBus() *EventBus {
	return &EventBus{
		subs: make(map[reflect.Type][]chan interface{}),
	}
}

func (e *EventBus) Subscript(t interface{}) *Subscription {
	e.rw.Lock()
	defer e.rw.Unlock()
	rtyp := reflect.TypeOf(t)
	subtion := &Subscription{
		typ: rtyp,
		c:   make(chan interface{}),
		eb:  e,
	}
	if prev, found := e.subs[rtyp]; found {
		nextIndex := len(prev)
		subtion.index = nextIndex
		e.subs[rtyp] = append(prev, subtion.c)
	} else {
		subtion.index = 0
		e.subs[rtyp] = append([]chan interface{}{}, subtion.c)
	}
	return subtion
}

func (e *EventBus) Publish(data interface{}) {
	e.rw.RLock()
	defer e.rw.RUnlock()
	rtyp := reflect.TypeOf(data)
	if cs, found := e.subs[rtyp]; found {
		go func(d interface{}, cs []chan interface{}) {
			for _, ch := range cs {
				ch <- d
			}
		}(data, cs)
	}
}

func (e *EventBus) unsubscribe(data interface{}, index int) {
	e.rw.Lock()
	defer e.rw.Unlock()
	rtyp := reflect.TypeOf(data)
	if old, found := e.subs[rtyp]; found {
		e.subs[rtyp] = append(old[:index], old[index+1:]...)
	}
}
