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
package slack

// This is a handler for Slack Webhooks API requests.
// Slack Webhooks API generates next URL:
//      https://hooks.slack.com/services/TRANDOM1/BRANDOM2/randomlongstring
//
// Where:
//   * RANDOM1 and RANDOM2 - random 8-char thing.
//   * randomlongstring - random 24-char thing.
//
// These strings should be defined in configuration, so we can create proper
// handler for it.

import (
	// stdlib

	"context"
	"net/http"
	"time"
)

type APIServer struct{}

func (sh APIServer) Initialize() {
	c.Log.Infoln("Initializing Slack API handler...")

	// Start HTTP server.
	// As OpenSAPS designed to be behind some proxy (nginx, Caddy, etc.)
	// we will listen only to plain HTTP.
	// Note to those who wants HTTPS - proxify with nginx, Caddy, etc!
	// Don't send pull requests, patches, don't create issues! :)
	cfg := c.Config.GetConfig()

	httpsrv = &http.Server{
		Addr: cfg.SlackHandler.Listener.Address,
		// This handler will figure out from where request has come and will
		// send it to appropriate pusher. Pusher should also determine to which
		// connection data should be sent.
		Handler:        Handler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		_ = httpsrv.ListenAndServe()
	}()

	c.Log.Infof("Slack Webhooks API server starting to listen on %s", cfg.SlackHandler.Listener.Address)
}

func (sh APIServer) Shutdown() {
	c.Log.Infoln("Shutting down Slack API handler...")

	_ = httpsrv.Shutdown(context.TODO())

	c.Log.Infoln("Slack API HTTP server shutted down")
}
