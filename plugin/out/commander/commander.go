package commander

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/message"
	"github.com/go-ini/ini"
	"github.com/juju/errors"
)

var (
	name = "commander"
)

// Commander ...
type Commander struct {
	name string
	app  string
	Root string
	Cmd  string
}

// Option for Commander
type Option struct {
	App  string
	Root string
	Cmd  string
}

// NewCommander returns commander
func NewCommander(opt *Option) (*Commander, error) {
	c := new(Commander)
	c.name = name
	c.app = opt.App
	c.Root = opt.Root
	c.Cmd = opt.Cmd
	return c, nil
}

// NewOption returns commander option
func NewOption(config *ini.File, app string) (*Option, error) {
	opt := new(Option)
	opt.App = app
	sec := config.Section(fmt.Sprintf("%s_%s", app, name))
	opt.Root = sec.Key("root").String()
	opt.Cmd = sec.Key("cmd").String()
	return opt, nil
}

// Execute ...
func (c *Commander) Execute(msg *message.Message) error {
	err := os.Chdir(c.Root)
	if err != nil {
		return errors.Trace(err)
	}
	if !msg.OK {
		return errors.Trace(msg.Err)
	}
	wd, _ := os.Getwd()
	logrus.Debugf("[%s, %s] wd: %s", c.App(), c.Name(), wd)
	cmd := exec.Command("/bin/sh", c.Cmd, msg.Content)
	logrus.Debugf("[%s, %s] cmd: %v", c.App(), c.Name(), cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Trace(err)
	}
	logrus.Infof("[%s] commander.Execute, root: %s, cmd: %v, result: %s", c.App(), c.Root, cmd.Args, string(out))
	return nil
}

// Name returns plugin name
func (c *Commander) Name() string {
	return c.name
}

// App returns app name
func (c *Commander) App() string {
	return c.app
}
