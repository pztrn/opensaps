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

package configstruct

// ConfigStruct is a config's root.
type ConfigStruct struct {
	SlackHandler ConfigSlackHandler        `yaml:"slackhandler"`
	Webhooks     map[string]ConfigWebhook  `yaml:"webhooks"`
	Matrix       map[string]ConfigMatrix   `yaml:"matrix"`
	Telegram     map[string]ConfigTelegram `yaml:"telegram"`
}

// Slack handler configuration.
type ConfigSlackHandler struct {
	Listener ConfigSlackHandlerListener `yaml:"listener"`
}

type ConfigSlackHandlerListener struct {
	Address string `yaml:"address"`
}

// Webhook configuration.
type ConfigWebhook struct {
	Slack  ConfigWebhookSlack  `yaml:"slack"`
	Remote ConfigWebhookRemote `yaml:"remote"`
}

type ConfigWebhookSlack struct {
	Random1    string `yaml:"random1"`
	Random2    string `yaml:"random2"`
	LongRandom string `yaml:"longrandom"`
}

type ConfigWebhookRemote struct {
	Pusher string `yaml:"pusher"`
	PushTo string `yaml:"push_to"`
}

// Matrix pusher configuration.
type ConfigMatrix struct {
	ApiRoot  string `yaml:"api_root"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Room     string `yaml:"room"`
}

// ConfigTelegram is a telegram pusher configuration
type ConfigTelegram struct {
	BotID  string      `yaml:"bot_id"`
	ChatID string      `yaml:"chat_id"`
	Proxy  ConfigProxy `yaml:"proxy"`
}

// ConfigProxy represents proxy server configuration.
type ConfigProxy struct {
	// ProxyType is a proxy type. Currently ignored.
	Enabled   bool   `yaml:"enabled"`
	ProxyType string `yaml:"proxy_type"`
	Address   string `yaml:"address"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
}
