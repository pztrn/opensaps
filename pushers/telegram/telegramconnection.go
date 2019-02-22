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
package telegrampusher

import (
	// stdlib
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	// local
	"gitlab.com/pztrn/opensaps/slack/message"
)

type TelegramConnection struct {
	botId    string
	chatId   string
	connName string
}

func (tc *TelegramConnection) Initialize(connName string, botId string, chatId string) {
	tc.connName = connName
	tc.chatId = chatId
	tc.botId = botId
}

func (tc *TelegramConnection) ProcessMessage(message slackmessage.SlackMessage) {
	// Prepare message body.
	message_data := c.SendToParser(message.Username, message)

	// Get message template.
	msg_tpl := message_data["message"]
	delete(message_data, "message")

	// Repeatables.
	var repeatables []string
	repeatables_raw, repeatables_found := message_data["repeatables"]
	if repeatables_found {
		repeatables = strings.Split(repeatables_raw, ",")
		c.Log.Debugln("Repeatable keys:", repeatables, ", length:", len(repeatables))
	}

	// Process keys.
	for key, value := range message_data {
		// Do nothing for keys with "_url" appendix.
		if strings.Contains(key, "_url") {
			c.Log.Debugln("_url key found in pre-stage, skipping:", key)
			continue
		}
		// Do nothing (yet) on repeatables.
		if strings.Contains(key, "repeatable") {
			c.Log.Debugln("Key containing 'repeatable' in pre-stage, skipping:", key)
			continue
		}

		if len(repeatables) > 0 {
			if strings.Contains(key, "repeatable_item_") {
				c.Log.Debugln("Repeatable key in pre-stage, skipping:", key)
				continue
			}
		}
		c.Log.Debugln("Processing message data key:", key)

		// Check if we have an item with "_url" appendix. This means
		// that we should generate a link.
		val_url, found := message_data[key+"_url"]
		// Generate a link and put into message if key with "_url"
		// was found.
		var s string = ""
		if found {
			c.Log.Debugln("Found _url key, will create HTML link")
			s = fmt.Sprintf("<a href='%s'>%s</a>", val_url, value)
		} else {
			c.Log.Debugln("Found no _url key, will use as-is")
			s = value
		}
		msg_tpl = strings.Replace(msg_tpl, "{"+key+"}", s, -1)
	}

	// Process repeatables.
	repeatable_tpl, repeatable_found := message_data["repeatable_message"]
	if repeatable_found {
		var repeatables_string string = ""
		repeatables_count, _ := strconv.Atoi(message_data["repeatables_count"])
		idx := 0
		for {
			if idx == repeatables_count {
				c.Log.Debug("IDX goes above repeatables_count, breaking loop")
				break
			}

			var repstring string = repeatable_tpl
			for i := range repeatables {
				c.Log.Debugln("Processing repeatable variable:", repeatables[i]+strconv.Itoa(idx))
				var data string = ""
				rdata := message_data["repeatable_item_"+repeatables[i]+strconv.Itoa(idx)]
				rurl, rurl_found := message_data["repeatable_item_"+repeatables[i]+strconv.Itoa(idx)+"_url"]
				if rurl_found {
					c.Log.Debugln("Found _url key, will create HTML link")
					data = fmt.Sprintf("<a href='%s'>%s</a>", rurl, rdata)
				} else {
					c.Log.Debugln("Found no _url key, will use as-is")
					data = rdata
				}
				repstring = strings.Replace(repstring, "{"+repeatables[i]+"}", data, -1)
			}

			repeatables_string += repstring
			c.Log.Debugln("Repeatable string:", repstring)
			idx += 1
		}

		msg_tpl = strings.Replace(msg_tpl, "{repeatables}", repeatables_string, -1)
	}

	msg_tpl = strings.Replace(msg_tpl, "{newline}", "\n", -1)

	c.Log.Debugln("Crafted message:", msg_tpl)

	// Send message.
	tc.SendMessage(msg_tpl)
}

func (tc *TelegramConnection) SendMessage(message string) {
	msgdata := url.Values{}
	msgdata.Set("chat_id", tc.chatId)
	msgdata.Set("text", message)
	msgdata.Set("parse_mode", "HTML")

	client := &http.Client{}
	botUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tc.botId)
	c.Log.Debugln("Bot URL:", botUrl)
	response, _ := client.PostForm(botUrl, msgdata)
	c.Log.Debugln("Status:", response.Status)
}

func (tc *TelegramConnection) Shutdown() {
	// There is nothing we can do actually.
}
