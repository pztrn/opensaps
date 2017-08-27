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
package matrixpusher

import (
    // local
    "lab.pztrn.name/pztrn/opensaps/slack/message"
)

type MatrixPusher struct {}

func (mp MatrixPusher) Initialize() {
    c.Log.Infoln("Initializing Matrix protocol pusher...")

    // Get configuration for pushers and initialize every connection.
    cfg := c.Config.GetConfig()
    for name, config := range cfg.Matrix {
        c.Log.Infof("Initializing connection: '%s'", name)
        conn := MatrixConnection{}
        connections[name] = &conn
        go conn.Initialize(name, config.ApiRoot, config.User, config.Password, config.Room)
    }
}

func (mp MatrixPusher) Push(connection string, data slackmessage.SlackMessage) {
    conn, found := connections[connection]
    if !found {
        c.Log.Errorf("Connection not found: '%s'!", connection)
        return
    }
    c.Log.Debugf("Pushing data to '%s'", connection)
    conn.ProcessMessage(data)
}

func (mp MatrixPusher) Shutdown() {
    c.Log.Infoln("Shutting down Matrix pusher...")

    for _, conn := range connections {
        conn.Shutdown()
    }
}
