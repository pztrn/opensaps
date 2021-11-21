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

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	slackmessage "go.dev.pztrn.name/opensaps/slack/message"
)

type Handler struct{}

func (sh Handler) ServeHTTP(respwriter http.ResponseWriter, req *http.Request) {
	c.Log.Debug().Str("method", req.Method).Str("host", req.Host).Str("path", req.URL.Path).Msg("Received HTTP request")

	// We should catch only POST requests. Otherwise return HTTP 404.
	if req.Method != "POST" {
		c.Log.Debug().Msg("Not a POST request, returning HTTP 404")
		respwriter.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(respwriter, "NOT FOUND")

		return
	}

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	c.Log.Debug().Msgf("Received body: %s", string(body))

	// Try to figure out where we should push received data.
	cfg := c.Config.GetConfig()

	var sentToPusher bool

	for name, config := range cfg.Webhooks {
		if strings.Contains(req.URL.Path, config.Slack.Random1) &&
			strings.Contains(req.URL.Path, config.Slack.Random2) &&
			strings.Contains(req.URL.Path, config.Slack.LongRandom) {
			c.Log.Debug().Msgf("Passed data belongs to '%s' and should go to '%s' pusher, protocol '%s'",
				name, config.Remote.PushTo, config.Remote.Pusher)
			// Parse message into SlackMessage structure.
			if strings.Contains(string(body)[0:7], "payload") {
				// We have HTTP form payload. It still should be a
				// parseable JSON string, we just need to do some
				// preparations.
				// First - remove "payload=" from the beginning.
				tempBody := string(body)
				tempBody = strings.Replace(tempBody, "payload=", "", 1)
				// Second - unescape data.
				tempBody, err := url.QueryUnescape(tempBody)
				if err != nil {
					c.Log.Error().Msg("Failed to decode body into parseable string!")

					return
				}

				// And finally - convert body back to bytes.
				body = []byte(tempBody)
			}

			// nolint:exhaustivestruct
			slackmsg := slackmessage.SlackMessage{}

			err := json.Unmarshal(body, &slackmsg)
			if err != nil {
				c.Log.Error().Err(err).Msg("Failed to decode JSON into SlackMessage struct")

				return
			}

			c.Log.Debug().Msgf("Received message: %+v", slackmsg)
			c.SendToPusher(config.Remote.Pusher, config.Remote.PushTo, slackmsg)

			sentToPusher = true
		}
	}

	if !sentToPusher {
		c.Log.Debug().Msg("Don't know where to push data. Ignoring with HTTP 404")
		respwriter.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(respwriter, "NOT FOUND")
	}
}
