package gfig

import "strings"

func (c *Collection) NewList(name string, defaultValue []string, description string) *List {
	s := &List{
		value: defaultValue,
		Base: Base{
			rawValue: formatList(defaultValue),
			meta: Metadata{
				Name:        name,
				Default:     formatList(defaultValue),
				Description: description,
			},
		},
	}
	c.settings = append(c.settings, s)
	return s
}

func formatList(l []string) string {
	return strings.Join(l, "\n")
}

func NewList(name string, defaultValue []string, description string) *List {
	return defaultCollection.NewList(name, defaultValue, description)
}

type List struct {
	value []string
	Base
}

func (l *List) Value() []string {
	return l.value
}

func (l *List) Set(in string) (interface{}, error) {
	if in == "" {
		l.value = nil
	}
	l.value = strings.Split(in, "\n")
	l.rawValue = in
	return in, nil
}

func (l *List) SetDefault(t Tier, d Datacenter, v []string) *List {
	l.setDefault(t, d, formatList(v))
	return l
}

func (l *List) OnChange(f func([]string)) *List {
	l.addCallback(f)
	return l
}

func (l *List) Validate(raw string) error {
	return nil
}
