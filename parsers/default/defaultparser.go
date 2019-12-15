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
package defaultparser

import (
	// stdlib
	"regexp"
	"strings"

	// local
	slackmessage "go.dev.pztrn.name/opensaps/slack/message"
)

type DefaultParser struct{}

func (dp DefaultParser) Initialize() {
	c.Log.Infoln("Initializing default parser...")
}

func (dp DefaultParser) ParseMessage(message slackmessage.SlackMessage) map[string]interface{} {
	c.Log.Debugln("Parsing default message...")

	msg := message.Text + "\n"
	for _, attachment := range message.Attachments {
		msg += attachment.Text + "\n"
	}

	// Remove line break in very beginning, if present.
	if strings.Contains(msg[0:3], "\n") {
		c.Log.Debugln("Initial br found, removing")

		msg = strings.Replace(msg, "\n", "", 1)
	}

	// Get all links from message.
	r := regexp.MustCompile(`<{1}([\pL\pP\pN]+)\|{1}([\pL\pP\pN\pZs]+)>{1}`)
	foundLinks := r.FindAllStringSubmatch(msg, -1)
	c.Log.Debugln("Found links:", foundLinks)

	data := make(map[string]interface{})
	data["message"] = msg
	data["links"] = foundLinks

	return data
}
