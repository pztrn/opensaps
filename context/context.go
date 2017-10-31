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
package context

import (
	// stdlib
	"os"
	"strings"

	// local
	"lab.pztrn.name/pztrn/opensaps/config/interface"
	"lab.pztrn.name/pztrn/opensaps/parsers/interface"
	"lab.pztrn.name/pztrn/opensaps/pushers/interface"
	"lab.pztrn.name/pztrn/opensaps/slack/apiserverinterface"
	"lab.pztrn.name/pztrn/opensaps/slack/message"

	// other
	"lab.pztrn.name/golibs/flagger"
	"lab.pztrn.name/golibs/mogrus"
)

type Context struct {
	Config         configurationinterface.ConfigurationInterface
	Flagger        *flagger.Flagger
	Log            *mogrus.LoggerHandler
	Parsers        map[string]parserinterface.ParserInterface
	Pushers        map[string]pusherinterface.PusherInterface
	SlackAPIServer slackapiserverinterface.SlackAPIServerInterface
}

func (c *Context) Initialize() {
	c.Parsers = make(map[string]parserinterface.ParserInterface)
	c.Pushers = make(map[string]pusherinterface.PusherInterface)

	l := mogrus.New()
	l.Initialize()
	c.Log = l.CreateLogger("opensaps")
	c.Log.CreateOutput("stdout", os.Stdout, true)

	c.Flagger = flagger.New(c.Log)
	c.Flagger.Initialize()
}

// Registers configuration interface.
func (c *Context) RegisterConfigurationInterface(ci configurationinterface.ConfigurationInterface) {
	c.Config = ci
	c.Config.Initialize()
}

// Registers parser interface.
func (c *Context) RegisterParserInterface(name string, iface parserinterface.ParserInterface) {
	c.Parsers[name] = iface
	c.Parsers[name].Initialize()
}

// Registers Pusher interface.
func (c *Context) RegisterPusherInterface(name string, iface pusherinterface.PusherInterface) {
	c.Pushers[name] = iface
	c.Pushers[name].Initialize()
}

// Registers Slack API HTTP server control structure.
// Russians will have pretty good luff on variable name.
func (c *Context) RegisterSlackAPIServerInterface(sasi slackapiserverinterface.SlackAPIServerInterface) {
	c.SlackAPIServer = sasi
	c.SlackAPIServer.Initialize()
}

func (c *Context) SendToParser(name string, message slackmessage.SlackMessage) map[string]string {
	parser, found := c.Parsers[strings.ToLower(name)]
	if !found {
		c.Log.Errorf("Parser '%s' not found, will use default one!", name)
		return c.Parsers["default"].ParseMessage(message)
	}

	return parser.ParseMessage(message)
}

func (c *Context) SendToPusher(protocol string, connection string, data slackmessage.SlackMessage) {
	pusher, ok := c.Pushers[protocol]
	if !ok {
		c.Log.Errorf("Pusher not found (or initialized) for protocol '%s'!", protocol)
	}

	pusher.Push(connection, data)
}

// Shutdown everything.
func (c *Context) Shutdown() {
	c.SlackAPIServer.Shutdown()

	for _, pusher := range c.Pushers {
		pusher.Shutdown()
	}
}
