package fswatcher

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/message"
	"github.com/fsnotify/fsnotify"
	"github.com/go-ini/ini"
	"github.com/juju/errors"
)

var (
	name = "fswatcher"
)

// Event ops
const (
	Create fsnotify.Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

// Fswatcher ...
type Fswatcher struct {
	name    string
	app     string
	Dir     string
	Ops     map[string]fsnotify.Op
	Watcher *fsnotify.Watcher
}

// NewFswatcher returns fswatcher
func NewFswatcher(opt *Option) (*Fswatcher, error) {
	w := new(Fswatcher)
	w.name = name
	w.app = opt.App
	w.Dir = opt.Dir
	w.Ops = opt.Ops
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if err = watcher.Add(w.Dir); err != nil {
		return nil, errors.Trace(err)
	}
	w.Watcher = watcher
	return w, nil
}

// Event ...
func (f *Fswatcher) Event(q *chan *message.Message) {
	for {
		select {
		case event := <-f.Watcher.Events:
			logrus.Debugf("fswatcher.Event: %s", event.String())
			for _, v := range f.Ops {
				if event.Op&v == v {
					*q <- message.NewMessage(true, event.Name, nil)
					break
				}
			}
		}
	}
}

// Name returns plugin name
func (f *Fswatcher) Name() string {
	return f.name
}

// App returns app name
func (f *Fswatcher) App() string {
	return f.app
}

// Option for fswatcher
type Option struct {
	App string
	Dir string
	Ops map[string]fsnotify.Op
}

// NewOption returns option
func NewOption(config *ini.File, app string) (*Option, error) {
	opt := new(Option)
	opt.App = app
	sec := config.Section(fmt.Sprintf("%s_%s", app, name))
	opt.Dir = sec.Key("dir").String()
	events := sec.Key("events").String()
	if opt.Dir == "" {
		return nil, errors.New("empty directory for watching")
	}
	if events == "" {
		return nil, errors.New("must specify event operation")
	}
	opt.Ops = make(map[string]fsnotify.Op)
	ss := strings.Split(events, ",")
	for _, s := range ss {
		o := strings.TrimSpace(s)
		switch o {
		case "create":
			opt.Ops[o] = Create
		case "write":
			opt.Ops[o] = Write
		case "remove":
			opt.Ops[o] = Remove
		case "rename":
			opt.Ops[o] = Rename
		case "chmod":
			opt.Ops[o] = Chmod
		default:
			return nil, errors.Errorf("no op: %s", o)
		}
	}
	return opt, nil
}
