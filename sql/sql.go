package sql

import (
	"database/sql"
	"fmt"

	"time"

	"github.com/captncraig/gfig"
)

type sqlBackend struct {
	d            *sql.DB
	appName      string
	overrides    gfig.TaggedValues
	hash         uint32
	subs         []func(gfig.TaggedValues)
	errSubs      []func(error)
	pollInterval time.Duration
	trigger      chan bool
}

func New(d *sql.DB, appName string) (gfig.Backend, error) {
	b := &sqlBackend{
		d:            d,
		appName:      appName,
		pollInterval: 30 * time.Second,
		trigger:      make(chan bool),
	}
	if err := b.initSchema(); err != nil {
		return nil, err
	}
	ch := make(chan error)
	go b.startPolling(ch)
	if err := <-ch; err != nil {
		return nil, err
	}
	return b, nil
}

func (s *sqlBackend) startPolling(started chan<- error) {
	var err error
	s.overrides, err = s.getInternal()
	s.hash = s.overrides.Hash()
	started <- err
	if err != nil {
		return
	}
	for {
		select {
		case <-s.trigger:
		case <-time.After(s.pollInterval):
		}
		o, err := s.getInternal()
		if err != nil {
			s.onErr(err)
			continue
		}
		if o.Hash() != s.hash {
			s.overrides = o
			s.onChange()
		}
	}
}

func (s *sqlBackend) onChange() {
	for _, f := range s.subs {
		go f(s.overrides)
	}
}

func (s *sqlBackend) onErr(err error) {
	for _, f := range s.errSubs {
		go f(err)
	}
}

func (s *sqlBackend) getInternal() (gfig.TaggedValues, error) {
	const q = `select Datacenter, Value, Tier, Name from [dbo].[%s.settings]`
	query := fmt.Sprintf(q, s.appName)
	rows, err := s.d.Query(query)
	if err != nil {
		return nil, err
	}
	setts := gfig.TaggedValues{}
	defer rows.Close()
	for rows.Next() {
		setting := &gfig.TaggedValue{}
		err = rows.Scan(&setting.Datacenter, &setting.Value, &setting.Tier, &setting.Name)
		if err != nil {
			return nil, err
		}
		setts = append(setts, setting)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return setts, nil
}

func (s *sqlBackend) initSchema() error {
	tname := s.appName + ".settings"
	const q = `
IF NOT EXISTS 
(SELECT * FROM sys.tables t JOIN sys.schemas s on (t.schema_id = s.schema_id) WHERE s.name = 'dbo' and t.name = '%s') 
CREATE TABLE [dbo].[%s] ( 
	Name nvarchar(200) NOT NULL,
	Value nvarchar(max) NOT NULL,
	Datacenter nvarchar(200) NOT NULL,
	Tier nvarchar(200) NOT NULL,
	LastUpdated datetime NOT NULL
)`
	query := fmt.Sprintf(q, tname, tname)
	_, err := s.d.Exec(query)
	return err
}

func (s *sqlBackend) GetOverrides() (gfig.TaggedValues, error) {
	return s.overrides, nil
}

func (s *sqlBackend) Subscribe(sub func(gfig.TaggedValues), e func(error)) {
	if sub != nil {
		s.subs = append(s.subs, sub)
		go sub(s.overrides)
	}
	if e != nil {
		s.errSubs = append(s.errSubs, e)
	}
}

func (s *sqlBackend) SetOverride(t *gfig.TaggedValue) error {
	const q = `
IF EXISTS(SELECT * FROM [dbo].[%s.settings] WHERE Datacenter = ?1 AND Tier = ?2 AND Name = ?3)
BEGIN
    UPDATE [dbo].[%s.settings] SET Value = ?4, LastUpdated = GETUTCDATE()
    WHERE Datacenter = ?1 AND Tier = ?2 AND Name = ?3
END
ELSE
BEGIN
    INSERT INTO [dbo].[%s.settings] (Datacenter, Tier, Name, Value,LastUpdated)
    VALUES (?1,?2,?3,?4,GETUTCDATE())
END`
	query := fmt.Sprintf(q, s.appName, s.appName, s.appName)
	_, err := s.d.Exec(query, t.Datacenter, t.Tier, t.Name, t.Value)
	if err == nil {
		select {
		case s.trigger <- true:
		default:
		}
	}
	return err
}

func (s *sqlBackend) ClearOverride(t *gfig.TaggedValue) error {
	const q = "DELETE FROM [dbo].[%s.settings] WHERE Datacenter = ?1 AND Tier = ?2 AND Name = ?3"
	_, err := s.d.Exec(fmt.Sprintf(q, s.appName), t.Datacenter, t.Tier, t.Name)
	if err == nil {
		select {
		case s.trigger <- true:
		default:
		}
	}
	return err
}
