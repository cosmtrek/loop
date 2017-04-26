package monitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/message"
	"github.com/go-ini/ini"
	"github.com/juju/errors"
)

var (
	name = "monitor"
)

// Monitor ...
type Monitor struct {
	name     string
	app      string
	Health   string
	Interval time.Duration
	Timeout  time.Duration
	client   http.Client
}

// Option for Monitor
type Option struct {
	App      string
	Health   string
	Interval time.Duration
	Timeout  time.Duration
}

// NewMonitor returns monitor
func NewMonitor(opt *Option) (*Monitor, error) {
	m := new(Monitor)
	m.name = name
	m.app = opt.App
	m.Health = opt.Health
	m.Interval = opt.Interval
	m.Timeout = opt.Timeout
	m.client = http.Client{
		Timeout: m.Timeout,
	}
	return m, nil
}

// NewOption returns option
func NewOption(config *ini.File, app string) (*Option, error) {
	opt := new(Option)
	opt.App = app
	sec := config.Section(fmt.Sprintf("%s_%s", app, name))
	opt.Health = sec.Key("health").String()
	itv, err := time.ParseDuration(sec.Key("interval").String())
	if err != nil {
		return nil, errors.Trace(err)
	}
	opt.Interval = itv
	to, err := time.ParseDuration(sec.Key("timeout").String())
	if err != nil {
		return nil, errors.Trace(err)
	}
	opt.Timeout = to
	return opt, nil
}

// Event ...
func (m *Monitor) Event(q *chan *message.Message) {
	t := time.NewTicker(m.Interval)
	for {
		select {
		case <-t.C:
			logrus.Debug("monitor.Event")
			var err error
			resp, err := m.client.Get(m.Health)
			if err != nil {
				*q <- message.NewMessage(false, "", err)
				break
			}
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				*q <- message.NewMessage(false, "", err)
				break
			}
			if resp.StatusCode == http.StatusOK {
				*q <- message.NewMessage(true, string(body), nil)
			} else {
				*q <- message.NewMessage(false, string(body), nil)
			}
		}
	}
}

// Name returns plugin name
func (m *Monitor) Name() string {
	return m.name
}

// App returns app name
func (m *Monitor) App() string {
	return m.app
}
