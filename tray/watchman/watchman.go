package watchman

import (
	"sync"

	"github.com/phlipse/systray"
)

// Watchman contains registry of menu items and their current state.
type Watchman struct {
	mutex       sync.RWMutex
	mainItems   []*systray.MenuItem
	deviceItems []*systray.MenuItem
	prevState   []int
	locked      bool // stateful locking flag
}

var (
	watchmanInstance *Watchman
	watchmanInit     sync.Once
)

// Get returns the watchman instance.
func Get() *Watchman {
	watchmanInit.Do(func() {
		watchmanInstance = &Watchman{}
	})

	return watchmanInstance
}

// Register adds a menu item to the registry.
func (w *Watchman) Register(m *systray.MenuItem) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if m.GetTitle() == "DEVICE" {
		w.deviceItems = append(w.deviceItems, m)
	} else {
		w.mainItems = append(w.mainItems, m)
	}
}

// GetDevices returns a list of device menu items.
func (w *Watchman) GetDevices() []*systray.MenuItem {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.deviceItems
}

// Lock locks the registered menu items.
func (w *Watchman) Lock() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.locked = true

	w.prevState = nil // wipe old state

	t := append(w.mainItems, w.deviceItems...)

	for idx := range t {
		if !t[idx].Disabled() {
			w.prevState = append(w.prevState, idx)
			t[idx].Disable()
		}
	}
}

// Unlock unlocks the registered menu items.
func (w *Watchman) Unlock() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	t := append(w.mainItems, w.deviceItems...)

	for idx := range w.prevState {
		// we need a dedicated check for empty device slots
		// reload could reveal an absent device that was cached in prevState
		if t[idx].GetTitle() != "DEVICE" {
			t[idx].Enable()
		}
	}

	w.locked = false
}

// Locked returns if the menu items are currently locked.
func (w *Watchman) Locked() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.locked
}
