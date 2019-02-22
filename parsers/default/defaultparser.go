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
	// local
	"gitlab.com/pztrn/opensaps/slack/message"
)

type DefaultParser struct{}

func (dp DefaultParser) Initialize() {
	c.Log.Infoln("Initializing default parser...")
}

func (dp DefaultParser) ParseMessage(message slackmessage.SlackMessage) map[string]string {
	c.Log.Debugln("Parsing default message...")

	data := make(map[string]string)
	data["message"] = message.Text
	return data
}
