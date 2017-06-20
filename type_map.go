package gfig

import (
	"fmt"
	"sort"
	"strings"
)

func (c *Collection) NewMap(name string, defaultValue map[string]string, description string) *Map {
	s := &Map{
		value: defaultValue,
		Base: Base{
			rawValue: fmt.Sprint(defaultValue),
			meta: Metadata{
				Name:        name,
				Default:     formatMap(defaultValue),
				Description: description,
			},
		},
	}
	c.settings = append(c.settings, s)
	return s
}

func formatMap(m map[string]string) string {
	s := ""
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		keys[i] += " " + m[k]
	}
	return strings.Join(keys, "\n")
}

func NewMap(name string, defaultValue map[string]string, description string) *Map {
	return defaultCollection.NewMap(name, defaultValue, description)
}

type Map struct {
	value map[string]string
	Base
}

func (m *Map) Value() map[string]string {
	return m.value
}

func (m *Map) OnChange(f func(map[string]string)) *Map {
	m.addCallback(f)
	return m
}

func (m *Map) Set(in string) (interface{}, error) {

	return in, nil
}

func (m *Map) SetDefault(t Tier, dc Datacenter, v map[string]string) *Map {
	m.setDefault(t, dc, formatMap(v))
	return m
}
