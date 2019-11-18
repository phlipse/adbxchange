package tray

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/phlipse/adbxchange/configuration"
	"github.com/phlipse/adbxchange/tray/watchman"
	"github.com/phlipse/go-adb"
	"github.com/phlipse/go-silo"
	"github.com/phlipse/systray"
	"github.com/sqweek/dialog"
)

func reloadFunction(m *systray.MenuItem) {
	// get watchman to maintain menu locks
	w := watchman.Get()
	c := configuration.Get()

	// should be simplified - later
	for {
		select {
		case <-m.ClickedCh:
			w.Lock()
			m.SetTitle("reloading...")
			systray.SetIcon(iconDeviceRefresh)

			reloadDevices(true)

			m.SetTitle("reload")
			w.Unlock()
		case <-time.After(time.Duration(c.Refresh) * time.Second):
			if !w.Locked() {
				w.Lock()
				m.SetTitle("auto re-load...")
				systray.SetIcon(iconDeviceRefresh)

				reloadDevices(false)

				m.SetTitle("reload")
				w.Unlock()
			}
		}
	}
}

func quitFunction(m *systray.MenuItem) {
	<-m.ClickedCh
	systray.Quit()
}

// tbd: should be refactored
func deviceFunction(m *systray.MenuItem) {
	// get watchman to maintain menu locks
	w := watchman.Get()
	c := configuration.Get()
	l := silo.Get()

	for {
		<-m.ClickedCh

		w.Lock()
		systray.SetIcon(iconDeviceRefresh) // will be re-set during reload

		// first check if the serial of the menu item that was clicked is still present
		// could be the case that device was disconnected and new one connected but no reload was done
		if state, _ := adb.GetDeviceState(m.GetTitle()); state == adb.DeviceUnknown {
			l.Warn("device menu seems to be corrupted, reloading devices")
			// we can't do anything useful with this error
			reloadDevices(true)
			w.Unlock()
			continue
		}

		// device individual key
		loadADBKey(m.GetTitle(), c.ADBKeysPath, c.AndroidDirectory)
		if state, err := adb.GetDeviceState(m.GetTitle()); err == nil && state == adb.DeviceOnline {
			reloadDevices(false)
			w.Unlock()
			continue
		}

		// developement keys as fallback
		if len(c.ADBDefaultKeys) > 0 {
			// bad workaround with found, should be coded in a smarter way
			found := false
			for idx, _ := range c.ADBDefaultKeys {
				if !found {
					loadADBKey(c.ADBDefaultKeys[idx], c.ADBKeysPath, c.AndroidDirectory)
					if state, err := adb.GetDeviceState(m.GetTitle()); err == nil {
						if state == adb.DeviceOnline {
							found = true
						} else if state == adb.DeviceOffline {
							// if device is offline and not unauthorized we have a not working ADB client
							// some clients appear via virtual USB port but are disabled in some way
							dialog.Message("Please check if ADB is enabled on the device!").Title("ADB seems to be disabled").Info()
						}
					}
				}
			}
			if found {
				reloadDevices(false)
				w.Unlock()
				continue
			}
		} else {
			l.Info("no ADB default keys configured")
		}

		// let the user pick one
		// in most cases we are lost
		p, err := dialog.File().SetStartDir(c.ADBKeysPath).Title(
			fmt.Sprintf("Select ADB Private Key for Serial %s", m.GetTitle())).Load()
		if err != nil {
			l.Debug("could not pick file from dialog: %v", err)
			reloadDevices(false)
			w.Unlock()
			continue
		}
		p = filepath.FromSlash(p)

		loadADBKey(p, c.ADBKeysPath, c.AndroidDirectory)
		if state, err := adb.GetDeviceState(m.GetTitle()); err == nil {
			if state == adb.DeviceOnline {
				// only try to copy ADB key to ADB Keys Folder when it worked
				if dialog.Message("Selected ADB key %s works and is not present in %s. Do you want to copy it?",
					filepath.Base(p), c.ADBKeysPath).
					Title("Copy new Key to Key Folder?").YesNo() {
					err = copyRegularFile(p, filepath.FromSlash(
						fmt.Sprintf("%s/%s_%s", c.ADBKeysPath, filepath.Base(p), m.GetTitle())))
					if err != nil {
						l.Debug("could not copy new ADB key to %s: %v", c.ADBKeysPath, err)
					}
				}
			} else if state == adb.DeviceUnauthorized {
				dialog.Message("Wrong ADB key provided. Please try again and provide a valid key.").Title("No working ADB key found").Info()
			} else if state == adb.DeviceOffline {
				// we need to check it here as well
				// could be the case that we have no default keys configured
				// and therefore we can not check it before
				dialog.Message("Please check if ADB is enabled on the device!").Title("ADB seems to be disabled").Info()
			}
		}

		// no need for restart because we do it before when we gather the state
		reloadDevices(false)
		w.Unlock()
	}
}

func prepareWorkspaceFunction(m *systray.MenuItem) {
	// get watchman to maintain menu locks
	w := watchman.Get()
	c := configuration.Get()
	l := silo.Get()

	for {
		<-m.ClickedCh

		// tbd: skip if we have more than one or no device connected

		w.Lock()
		m.SetTitle("prepairing...")
		systray.SetIcon(iconDeviceUpload)

		err := prepareWorkspace(c.WorkspaceSrc)
		if err != nil {
			l.Error("workspace was not successfully prepaired")
		}

		reloadDevices(false) // needed to reset icon

		m.SetTitle("prepare workspace")
		w.Unlock()
	}
}
