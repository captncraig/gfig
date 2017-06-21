package gfig

import "fmt"

func (c *Collection) NewBool(name string, defaultValue bool, description string) *Bool {
	s := &Bool{
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

func NewBool(name string, defaultValue bool, description string) *Bool {
	return defaultCollection.NewBool(name, defaultValue, description)
}

type Bool struct {
	value bool
	Base
}

func (b *Bool) Value() bool {
	return b.value
}

func (b *Bool) OnChange(f func(bool)) *Bool {
	b.addCallback(f)
	return b
}

func (b *Bool) Set(in string) (interface{}, error) {
	v := false
	if in == "true" {
		v = true
	} else if in == "false" {
		v = false
	} else {
		return nil, fmt.Errorf("Bool must be 'true' or 'false'")
	}
	b.rawValue = in
	b.value = v
	return v, nil
}

func (d *Bool) SetDefault(t Tier, dc Datacenter, v bool) *Bool {
	d.setDefault(t, dc, fmt.Sprint(v))
	return d
}

func (b *Bool) Validate(raw string) error {
	if raw != "true" && raw != "false" {
		return fmt.Errorf("Bool must be 'true' or 'false'")
	}
	return nil
}
