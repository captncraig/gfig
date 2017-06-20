package gfig

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Collection) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/settings.js" {
			http.ServeFile(w, r, "settings.js")
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(indexHtml))
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/get" {
			c.serveJson(w)
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == "/set" {
			c.addOverride(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

func Handler() http.Handler {
	return defaultCollection.Handler()
}

func (c *Collection) addOverride(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("settingName")
	dc := Datacenter(r.FormValue("dataCenter"))
	val := r.FormValue("value")
	tv := &TaggedValue{
		Datacenter: dc,
		Name:       name,
		Tier:       c.Tier,
		Value:      val,
	}
	err := c.Backend.SetOverride(tv)
	if err != nil {
		c.onError(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (c *Collection) serveJson(w http.ResponseWriter) {
	setts := &jsonSettings{
		CurrentTier:          c.Tier,
		AvailibleDatacenters: c.AllDatacenters,
	}
	allow := map[Datacenter]bool{}
	for _, dc := range c.AllDatacenters {
		allow[dc] = true
	}
	for _, s := range c.settings {
		m := s.Metadata()
		js := &jsonSetting{
			TypeName:    fmt.Sprintf("%T", s),
			Name:        s.Name(),
			Description: m.Description,
			DefaultValue: &jsonOverride{
				IsDefault:  true,
				Tier:       ANYTIER,
				DataCenter: ANYDC,
				Value:      m.Default,
			},
			AllowsOverrides: allow,
			AllDefaults:     []*jsonOverride{},
			AllOverrides:    []*jsonOverride{},
		}
		setts.Settings = append(setts.Settings, js)
		for _, d := range s.Metadata().Defaults {
			js.AllDefaults = append(js.AllDefaults, &jsonOverride{
				Tier:       d.Tier,
				DataCenter: d.Datacenter,
				Name:       d.Name,
				Value:      d.Value,
				IsDefault:  true,
			})
		}
		os, err := c.Backend.GetOverrides()
		if err != nil {
			c.onError(err)
			http.Error(w, err.Error(), 500)
			return
		}
		for _, o := range os {
			if o.Name != m.Name {
				continue
			}
			js.AllOverrides = append(js.AllOverrides, &jsonOverride{
				Tier:       o.Tier,
				DataCenter: o.Datacenter,
				Name:       o.Name,
				Value:      o.Value,
				IsOverride: true,
			})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(setts)
}

type jsonSettings struct {
	Settings             []*jsonSetting `json:"settings"`
	CurrentTier          Tier           `json:"currentTier"`
	AvailibleDatacenters []Datacenter   `json:"availableDataCenters"`
}

type jsonSetting struct {
	RequiresRestart bool                `json:"requiresRestart"`
	TypeName        string              `json:"typeName"`
	AllDefaults     []*jsonOverride     `json:"allDefaults"`
	DefaultValue    *jsonOverride       `json:"defaultValue"`
	AllOverrides    []*jsonOverride     `json:"allOverrides"`
	Description     string              `json:"description"`
	Name            string              `json:"name"`
	AllowsOverrides map[Datacenter]bool `json:"allowsOverrides"`
}

type jsonOverride struct {
	IsOverride      bool       `json:"isOverride"`
	HasTier         bool       `json:"hasTier"`
	HasDataCenter   bool       `json:"hasDataCenter"`
	AllowsOverrides bool       `json:"allowsOverrides"`
	IsDefault       bool       `json:"isDefault"`
	Tier            Tier       `json:"tier"`
	DataCenter      Datacenter `json:"dataCenter"`
	Value           string     `json:"value"`
	Name            string     `json:"name"`
}

const indexHtml = `<!DOCTYPE html>
<html>
<head>

    
<script src="https://cdn.sstatic.net/clc/admin/settings-panel.min.js?v=149fe996d3f8"></script>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
</head>
<body>
<div class="container">
<div id="settings"></div>
</div>
<script>
	 window.settingsList = SettingsPanel.init(
        document.getElementById('settings'),
        {
            settingsUrl: 'get',
            setUrl: 'set',
            clearUrl: 'clear',
			
        }
    );
	</script>
</body>
</html>`
