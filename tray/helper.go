package tray

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/phlipse/adbxchange/tray/watchman"
	"github.com/phlipse/go-adb"
	"github.com/phlipse/go-silo"
	"github.com/phlipse/systray"
)

func loadADBKey(s, keysPath, android string) error {
	l := silo.Get()

	// build up source and destination path
	var src string

	if adb.RegexpSerial.MatchString(s) {
		// match adbkey file
		pattern := fmt.Sprintf("%s%sadbkey*_%s", keysPath, string(filepath.Separator), s)
		var f []string
		if err := filepath.Walk(keysPath, func(p string, s os.FileInfo, err error) error {
			if s.Mode().IsRegular() {
				if m, err := filepath.Match(pattern, p); m && err == nil {
					f = append(f, p)
				}
			}

			return nil
		}); err != nil {
			return err
		}

		// ensure we matched only one
		if len(f) < 1 {
			return fmt.Errorf("no ADB key found")
		} else if len(f) > 1 {
			l.Warn("%s seems to be corrupted", keysPath)
			// we will take first one and hope for the best
		}
		src = f[0]
	} else if p, err := os.Stat(filepath.FromSlash(s)); err == nil && p.Mode().IsRegular() {
		src = filepath.FromSlash(s)
	} else {
		return fmt.Errorf("no valid ADB key provided")
	}

	dst := filepath.FromSlash(
		fmt.Sprintf("%s/adbkey", android))

	err := copyRegularFile(src, dst)
	if err != nil {
		return err
	}

	return adb.RestartServer()
}

func reloadDevices(restart bool) error {
	if restart {
		if err := adb.RestartServer(); err != nil {
			return err
		}
	}

	c, err := adb.GetDevices()
	if err != nil {
		return err
	}

	// get watchman only to maintain device menu items
	// helper functions shall never maintain a menu lock on their own
	w := watchman.Get()
	m := w.GetDevices()

	cur := 0
	for idx, _ := range m {
		// remove old devices
		m[idx].Uncheck()
		m[idx].Disable()
		m[idx].SetTitle("DEVICE") // needs to be set before Hide()
		m[idx].Hide()

		// populate new ones
		if cur < len(c) {
			m[idx].SetTitle(c[cur].Serial)
			m[idx].Enable()
			m[idx].Show()

			if c[cur].State == adb.DeviceOnline {
				m[idx].Check()
			}

			cur++
		}
	}

	// set systray icon according to device state
	switch cur {
	case 0:
		systray.SetIcon(iconDeviceNone)
	case 1:
		if c[0].State == adb.DeviceOnline {
			systray.SetIcon(iconDeviceConnected)
		} else {
			systray.SetIcon(iconDeviceUnauthorized)
		}
	default:
		systray.SetIcon(iconDeviceMultiple)
	}

	return nil
}

func copyRegularFile(src, dst string) error {
	// check source file
	sf, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sf.Mode().IsRegular() {
		return fmt.Errorf("non-regular source file %s (%q)", sf.Name(), sf.Mode().String())
	}

	// check destination file
	df, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !df.Mode().IsRegular() {
			return fmt.Errorf("non-regular destination file %s (%q)", df.Name(), df.Mode().String())
		}

		// should not happen
		if os.SameFile(sf, df) {
			return nil
		}
	}

	// try to create a hard link
	// shortcut on some OS
	if err = os.Link(src, dst); err == nil {
		return nil
	}

	// copy the stuff over
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		// check for error when closing out file
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func deviceConnected(serial string, devices []*adb.Device) bool {
	for _, d := range devices {
		if serial == d.Serial {
			return true
		}
	}

	return false
}

func prepareWorkspace(src string) error {
	_, err := adb.Exec(3, "shell", "mount", "-o", "remount,exec", "/tmp")
	if err != nil {
		return err
	}
	_, err = adb.Exec(3, "shell", "mkdir", "/tmp/workspace")
	if err != nil {
		return err
	}
	_, err = adb.Exec(30, "push", src, "/tmp/workspace")
	if err != nil {
		return err
	}
	_, err = adb.Exec(3, "shell", "chmod", "+x", "/tmp/workspace/*")
	if err != nil {
		return err
	}

	return nil
}
