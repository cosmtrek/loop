package emitter

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/message"
	"github.com/go-ini/ini"
	"github.com/juju/errors"
)

var (
	name = "emitter"
)

// Emitter emits at certain interval
type Emitter struct {
	name     string
	app      string
	Interval time.Duration
}

// Option for Emitter
type Option struct {
	App      string
	Interval time.Duration
}

// NewEmitter returns emitter
func NewEmitter(opt *Option) (*Emitter, error) {
	emt := new(Emitter)
	emt.name = name
	emt.app = opt.App
	emt.Interval = opt.Interval
	return emt, nil
}

// NewEmitterOption returns emitter option
func NewEmitterOption(config *ini.File, app string) (*Option, error) {
	opt := new(Option)
	opt.App = app
	itv, err := time.ParseDuration(config.Section(fmt.Sprintf("%s_%s", app, name)).Key("interval").String())
	if err != nil {
		return nil, errors.Trace(err)
	}
	opt.Interval = itv
	return opt, nil
}

// Event ...
func (e *Emitter) Event(q *chan *message.Message) {
	t := time.NewTicker(e.Interval)
	for {
		select {
		case <-t.C:
			logrus.Debug("emitter.Event")
			*q <- message.NewMessage(true, "", nil)
		}
	}
}

// Name returns plugin name
func (e *Emitter) Name() string {
	return e.name
}

// App returns app name
func (e *Emitter) App() string {
	return e.app
}
