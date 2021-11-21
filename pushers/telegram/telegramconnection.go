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
	"fmt"
	"net/http"
	"net/url"
	"strings"

	configstruct "go.dev.pztrn.name/opensaps/config/struct"
	slackmessage "go.dev.pztrn.name/opensaps/slack/message"
)

type TelegramConnection struct {
	config configstruct.ConfigTelegram
}

func (tc *TelegramConnection) Initialize(connName string, cfg configstruct.ConfigTelegram) {
	tc.config = cfg
}

func (tc *TelegramConnection) ProcessMessage(message slackmessage.SlackMessage) {
	// Prepare message body.
	messageData := c.SendToParser(message.Username, message)

	messageToSend, _ := messageData["message"].(string)
	// We'll use HTML, so reformat links accordingly (if any).
	linksRaw, linksFound := messageData["links"]
	if linksFound {
		links, _ := linksRaw.([][]string)
		for _, link := range links {
			messageToSend = strings.ReplaceAll(messageToSend, link[0], `<a href="`+link[1]+`">`+link[2]+`</a>`)
		}
	}

	c.Log.Debug().Msgf("Crafted message: %s", messageToSend)

	// Send message.
	tc.SendMessage(messageToSend)
}

func (tc *TelegramConnection) SendMessage(message string) {
	msgdata := url.Values{}
	msgdata.Set("chat_id", tc.config.ChatID)
	msgdata.Set("text", message)
	msgdata.Set("parse_mode", "HTML")

	// Are we should use proxy?
	// nolint:exhaustivestruct
	httpTransport := &http.Transport{}

	// nolint:nestif
	if tc.config.Proxy.Enabled {
		// Compose proxy URL.
		proxyURL := "http://"
		if tc.config.Proxy.User != "" {
			proxyURL += tc.config.Proxy.User
			if tc.config.Proxy.Password != "" {
				proxyURL += ":" + tc.config.Proxy.Password
			}

			proxyURL += "@"
		}

		proxyURL += tc.config.Proxy.Address

		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			c.Log.Error().Err(err).Msg("Error while constructing/parsing proxy URL")
		} else {
			httpTransport.Proxy = http.ProxyURL(proxyURLParsed)
		}
	}

	// nolint:exhaustivestruct
	client := &http.Client{Transport: httpTransport}
	botURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tc.config.BotID)

	c.Log.Debug().Msgf("Bot URL: %s", botURL)

	// ToDo: fix it.
	// nolint
	response, err := client.PostForm(botURL, msgdata)
	if err != nil {
		c.Log.Error().Err(err).Msg("Error occurred while sending data to Telegram")
	} else {
		c.Log.Debug().Msgf("Status: %s", response.Status)
		if response.StatusCode != http.StatusOK {
			body := []byte{}
			_, _ = response.Body.Read(body)
			response.Body.Close()
			c.Log.Debug().Msg(string(body))
		}
	}
}

func (tc *TelegramConnection) Shutdown() {
	// There is nothing we can do actually.
}
