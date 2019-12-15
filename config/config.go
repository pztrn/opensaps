// OpenSAPS - Open Slack API server for everyone.
//
// Copyright (c) 2017, Stanislav N. aka pztrn.
// All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package config

import (
	// stdlib
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"

	// local
	configstruct "go.dev.pztrn.name/opensaps/config/struct"

	// other
	"go.dev.pztrn.name/flagger"
	"gopkg.in/yaml.v2"
)

type Configuration struct{}

// Returns configuration to caller.
func (conf Configuration) GetConfig() *configstruct.ConfigStruct {
	return config
}

// Gets value from temporary configuration storage.
// If value isn't found - returns empty string with error.
func (conf Configuration) GetTempValue(key string) (string, error) {
	value, found := tempconfig[key]
	if !found {
		return "", errors.New("No such key in temporary configuration storage: " + key)
	}

	// If we have path with tilde in front (home directory) - replace
	// tilde with actual home directory.
	if value[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			c.Log.Fatalln("Failed to get current user data: " + err.Error())
		}

		value = strings.Replace(value, "~", usr.HomeDir, 1)
	}

	return value, nil
}

func (conf Configuration) Initialize() {
	c.Log.Infoln("Initializing configuration storage...")

	tempconfig = make(map[string]string)

	flagConfigpath := flagger.Flag{
		Name:         "config",
		Description:  "Path to configuration file.",
		Type:         "string",
		DefaultValue: "~/.config/OpenSAPS/config.yaml",
	}

	_ = c.Flagger.AddFlag(&flagConfigpath)
}

// Initializes configuration root path for later usage.
func (conf Configuration) initializeConfigurationFilePath() {
	c.Log.Debugln("Asking flagger about configuration root path supplied by user...")

	configpath, err := c.Flagger.GetStringValue("config")
	if err != nil {
		c.Log.Fatalln("Something went wrong - Flagger doesn't know about \"-config\" parameter!")
	}

	c.Log.Infoln("Will use configuration file: '" + configpath + "'")
	conf.SetTempValue("CONFIGURATION_FILE", configpath)
}

// Asking Flagger about flags, initialize internal variables.
// Should be called **after** Flagger.Parse().
func (conf Configuration) InitializeLater() {
	c.Log.Infoln("Completing configuration initialization...")

	conf.initializeConfigurationFilePath()
}

// Loads configuration from file.
func (conf Configuration) LoadConfigurationFromFile() {
	configpath, err := conf.GetTempValue("CONFIGURATION_FILE")
	if err != nil {
		c.Log.Fatalln("Failed to get configuration file path from internal temporary configuration storage! OpenSAPS is BROKEN!")
	}

	c.Log.Infof("Loading configuration from '%s'...", configpath)

	// Read file into memory.
	configBytes, err1 := ioutil.ReadFile(configpath)
	if err1 != nil {
		c.Log.Fatalf("Error occurred while reading configuration file: %s", err1.Error())
	}

	config = &configstruct.ConfigStruct{}
	// Parse YAML.
	err2 := yaml.Unmarshal(configBytes, config)
	if err2 != nil {
		c.Log.Fatalf("Failed to parse configuration file: %s", err2.Error())
	}

	c.Log.Debugln("Loaded configuration:", fmt.Sprintf("%+v", config))
}

// Sets value to key in temporary configuration storage.
// If key already present in map - value will be replaced.
func (conf Configuration) SetTempValue(key, value string) {
	tempconfig[key] = value
}
