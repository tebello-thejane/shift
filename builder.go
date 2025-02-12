package shift

import (
	"github.com/luno/reflex/rsql"
)

// TODO(corver): Possibly support explicit shifting to status X from different
//  statuses (Y and Z) each with different requests (XFromYReq, XFromZReq).

type option func(*FSM)

// WithMetadata provides an option to enable event metadata with a FSM.
func WithMetadata() option {
	return func(fsm *FSM) {
		fsm.withMetadata = true
	}
}

// NewFSM returns a new FSM builder.
func NewFSM(events rsql.EventsTableInt, opts ...option) initer {
	fsm := FSM{
		states: make(map[Status]status),
		events: events,
	}

	for _, opt := range opts {
		opt(&fsm)
	}

	return initer(builder(fsm))
}

type builder FSM

type initer builder

// Insert returns a FSM builder with the provided insert status.
func (c initer) Insert(st Status, inserter Inserter, next ...Status) builder {
	c.states[st] = status{
		st:     st,
		req:    inserter,
		t:      st,
		insert: false,
		next:   toMap(next),
	}
	c.insertStatus = st
	return builder(c)
}

// Update returns a FSM builder with the provided status update added.
func (b builder) Update(st Status, updater Updater, next ...Status) builder {
	if _, has := b.states[st]; has {
		// Ok to panic since it is build time.
		panic("state already added")
	}
	b.states[st] = status{
		st:     st,
		req:    updater,
		t:      st,
		insert: false,
		next:   toMap(next),
	}
	return b
}

// Build returns the built FSM.
func (b builder) Build() *FSM {
	fsm := FSM(b)
	return &fsm
}

func toMap(sl []Status) map[Status]bool {
	m := make(map[Status]bool)
	for _, s := range sl {
		m[s] = true
	}
	return m
}
