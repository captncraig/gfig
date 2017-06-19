package main

import (
	"database/sql"
	"time"

	"net/http"

	"github.com/captncraig/gfig"
)

// make some simple fields
var myStringVar = gfig.NewString("core.UserAgent", "abc", "User agent to use for http requests or whatever")
var myIntVar = gfig.NewInt("core.MaxConns", 42, "Maximum number of outbound connections. **Markdown** in descriptions works!")

// define multiline markdown description!
var descDuration = `How often to poll the db for changes

Newlines work

 - So
 - Do
 - Bullets.`
var myDurationVar = gfig.NewDuration("core.PollingInterval", 5*time.Minute, descDuration)

// Set Specific defaults per tier/dc.
// And callback functions
var db *sql.DB
var maxActiveConns = gfig.NewInt("db.MaxActiveConns", 10, "Maximum open db conns").
	SetDefault(TDev, DCAny, 20).
	SetDefault(TProd, DCAny, 100).
	OnChange(func(i int) {
		if db != nil {
			db.SetMaxOpenConns(i)
		}
	})

var maxIdleConns = gfig.NewInt("db.MaxIdleConns", 10, "Max idle sql connections").
	OnChange(func(i int) {
		if db != nil {
			db.SetMaxIdleConns(i)
		}
	})

var sqlMaxLife = gfig.NewDuration("db.MaxLifetime", time.Minute, "Time after which to close db connections").
	OnChange(func(d time.Duration) {
		if db != nil {
			db.SetConnMaxLifetime(d)
		}
	})

const (
	DCAny gfig.Datacenter = "Any"
	DCNY  gfig.Datacenter = "New York"
	DCCO  gfig.Datacenter = "Colorado"
)

const (
	TAny   gfig.Tier = "Any"
	TLocal gfig.Tier = "Local"
	TDev   gfig.Tier = "Dev"
	TProd  gfig.Tier = "Prod"
)

func main() {
	gfig.Init(gfig.Options{
		AppName:        "MyTestApp",
		Datacenter:     DCNY,
		Tier:           TDev,
		Backend:        gfig.NewMemoryStore(),
		AllDatacenters: []gfig.Datacenter{DCAny, DCNY, DCCO},
	})
	http.ListenAndServe(":8080", gfig.Handler())
}
