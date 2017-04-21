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
	Root string
	Cmd  string
}

// Option for Commander
type Option struct {
	Root string
	Cmd  string
}

// NewCommander returns commander
func NewCommander(opt *Option) (*Commander, error) {
	c := new(Commander)
	c.name = name
	c.Root = opt.Root
	c.Cmd = opt.Cmd
	return c, nil
}

// NewOption returns commander option
func NewOption(config *ini.File, app string) (*Option, error) {
	opt := new(Option)
	opt.Root = config.Section(fmt.Sprintf("%s_%s", app, name)).Key("root").String()
	opt.Cmd = config.Section(fmt.Sprintf("%s_%s", app, name)).Key("cmd").String()
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
	logrus.Debug(wd, c.Cmd)
	cmd := exec.Command("/bin/sh", "-c", c.Cmd, msg.Content)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Trace(err)
	}
	logrus.Infof("commander.Execute, root: %s, cmd: %s, result: %s", c.Root, c.Cmd, string(out))
	return nil
}

// Name returns plugin name
func (c *Commander) Name() string {
	return c.name
}
