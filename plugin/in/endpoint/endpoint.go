package endpoint

import (
	"fmt"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/message"
	"github.com/go-ini/ini"
	"github.com/juju/errors"
)

var (
	name = "endpoint"
)

// Endpoint accepts requests
type Endpoint struct {
	name  string
	app   string
	port  string
	route string
	reqCh chan reqMsg
}

// Option for Endpoint
type Option struct {
	App   string
	Port  string
	Route string
}

type reqMsg struct {
	content string
}

// NewEndpoint returns emitter
func NewEndpoint(opt *Option) (*Endpoint, error) {
	ept := new(Endpoint)
	ept.name = name
	ept.app = opt.App
	ept.port = opt.Port
	ept.route = opt.Route
	ept.reqCh = make(chan reqMsg, 100)
	go func() {
		if err := ept.runServer(); err != nil {
			logrus.Fatal(err)
		}
	}()
	return ept, nil
}

// NewOption returns endpoint option
func NewOption(config *ini.File, app string) (*Option, error) {
	opt := new(Option)
	opt.App = app
	opt.Port = config.Section(fmt.Sprintf("%s_%s", app, name)).Key("port").String()
	opt.Route = config.Section(fmt.Sprintf("%s_%s", app, name)).Key("route").String()
	return opt, nil
}

func (e *Endpoint) runServer() error {
	h := http.NewServeMux()
	h.HandleFunc(e.route, func(w http.ResponseWriter, r *http.Request) {
		e.reqCh <- reqMsg{
			content: fmt.Sprintf("route: %s", r.URL.Path),
		}
	})
	return errors.Trace(http.ListenAndServe(":"+e.port, h))
}

// Event ...
func (e *Endpoint) Event(q *chan *message.Message) {
	for {
		select {
		case msg := <-e.reqCh:
			logrus.Debug("endpoint.Event")
			*q <- message.NewMessage(true, msg.content, nil)
		}
	}
}

// Name returns plugin name
func (e *Endpoint) Name() string {
	return e.name
}

// App returns app name
func (e *Endpoint) App() string {
	return e.app
}
