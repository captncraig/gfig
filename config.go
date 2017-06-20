package gfig

import (
	"fmt"
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
	Set(raw string) (interface{}, error)
	Trigger(interface{})
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
	c.reconcile(nil)
	c.Backend.Subscribe(c.reconcile, c.onError)
}

func Init(opts Options) {
	defaultCollection.Init(opts)
}

const (
	ANYDC   Datacenter = "Any"
	ANYTIER Tier       = "Any"
)

func (c *Collection) isMatch(tv *TaggedValue) bool {
	if tv.Tier == ANYTIER || tv.Tier == c.Tier {
		if tv.Datacenter == ANYDC || tv.Datacenter == c.Datacenter {
			return true
		}
	}
	return false
}

func (c *Collection) reconcile(overrides []*TaggedValue) {
	for _, setting := range c.settings {
		// base default
		chosenValue := setting.Metadata().Default
		// specific defaults
		for _, tv := range setting.Metadata().Defaults {
			if tv.Name == setting.Name() && c.isMatch(tv) {
				chosenValue = tv.Value
			}
		}
		// overrides
		for _, tv := range overrides {
			if tv.Name == setting.Name() && c.isMatch(tv) {
				chosenValue = tv.Value
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
