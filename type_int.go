package gfig

import (
	"fmt"
	"strconv"
)

func (c *Collection) NewInt(name string, defaultValue int, description string) *Int {
	s := &Int{
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

func NewInt(name string, defaultValue int, description string) *Int {
	return defaultCollection.NewInt(name, defaultValue, description)
}

type Int struct {
	value int
	Base
}

func (i *Int) Set(in string) (interface{}, error) {
	v, err := strconv.Atoi(in)
	if err != nil {
		return nil, err
	}
	i.value = v
	i.rawValue = in
	return v, nil
}

func (i *Int) Value() int {
	return i.value
}

func (i *Int) OnChange(f func(int)) *Int {
	i.addCallback(f)
	return i
}

func (i *Int) SetDefault(t Tier, d Datacenter, v int) *Int {
	i.setDefault(t, d, fmt.Sprint(v))
	return i
}
