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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.dev.pztrn.name/flagger"
	configurationinterface "go.dev.pztrn.name/opensaps/config/interface"
	parserinterface "go.dev.pztrn.name/opensaps/parsers/interface"
	pusherinterface "go.dev.pztrn.name/opensaps/pushers/interface"
	slackapiserverinterface "go.dev.pztrn.name/opensaps/slack/apiserverinterface"
	slackmessage "go.dev.pztrn.name/opensaps/slack/message"
)

type Context struct {
	Config         configurationinterface.ConfigurationInterface
	SlackAPIServer slackapiserverinterface.SlackAPIServerInterface
	Flagger        *flagger.Flagger
	Parsers        map[string]parserinterface.ParserInterface
	Pushers        map[string]pusherinterface.PusherInterface
	Log            zerolog.Logger
}

func (c *Context) Initialize() {
	c.Parsers = make(map[string]parserinterface.ParserInterface)
	c.Pushers = make(map[string]pusherinterface.PusherInterface)

	// nolint:exhaustruct
	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: time.RFC3339}
	output.FormatLevel = func(lvlRaw interface{}) string {
		var formattedLvl string

		if lvl, ok := lvlRaw.(string); ok {
			lvl = strings.ToUpper(lvl)
			switch lvl {
			case "DEBUG":
				formattedLvl = fmt.Sprintf("\x1b[30m%-5s\x1b[0m", lvl)
			case "ERROR":
				formattedLvl = fmt.Sprintf("\x1b[31m%-5s\x1b[0m", lvl)
			case "FATAL":
				formattedLvl = fmt.Sprintf("\x1b[35m%-5s\x1b[0m", lvl)
			case "INFO":
				formattedLvl = fmt.Sprintf("\x1b[32m%-5s\x1b[0m", lvl)
			case "PANIC":
				formattedLvl = fmt.Sprintf("\x1b[36m%-5s\x1b[0m", lvl)
			case "WARN":
				formattedLvl = fmt.Sprintf("\x1b[33m%-5s\x1b[0m", lvl)
			default:
				formattedLvl = lvl
			}
		}

		return fmt.Sprintf("| %s |", formattedLvl)
	}

	c.Log = zerolog.New(output).With().Timestamp().Logger()

	flaggerLogger := &FlaggerLogger{log: c.Log}
	c.Flagger = flagger.New("opensaps", flagger.LoggerInterface(flaggerLogger))
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

func (c *Context) SendToParser(name string, message slackmessage.SlackMessage) map[string]interface{} {
	parser, found := c.Parsers[strings.ToLower(name)]
	if !found {
		c.Log.Error().Msgf("Parser '%s' not found, will use default one!", name)

		return c.Parsers["default"].ParseMessage(message)
	}

	return parser.ParseMessage(message)
}

func (c *Context) SendToPusher(protocol string, connection string, data slackmessage.SlackMessage) {
	pusher, ok := c.Pushers[protocol]
	if !ok {
		c.Log.Error().Msgf("Pusher not found (or initialized) for protocol '%s'!", protocol)
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
