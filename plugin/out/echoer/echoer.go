package echoer

import (
	"fmt"

	"github.com/Sirupsen/logrus"
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

// NewEchoerOption returns echoer option
func NewEchoerOption(config *ini.File, app string) (*Option, error) {
	opt := new(Option)
	opt.Text = config.Section(fmt.Sprintf("%s_%s", app, name)).Key("text").String()
	return opt, nil
}

// Execute ...
func (c *Echoer) Execute() error {
	logrus.Debug("echoer.Execute")
	fmt.Println(c.Text)
	return nil
}

// Name returns plugin name
func (c *Echoer) Name() string {
	return c.name
}
