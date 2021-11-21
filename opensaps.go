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
	"os"
	"os/signal"
	"syscall"

	"go.dev.pztrn.name/opensaps/config"
	"go.dev.pztrn.name/opensaps/context"
	defaultparser "go.dev.pztrn.name/opensaps/parsers/default"
	matrixpusher "go.dev.pztrn.name/opensaps/pushers/matrix"
	telegrampusher "go.dev.pztrn.name/opensaps/pushers/telegram"
	"go.dev.pztrn.name/opensaps/slack"
)

func main() {
	ctx := context.New()
	ctx.Initialize()

	config.New(ctx)

	ctx.Log.Info().Msg("Launching OpenSAPS...")

	ctx.Flagger.Parse()
	ctx.Config.InitializeLater()
	ctx.Config.LoadConfigurationFromFile()

	slack.New(ctx)

	// Initialize parsers.
	defaultparser.New(ctx)

	// Initialize pushers.
	matrixpusher.New(ctx)
	telegrampusher.New(ctx)

	// CTRL+C handler.
	signalHandler := make(chan os.Signal, 1)
	shutdownDone := make(chan bool, 1)

	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalHandler
		ctx.Shutdown()
		shutdownDone <- true
	}()

	<-shutdownDone
	os.Exit(0)
}
