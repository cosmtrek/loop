package echoer

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/message"
	"github.com/go-ini/ini"
)

var (
	name = "echoer"
)

// Echoer prints text
type Echoer struct {
	name string
	Text string
}

// Option for echoer
type Option struct {
	Text string
}

// NewEchoer returns echoer
func NewEchoer(opt *Option) (*Echoer, error) {
	c := new(Echoer)
	c.name = name
	c.Text = opt.Text
	return c, nil
}

// NewOption returns echoer option
func NewOption(config *ini.File, app string) (*Option, error) {
	opt := new(Option)
	opt.Text = config.Section(fmt.Sprintf("%s_%s", app, name)).Key("text").String()
	return opt, nil
}

// Execute ...
func (c *Echoer) Execute(msg *message.Message) error {
	logrus.Debug("echoer.Execute")
	if !msg.OK {
		return msg.Err
	}
	if c.Text == "-" {
		logrus.Info(msg.Content)
	} else {
		logrus.Info(c.Text)
	}
	return nil
}

// Name returns plugin name
func (c *Echoer) Name() string {
	return c.name
}
