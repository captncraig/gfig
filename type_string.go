package gfig

func (c *Collection) NewString(name string, defaultValue string, description string) *String {
	s := &String{
		value: defaultValue,
		Base: Base{
			rawValue: defaultValue,
			meta: Metadata{
				Name:        name,
				Default:     defaultValue,
				Description: description,
			},
		},
	}
	c.settings = append(c.settings, s)
	return s
}

func NewString(name string, defaultValue string, description string) *String {
	return defaultCollection.NewString(name, defaultValue, description)
}

type String struct {
	value string
	Base
}

func (s *String) Value() string {
	return s.value
}

func (s *String) Set(in string) (interface{}, error) {
	s.value = in
	s.rawValue = in
	return in, nil
}

func (s *String) SetDefault(t Tier, d Datacenter, v string) *String {
	s.setDefault(t, d, v)
	return s
}

func (s *String) OnChange(f func(string)) *String {
	s.addCallback(f)
	return s
}
