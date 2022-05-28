package events

import (
	"errors"
	"fmt"
	"sort"

	"github.com/KejawenLab/bima/v2/configs"
)

type Dispatcher struct {
	Events map[string][]configs.Listener
}

func (d *Dispatcher) Register(listeners []configs.Listener) {
	d.Events = make(map[string][]configs.Listener)
	sort.Slice(listeners, func(i, j int) bool {
		return listeners[i].Priority() > listeners[j].Priority()
	})

	for _, listener := range listeners {
		if _, ok := d.Events[listener.Listen()]; !ok {
			d.Events[listener.Listen()] = []configs.Listener{}
		}

		d.Events[listener.Listen()] = append(d.Events[listener.Listen()], listener)
	}
}

func (d *Dispatcher) Dispatch(event string, payload interface{}) error {
	if _, ok := d.Events[event]; !ok {
		return errors.New(fmt.Sprintf("Event '%s' not registered", event))
	}

	for _, listener := range d.Events[event] {
		listener.Handle(payload)
	}

	return nil
}
