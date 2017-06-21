package gfig

import (
	"fmt"
	"hash/fnv"
	"log"
	"reflect"
)

// Metadata represents a setting's type and data for the ui to render. Has no value or context.
type Metadata struct {
	Name        string
	Default     string
	Defaults    []*TaggedValue
	Description string
}

// Setting represents an instance of a setting FOR THIS APP.
// All overrides will be applied in this apps tier/dc context
type Setting interface {
	Name() string
	Raw() string
	Metadata() Metadata
	Trigger(interface{})

	// Implemented by specific types:

	Set(raw string) (interface{}, error)
	Validate(raw string) error
}

type Base struct {
	rawValue  string
	meta      Metadata
	callbacks []reflect.Value
}

func (b *Base) Name() string {
	return b.meta.Name
}

func (b *Base) Raw() string {
	return b.rawValue
}

func (b *Base) Metadata() Metadata {
	return b.meta
}

func (b *Base) Trigger(v interface{}) {
	for _, cb := range b.callbacks {
		cb.Call([]reflect.Value{reflect.ValueOf(v)})
	}
}

func (b *Base) addCallback(f interface{}) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("Callbacks must be functions. Got a %T", f))
	}
	if t.NumIn() != 1 {
		panic("Callbacks must be functions with exactly one argument")
	}
	if t.NumOut() != 0 {
		panic("Callbacks must not have returns")
	}
	b.callbacks = append(b.callbacks, reflect.ValueOf(f))
}

func (b *Base) setDefault(t Tier, d Datacenter, value string) {
	b.meta.Defaults = append(b.meta.Defaults, &TaggedValue{Tier: t, Datacenter: d, Value: value, Name: b.Name()})
}

type TaggedValue struct {
	Datacenter Datacenter
	Name       string
	Tier       Tier
	Value      string
}

func (t *TaggedValue) Hash() string {
	return fmt.Sprintf("%s|%s|%s|%s", t.Datacenter, t.Tier, t.Name, t.Value)
}

type TaggedValues []*TaggedValue

func (t TaggedValues) Hash() uint32 {
	h := fnv.New32a()
	for _, tv := range t {
		h.Write([]byte(tv.Hash() + "!"))
	}
	return h.Sum32()
}

type Datacenter string
type Tier string

type Collection struct {
	Options
	settings []Setting
}

type Options struct {
	AppName        string
	Tier           Tier
	Datacenter     Datacenter
	AllDatacenters []Datacenter
	AllTiers       []Tier
	Backend        Backend
}

func (c *Collection) Init(opts Options) {
	c.Options = opts
	c.Backend.Subscribe(c.reconcile, c.onError)
}

func (c *Collection) getByName(name string) Setting {
	for _, s := range c.settings {
		if s.Name() == name {
			return s
		}
	}
	return nil
}

func Init(opts Options) {
	defaultCollection.Init(opts)
}

const (
	ANYDC   Datacenter = "Any"
	ANYTIER Tier       = "Any"
)

func (c *Collection) reconcile(overrides TaggedValues) {
	for _, setting := range c.settings {
		// base default
		chosenValue := setting.Metadata().Default
		var bestMatch Datacenter
		// specific defaults
		for _, tv := range setting.Metadata().Defaults {
			if tv.Name == setting.Name() && tv.Tier == c.Tier {
				if tv.Datacenter == c.Datacenter || (bestMatch == "" && tv.Datacenter == ANYDC) {
					chosenValue = tv.Value
					bestMatch = tv.Datacenter
				}
			}
		}
		// overrides
		bestMatch = ""
		for _, tv := range overrides {
			if tv.Name == setting.Name() && tv.Tier == c.Tier {
				if tv.Datacenter == c.Datacenter || (bestMatch == "" && tv.Datacenter == ANYDC) {
					chosenValue = tv.Value
					bestMatch = tv.Datacenter
				}
			}
		}
		// see if there's any difference
		if chosenValue != setting.Raw() {
			// change!
			newVal, err := setting.Set(chosenValue)
			if err != nil {
				c.onError(err)
				continue
			}
			setting.Trigger(newVal)
		}
	}
}

func (c *Collection) onError(err error) {
	log.Println("ERROR in SETTINGS:", err)
}

var defaultCollection = Collection{}
