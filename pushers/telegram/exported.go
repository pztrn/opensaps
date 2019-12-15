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
	// local
	"go.dev.pztrn.name/opensaps/context"
	pusherinterface "go.dev.pztrn.name/opensaps/pushers/interface"
)

var (
	c           *context.Context
	connections map[string]*TelegramConnection
)

func New(cc *context.Context) {
	c = cc
	connections = make(map[string]*TelegramConnection)

	tp := TelegramPusher{}
	c.RegisterPusherInterface("telegram", pusherinterface.PusherInterface(tp))
}
