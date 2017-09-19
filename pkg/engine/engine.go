package engine

import (
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/message"
	"github.com/cosmtrek/loop/plugin/in/emitter"
	"github.com/cosmtrek/loop/plugin/in/endpoint"
	"github.com/cosmtrek/loop/plugin/in/fswatcher"
	"github.com/cosmtrek/loop/plugin/in/monitor"
	"github.com/cosmtrek/loop/plugin/out/commander"
	"github.com/cosmtrek/loop/plugin/out/echoer"
	"github.com/cosmtrek/loop/plugin/out/emailer"
	"github.com/go-ini/ini"
	"github.com/juju/errors"
)

var (
	emitterPlug   = "emitter"
	fswatcherPlug = "fswatcher"
	monitorPlug   = "monitor"
	endpointPlug  = "endpoint"
)

var (
	echoerPlug    = "echoer"
	commanderPlug = "commander"
	emailerPlug   = "emailer"
)

// Plug implements common plugin methods
type Plug interface {
	Name() string
	App() string
}

// In implements in plugin
type In interface {
	Event(q *chan *message.Message)
	Plug
}

// Out implements out plugin
type Out interface {
	Execute(msg *message.Message) error
	Plug
}

// Engine rules the world
type Engine struct {
	registerLock sync.RWMutex
	Pipes        map[string]*Pipe
}

// NewEngine returns engine
func NewEngine() *Engine {
	g := new(Engine)
	g.Pipes = make(map[string]*Pipe)
	return g
}

// RegisterPipe pushes pipe into the engine
func (g *Engine) RegisterPipe(pipe *Pipe) error {
	g.registerLock.Lock()
	g.Pipes[pipe.Name] = pipe
	g.registerLock.Unlock()
	return nil
}

// Run ...
func (g *Engine) Run() {
	var wg sync.WaitGroup
	for _, pipe := range g.Pipes {
		wg.Add(1)
		go pipe.Run(&wg)
	}
	wg.Wait() // wait forever
}

// InPlugin returns the in plugin
func InPlugin(config *ini.File, app string, in string) (In, error) {
	switch in {
	case emitterPlug:
		opt, err := emitter.NewEmitterOption(config, app)
		if err != nil {
			return nil, errors.Trace(err)
		}
		emitter, err := emitter.NewEmitter(opt)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return emitter, nil
	case fswatcherPlug:
		opt, err := fswatcher.NewOption(config, app)
		if err != nil {
			return nil, errors.Trace(err)
		}
		fswatcher, err := fswatcher.NewFswatcher(opt)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return fswatcher, nil
	case monitorPlug:
		opt, err := monitor.NewOption(config, app)
		if err != nil {
			return nil, errors.Trace(err)
		}
		monitor, err := monitor.NewMonitor(opt)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return monitor, nil
	case endpointPlug:
		opt, err := endpoint.NewOption(config, app)
		if err != nil {
			return nil, errors.Trace(err)
		}
		endpoint, err := endpoint.NewEndpoint(opt)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return endpoint, nil
	default:
		return nil, errors.New("not known in plugin")
	}
}

// OutPlugin returns the out plugin
func OutPlugin(config *ini.File, app string, out string) (Out, error) {
	switch out {
	case echoerPlug:
		opt, err := echoer.NewOption(config, app)
		if err != nil {
			return nil, errors.Trace(err)
		}
		echoer, err := echoer.NewEchoer(opt)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return echoer, nil
	case commanderPlug:
		opt, err := commander.NewOption(config, app)
		if err != nil {
			return nil, errors.Trace(err)
		}
		commander, err := commander.NewCommander(opt)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return commander, nil
	case emailerPlug:
		opt, err := emailer.NewOption(config, app)
		if err != nil {
			return nil, errors.Trace(err)
		}
		emailer, err := emailer.NewEmailer(opt)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return emailer, nil
	default:
		return nil, errors.New("not known out plugin")
	}
}

// Pipe is the flow
type Pipe struct {
	In     In
	Out    Out
	Name   string
	Enable bool
	msgQ   chan *message.Message
}

// NewPipe returns pipe
func NewPipe(config *ini.File, app string, inout string) (*Pipe, error) {
	if app == "" {
		return nil, errors.New("empty app name")
	}
	if inout == "" {
		return nil, errors.New("empty pipe")
	}
	inouts := strings.Split(inout, ",")
	in, err := InPlugin(config, app, strings.TrimSpace(inouts[0]))
	if err != nil {
		return nil, errors.Trace(err)
	}
	out, err := OutPlugin(config, app, strings.TrimSpace(inouts[1]))
	if err != nil {
		return nil, errors.Trace(err)
	}
	pipe := new(Pipe)
	pipe.In = in
	pipe.Out = out
	pipe.Name = app
	pipe.msgQ = make(chan *message.Message, 10)
	enable := config.Section(app).Key("enable").String()
	if enable == "true" {
		pipe.Enable = true
	}
	return pipe, nil
}

// Run ...
func (p *Pipe) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	if p == nil {
		return
	}
	if !p.Enable {
		logrus.Infof("[%s] OFF", p.Name)
		return
	}
	logrus.Infof("[%s] ON", p.Name)

	go p.In.Event(&p.msgQ)
	for {
		select {
		case msg := <-p.msgQ:
			if err := p.Out.Execute(msg); err != nil {
				logrus.Error(errors.ErrorStack(err))
				break
			}
		}
	}
}
