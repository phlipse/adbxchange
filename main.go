package main

import (
	"fmt"
	"os"

	"github.com/phlipse/adbxchange/configuration"
	"github.com/phlipse/adbxchange/tray"
	"github.com/phlipse/go-adb"
	"github.com/phlipse/go-silo"
	"github.com/phlipse/systray"
)

// will be set during compile time
var (
	VERSION string = "unknown"
	BUILD   string = "unknown"
)

func main() {
	// we do not need flags except for version info at the moment
	// so do it in a cheap way
	if len(os.Args) > 1 && os.Args[1] == "version" {
		silo.PrintConsole(fmt.Sprintf("VERSION:\t%s\nBUILD:\t\t%s", VERSION, BUILD))
		os.Exit(0)
	}

	// initialize configuration
	c := configuration.Get()

	// initialize logging
	var f *os.File
	var err error
	if c.Logfile == "console" {
		f = os.Stdout
	} else {
		f, err = os.OpenFile(c.Logfile, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			silo.PrintConsole(fmt.Errorf("could not write to logfile: %v", err).Error())
			os.Exit(1)
		}
		defer f.Close()
	}
	l := silo.Init(f, silo.INFO)
	l.Info("adbXchange %s %s", VERSION, BUILD)

	// check if ADB exists
	if adb.ADBPath == "" {
		l.Error("ADB executable not found in PATH")
		os.Exit(1)
	}

	// run systray
	systray.Run(tray.ReadyHandler, tray.ExitHandler)
}
