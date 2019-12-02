package tray

import (
	"github.com/phlipse/adbxchange/configuration"
	"github.com/phlipse/adbxchange/tray/watchman"
	"github.com/phlipse/systray"
)

const MaxDevices = 5

func ReadyHandler() {
	c := configuration.Get()
	w := watchman.Get()

	// setup basic stuff
	systray.SetIcon(iconDeviceRefresh) // correct icon will be set during reload of devices
	systray.SetTitle("adbXchange")

	menuReload := systray.AddMenuItem("reload", "reload")
	go reloadFunction(menuReload)
	w.Register(menuReload)

	if c.WorkspaceSrc != "" {
		// only add menu item if it is configured
		menuPrepareWorkspace := systray.AddMenuItem("prepare workspace", "prepare workspace")
		go prepareWorkspaceFunction(menuPrepareWorkspace)
		w.Register(menuPrepareWorkspace)
	}

	systray.AddSeparator()

	// prepare menu items for devices
	for i := 0; i < MaxDevices; i++ {
		d := systray.AddMenuItem("DEVICE", "DEVICE")
		d.Disable()
		d.Hide()

		go deviceFunction(d)
		w.Register(d)
	}

	systray.AddSeparator()

	menuQuit := systray.AddMenuItem("quit", "quit")
	go quitFunction(menuQuit)
	// do not add quit to watchman

	// do one hard reload at the end
	w.Lock()
	defer w.Unlock()
	reloadDevices(true)
}

func ExitHandler() {}
