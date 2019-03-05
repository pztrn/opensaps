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

	// local
	"gitlab.com/pztrn/opensaps/slack/message"
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

	// Get all links from message.
	r := regexp.MustCompile("((https??://[a-zA-Z0-9.#!*/ _-]+)\\|([a-zA-Z0-9.#!*/ _+-]+))")
	foundLinks := r.FindAllStringSubmatch(msg, -1)

	// Replace them.
	/*for _, link := range foundLinks {
		c.Log.Debugln("Link:", link)
		msg = strings.Replace(msg, "<"+link[0]+">", "<a href='"+link[2]+"'>"+link[3]+"</a>", 1)
	}*/

	data := make(map[string]interface{})
	data["message"] = msg
	data["links"] = foundLinks
	return data
}
