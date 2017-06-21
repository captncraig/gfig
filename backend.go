package gfig

import (
	"sync"
)

type Backend interface {
	// Get all overrides now. Try not to rely on cache.
	GetOverrides() (TaggedValues, error)

	// Subscribe to override changes.
	Subscribe(func(TaggedValues), func(error))

	SetOverride(*TaggedValue) error

	ClearOverride(*TaggedValue) error
}

func NewMemoryStore() Backend {
	return &memory{}
}

type memory struct {
	overrides TaggedValues
	subs      []func(TaggedValues)
	sync.Mutex
}

func (m *memory) GetOverrides() (TaggedValues, error) {
	return m.overrides, nil
}

// Subscribe to override changes.
func (m *memory) Subscribe(f func(TaggedValues), _ func(error)) {
	m.Lock()
	if f != nil {
		m.subs = append(m.subs, f)
	}
	m.Unlock()
}

func (m *memory) SetOverride(v *TaggedValue) error {
	m.Lock()
	m.overrides = append(m.overrides, v)
	m.publish()
	m.Unlock()
	return nil
}

//only call in locked context
func (m *memory) publish() {
	for _, sub := range m.subs {
		go sub(m.overrides)
	}
}

func (m *memory) ClearOverride(tv *TaggedValue) error {
	m.Lock()
	newO := make([]*TaggedValue, 0, len(m.overrides))
	for _, o := range m.overrides {
		if tv.Datacenter == o.Datacenter && tv.Tier == o.Tier {
			continue
		}
		newO = append(newO, o)
	}
	m.overrides = newO
	m.publish()
	m.Unlock()
	return nil
}
