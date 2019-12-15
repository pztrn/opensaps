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
package main

import (
	// stdlib
	"os"
	"os/signal"
	"syscall"

	// local
	"go.dev.pztrn.name/opensaps/config"
	"go.dev.pztrn.name/opensaps/context"
	defaultparser "go.dev.pztrn.name/opensaps/parsers/default"
	matrixpusher "go.dev.pztrn.name/opensaps/pushers/matrix"
	telegrampusher "go.dev.pztrn.name/opensaps/pushers/telegram"
	"go.dev.pztrn.name/opensaps/slack"
)

func main() {
	c := context.New()
	c.Initialize()

	config.New(c)

	c.Log.Infoln("Launching OpenSAPS...")

	c.Flagger.Parse()
	c.Config.InitializeLater()
	c.Config.LoadConfigurationFromFile()

	slack.New(c)

	// Initialize parsers.
	defaultparser.New(c)

	// Initialize pushers.
	matrixpusher.New(c)
	telegrampusher.New(c)

	// CTRL+C handler.
	signalHandler := make(chan os.Signal, 1)
	shutdownDone := make(chan bool, 1)

	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalHandler
		c.Shutdown()
		shutdownDone <- true
	}()

	<-shutdownDone
	os.Exit(0)
}
