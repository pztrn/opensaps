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
	configurationinterface "go.dev.pztrn.name/opensaps/config/interface"
	configstruct "go.dev.pztrn.name/opensaps/config/struct"
	"go.dev.pztrn.name/opensaps/context"
)

var (
	ctx *context.Context
	// Temporary configuration.
	tempconfig map[string]string
	// Configuration from YAML file.
	config *configstruct.ConfigStruct
)

func New(cc *context.Context) {
	ctx = cc
	conf := Configuration{}
	ctx.RegisterConfigurationInterface(configurationinterface.ConfigurationInterface(conf))
}
