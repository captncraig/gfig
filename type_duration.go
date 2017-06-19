package gfig

import (
	"fmt"
	"time"
)

func (c *Collection) NewDuration(name string, defaultValue time.Duration, description string) *Duration {
	s := &Duration{
		value: defaultValue,
		Base: Base{
			rawValue: fmt.Sprint(defaultValue),
			meta: Metadata{
				Name:        name,
				Default:     fmt.Sprint(defaultValue),
				Description: description,
			},
		},
	}
	c.settings = append(c.settings, s)
	return s
}

func NewDuration(name string, defaultValue time.Duration, description string) *Duration {
	return defaultCollection.NewDuration(name, defaultValue, description)
}

type Duration struct {
	value time.Duration
	Base
}

func (d *Duration) Value() time.Duration {
	return d.value
}

func (d *Duration) OnChange(f func(time.Duration)) *Duration {
	d.addCallback(f)
	return d
}

func (d *Duration) Set(in string) (interface{}, error) {
	dur, err := time.ParseDuration(in)
	if err != nil {
		return nil, err
	}
	d.rawValue = in
	d.value = dur
	return dur, nil
}

func (d *Duration) SetDefault(t Tier, dc Datacenter, v time.Duration) *Duration {
	d.setDefault(t, dc, fmt.Sprint(v))
	return d
}
