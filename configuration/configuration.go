package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/phlipse/go-silo"
	"gopkg.in/yaml.v2"
)

const (
	configEnv     = "ADBXCHANGE_CONFIG"
	defaultConfig = "./config.yml"
	minRefresh    = 15
)

// Configuration contains all configuration parameters from configuration file.
type Configuration struct {
	AndroidDirectory string   `yaml:"android_directory_path"`
	ADBKeysPath      string   `yaml:"adb_keys_path"`
	ADBDefaultKeys   []string `yaml:"adb_default_keys"`
	WorkspaceSrc     string   `yaml:"workspace_source_path"`
	Logfile          string   `yaml:"logfile"`
	Refresh          int      `yaml:"refresh_cycle"`
}

var (
	configInstance *Configuration
	configInit     sync.Once
)

// Get returns configuration instance.
func Get() *Configuration {
	configInit.Do(mustRead)

	return configInstance
}

func mustRead() {
	path := os.Getenv(configEnv)
	if path == "" {
		path = filepath.FromSlash(defaultConfig)
	}

	var c Configuration

	f, err := ioutil.ReadFile(path)
	if err != nil {
		silo.PrintConsole(fmt.Errorf("could not read config from file: %s", defaultConfig).Error())
		os.Exit(1)
	}

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		silo.PrintConsole(fmt.Errorf("could not parse config from file: %s", defaultConfig).Error())
		os.Exit(1)
	}

	c.AndroidDirectory = filepath.FromSlash(c.AndroidDirectory)
	// check first error is unnecessary, os.Stat() will fail
	if stat, err := os.Stat(c.AndroidDirectory); err != nil || !stat.IsDir() {
		silo.PrintConsole(fmt.Errorf("error verifying config: android directory not set correctly").Error())
		os.Exit(1)
	}

	c.ADBKeysPath = filepath.FromSlash(c.ADBKeysPath)
	if stat, err := os.Stat(c.ADBKeysPath); err != nil || !stat.IsDir() {
		silo.PrintConsole(fmt.Errorf("error verifying config: ADB keys directory not set correctly").Error())
		os.Exit(1)
	}

	if len(c.ADBDefaultKeys) > 0 {
		for idx := range c.ADBDefaultKeys {
			if stat, err := os.Stat(filepath.FromSlash(c.ADBDefaultKeys[idx])); err != nil || !stat.Mode().IsRegular() {
				// remove key
				c.ADBDefaultKeys = append(c.ADBDefaultKeys[:idx], c.ADBDefaultKeys[idx+1:]...)
			}
		}
	}

	c.WorkspaceSrc = filepath.FromSlash(c.WorkspaceSrc)
	if stat, err := os.Stat(c.WorkspaceSrc); err != nil || !stat.IsDir() {
		c.WorkspaceSrc = ""
	}

	if c.Logfile == "" || c.Logfile == "console" {
		c.Logfile = "console"
	} else {
		c.Logfile = filepath.FromSlash(c.Logfile)
	}

	if c.Refresh < minRefresh {
		c.Refresh = minRefresh
	}

	configInstance = &c
}
