package main

import (
	"database/sql"
	"flag"
	"fmt"
	"time"

	"net/http"

	"log"

	"github.com/captncraig/gfig"
	setsql "github.com/captncraig/gfig/sql"
	_ "github.com/denisenkom/go-mssqldb"
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

var x = gfig.NewString("AAAA.AAA", "FOOO", "AAAA").OnChange(func(s string) {
	fmt.Println("AAAAAA", s)
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

var conn = flag.String("db", "", "db connection string")

func main() {
	flag.Parse()
	db, err := sql.Open("mssql", *conn)
	if err != nil {
		log.Fatal(err)
	}
	bak, err := setsql.New(db, "MyTestApp")
	if err != nil {
		log.Fatal(err)
	}
	gfig.Init(gfig.Options{
		AppName:        "MyTestApp",
		Datacenter:     DCNY,
		Tier:           TDev,
		Backend:        bak,
		AllDatacenters: []gfig.Datacenter{DCAny, DCNY, DCCO},
	})
	go func() {
		for {
			fmt.Println("LLLL", x.Value())
			time.Sleep(time.Second)
		}
	}()
	http.ListenAndServe(":8080", gfig.Handler())
}
