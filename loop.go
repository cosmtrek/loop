package main

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/engine"
	"github.com/go-ini/ini"
	"github.com/juju/errors"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	cfg, err := ini.Load("./loop.ini")
	if err != nil {
		logrus.Error(err)
	}
	eng := engine.NewEngine()
	apps := strings.Split(cfg.Section("").Key("apps").String(), ",")
	for _, app := range apps {
		pipe, err := engine.NewPipe(cfg, app, cfg.Section(app).Key("pipe").String())
		if err != nil {
			logrus.Error(errors.ErrorStack(err))
			continue
		}
		eng.RegisterPipe(pipe)
	}
	eng.Run()
}
