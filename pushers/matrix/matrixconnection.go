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
	messageData := c.SendToParser(message.Username, message)

	messageToSend := messageData["message"].(string)
	// We'll use HTML, so reformat links accordingly (if any).
	linksRaw, linksFound := messageData["links"]
	if linksFound {
		links := linksRaw.([][]string)
		for _, link := range links {
			messageToSend = strings.Replace(messageToSend, "<"+link[0]+">", `<a href="`+link[2]+`">`+link[3]+`</a>`, -1)
		}
	}

	// "\n" should be "<br>".
	messageToSend = strings.Replace(messageToSend, "\n", "<br>", -1)

	c.Log.Debugln("Crafted message:", messageToSend)

	// Send message.
	mxc.SendMessage(messageToSend)
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
