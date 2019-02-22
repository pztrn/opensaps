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
	// stdlib
	"bytes"
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	// local
	"gitlab.com/pztrn/opensaps/slack/message"
)

// Constants for random transaction ID.
const (
	letterBytes   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 36 possibilities
	letterIdxBits = 6                                      // 6 bits to represent 64 possibilities / indexes
	letterIdxMask = 1<<letterIdxBits - 1                   // All 1-bits, as many as letterIdxBits
)

type MatrixMessage struct {
	MsgType       string `json:"msgtype"`
	Body          string `json:"body"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body"`
}

type MatrixConnection struct {
	// API root for connection.
	api_root string
	// Connection name.
	conn_name string
	// Our device ID.
	device_id string
	// Password for user.
	password string
	// Room ID.
	room_id string
	// Token we obtained after logging in.
	token string
	// Our username for logging in to server.
	username string
}

func (mxc *MatrixConnection) doPostRequest(endpoint string, data string) ([]byte, error) {
	c.Log.Debugln("Data to send:", data)

	api_root := mxc.api_root + endpoint
	if mxc.token != "" {
		api_root += fmt.Sprintf("?access_token=%s", mxc.token)
	}
	c.Log.Debugln("Request URL:", api_root)
	req, _ := http.NewRequest("POST", api_root, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to perform POST request to Matrix as '" + mxc.username + "' (conn " + mxc.conn_name + "): " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.Status == "200 OK" {
		// Return body.
		return body, nil
	} else {
		return nil, errors.New("Status: " + resp.Status + ", body: " + string(body))
	}
}

func (mxc *MatrixConnection) doPutRequest(endpoint string, data string) ([]byte, error) {
	c.Log.Debugln("Data to send:", data)
	api_root := mxc.api_root + endpoint
	if mxc.token != "" {
		api_root += fmt.Sprintf("?access_token=%s", mxc.token)
	}
	c.Log.Debugln("Request URL:", api_root)
	req, _ := http.NewRequest("PUT", api_root, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to perform PUT request to Matrix as '" + mxc.username + "' (conn " + mxc.conn_name + "): " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.Status == "200 OK" {
		// Return body.
		return body, nil
	} else {
		return nil, errors.New("Status: " + resp.Status + ", body: " + string(body))
	}
}

func (mxc *MatrixConnection) generateTnxId() string {
	// Random tnxid - 16 chars.
	length := 16

	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)

	// Making sure that we have only letters and numbers and resulted
	// string will be exactly requested length.
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = mxc.generateTnxIdSecureBytes(bufferSize)
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}
	return string(result)
}

func (mxc *MatrixConnection) generateTnxIdSecureBytes(length int) []byte {
	// Get random bytes.
	var randomBytes = make([]byte, length)
	_, err := crand.Read(randomBytes)
	if err != nil {
		c.Log.Fatalln("Unable to generate random bytes for transaction ID!")
	}

	return randomBytes
}

func (mxc *MatrixConnection) Initialize(conn_name string, api_root string, user string, password string, room_id string) {
	mxc.conn_name = conn_name
	mxc.api_root = api_root
	mxc.username = user
	mxc.password = password
	mxc.room_id = room_id
	mxc.token = ""

	c.Log.Debugln("Trying to connect to", mxc.conn_name, "("+api_root+")")

	loginStr := fmt.Sprintf(`{"type": "m.login.password", "user": "%s", "password": "%s"}`, mxc.username, mxc.password)
	c.Log.Debugln("Login string:", loginStr)
	reply, err := mxc.doPostRequest("/login", loginStr)
	if err != nil {
		c.Log.Fatalf("Failed to login to Matrix with user '%s' (conn %s): '%s'", mxc.username, mxc.conn_name, err.Error())
	}
	// Parse received JSON and get access token.
	data := make(map[string]string)
	err1 := json.Unmarshal(reply, &data)
	if err1 != nil {
		c.Log.Fatalf("Failed to parse received JSON from Matrix for user '%s' (conn %s): %s", mxc.username, mxc.conn_name, err1.Error())
	}
	mxc.token = data["access_token"]
	mxc.device_id = data["device_id"]

	c.Log.Debugf("Login successful for conn '%s', access token is '%s', our device_id is '%s'", mxc.conn_name, mxc.token, mxc.device_id)

	// We should check if we're already in room and, if not, join it.
	// We will do this by simply trying to join. We don't care about reply
	// here.
	_, err2 := mxc.doPostRequest("/rooms/"+mxc.room_id+"/join", "{}")
	if err2 != nil {
		c.Log.Fatalf("Failed to join room: %s", err2.Error())
	}

	// If we're here - everything is okay and we already in room. Send
	// greeting message.
	mxc.SendMessage("OpenSAPS is back in business for connection '" + mxc.conn_name + "'!")
}

// This function launches when new data was received thru Slack API.
// It will prepare a message which will be passed to mxc.SendMessage().
func (mxc *MatrixConnection) ProcessMessage(message slackmessage.SlackMessage) {
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

	msg_tpl = strings.Replace(msg_tpl, "{newline}", "<br />", -1)

	// Replace all "\n" with "<br />".
	msg_tpl = strings.Replace(msg_tpl, "\n", "<br />", -1)

	c.Log.Debugln("Crafted message:", msg_tpl)

	// Send message.
	mxc.SendMessage(msg_tpl)
}

// This function sends already prepared message to room.
func (mxc *MatrixConnection) SendMessage(message string) {
	c.Log.Debugf("Sending message to connection '%s': '%s'", mxc.conn_name, message)

	// We should send notices as it is preferred behaviour for bots and
	// appservices.
	//msgStr := fmt.Sprintf(`{"msgtype": "m.text", "body": "%s", "format": "org.matrix.custom.html", "formatted_body": "%s"}`, message, message)
	msg := MatrixMessage{}
	msg.MsgType = "m.notice"
	msg.Body = message
	msg.Format = "org.matrix.custom.html"
	msg.FormattedBody = message

	msgBytes, err := json.Marshal(&msg)
	if err != nil {
		c.Log.Errorln("Failed to marshal message into JSON:", err.Error())
		return
	}
	msgStr := string(msgBytes)

	reply, err := mxc.doPutRequest("/rooms/"+mxc.room_id+"/send/m.room.message/"+mxc.generateTnxId(), msgStr)
	if err != nil {
		c.Log.Fatalf("Failed to send message to room '%s' (conn: '%s'): %s", mxc.room_id, mxc.conn_name, err.Error())
	}
	c.Log.Debugf("Message sent, reply: %s", string(reply))
}

func (mxc *MatrixConnection) Shutdown() {
	c.Log.Infof("Shutting down connection '%s'...", mxc.conn_name)

	_, err := mxc.doPostRequest("/logout", "{}")
	if err != nil {
		c.Log.Errorf("Error occured while trying to log out from Matrix (conn %s): %s", mxc.conn_name, err.Error())
	}
	mxc.token = ""
	c.Log.Infof("Connection '%s' successfully shutted down", mxc.conn_name)
}
