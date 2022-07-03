package events

import (
	"errors"
	"sort"
	"strings"
)

type (
	Listener interface {
		Handle(event interface{}) interface{}
		Listen() string
		Priority() int
	}

	Dispatcher struct {
		Events map[string][]Listener
	}
)

func (d *Dispatcher) Register(listeners []Listener) {
	d.Events = make(map[string][]Listener)
	sort.Slice(listeners, func(i, j int) bool {
		return listeners[i].Priority() > listeners[j].Priority()
	})

	for k, listener := range listeners {
		if _, ok := d.Events[listener.Listen()]; !ok {
			d.Events[listener.Listen()] = make([]Listener, len(listeners))
		}

		d.Events[listener.Listen()][k] = listener
	}
}

func (d *Dispatcher) Dispatch(event string, payload interface{}) error {
	if _, ok := d.Events[event]; !ok {
		var message strings.Builder
		message.WriteString("event '")
		message.WriteString(event)
		message.WriteString("' not registered")

		return errors.New(message.String())
	}

	for _, listener := range d.Events[event] {
		listener.Handle(payload)
	}

	return nil
}
