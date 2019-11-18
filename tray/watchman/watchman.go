package watchman

import (
	"sync"

	"github.com/phlipse/systray"
)

type watchman struct {
	mutex       sync.RWMutex
	mainItems   []*systray.MenuItem
	deviceItems []*systray.MenuItem
	prevState   []int
	locked		bool	// stateful locking flag
}

var (
	watchmanInstance *watchman
	watchmanInit     sync.Once
)

func Get() *watchman {
	watchmanInit.Do(func() {
		watchmanInstance = &watchman{}
	})

	return watchmanInstance
}

func (w *watchman) Register(m *systray.MenuItem) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if m.GetTitle() == "DEVICE" {
		w.deviceItems = append(w.deviceItems, m)
	} else {
		w.mainItems = append(w.mainItems, m)
	}
}

func (w *watchman) GetDevices() []*systray.MenuItem {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.deviceItems
}

func (w *watchman) Lock() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.locked = true

	w.prevState = nil // wipe old state

	t := append(w.mainItems, w.deviceItems...)

	for idx, _ := range t {
		if !t[idx].Disabled() {
			w.prevState = append(w.prevState, idx)
			t[idx].Disable()
		}
	}
}

func (w *watchman) Unlock() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	t := append(w.mainItems, w.deviceItems...)

	for idx, _ := range w.prevState {
		// we need a dedicated check for empty device slots
		// reload could reveal an absent device that was cached in prevState
		if t[idx].GetTitle() != "DEVICE" {
			t[idx].Enable()
		}
	}

	w.locked = false
}

func (w *watchman) Locked() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.locked
}
