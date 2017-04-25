package main

import (
	"flag"
	"os"
	"runtime/debug"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cosmtrek/loop/pkg/engine"
	"github.com/go-ini/ini"
	"github.com/juju/errors"
)

var (
	debugMode bool
	cfgfile   string
	// BuildTimestamp ...
	BuildTimestamp string
	// Version ...
	Version string
)

func init() {
	flag.BoolVar(&debugMode, "D", false, "debug mode")
	flag.StringVar(&cfgfile, "C", "./loop.ini", "config file")
	flag.Parse()
}

func main() {
	logrus.Infoln(`
   __    ___  ___  ___
  / /   /___\/___\/ _ \
 / /   //  ///  // /_)/
/ /___/ \_// \_// ___/
\____/\___/\___/\/
	`)
	logrus.Infof("build timestamp: %s, version: %s", BuildTimestamp, Version)
	flag.Parse()
	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Info("DEBUG!!!")
	}
	cfg, err := ini.Load(cfgfile)
	if err != nil {
		logrus.Error(err)
	}
	eng := engine.NewEngine()
	apps := strings.Split(cfg.Section("").Key("apps").String(), ",")
	for _, app := range apps {
		name := strings.TrimSpace(app)
		pipe, err := engine.NewPipe(cfg, name, cfg.Section(name).Key("pipe").String())
		if err != nil {
			logrus.Error(errors.ErrorStack(err))
			continue
		}
		eng.RegisterPipe(pipe)
	}

	defer func() {
		if e := recover(); e != nil {
			logrus.Error(e, debug.Stack())
			os.Exit(1)
		}
	}()

	eng.Run()
}
