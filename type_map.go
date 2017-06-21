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
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// add value to each item in place
	for i, k := range keys {
		keys[i] += " " + m[k]
	}
	return strings.Join(keys, "\n")
}

func parseMap(raw string) (map[string]string, error) {
	m := map[string]string{}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid map line %s. Needs exactly one space", line)
		}
		k, v := parts[0], parts[1]
		if _, ok := m[k]; ok {
			return nil, fmt.Errorf("key %s specified multiple times in map", k)
		}
		m[k] = v
	}
	return m, nil
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
	return parseMap(in)
}

func (m *Map) SetDefault(t Tier, dc Datacenter, v map[string]string) *Map {
	m.setDefault(t, dc, formatMap(v))
	return m
}

func (m *Map) Validate(raw string) error {
	_, err := parseMap(raw)
	return err
}
