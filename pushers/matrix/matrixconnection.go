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
	slackmessage "go.dev.pztrn.name/opensaps/slack/message"
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
	apiRoot string
	// Connection name.
	connName string
	// Our device ID.
	deviceID string
	// Password for user.
	password string
	// Room ID.
	roomID string
	// Token we obtained after logging in.
	token string
	// Our username for logging in to server.
	username string
}

// nolint
func (mxc *MatrixConnection) doPostRequest(endpoint string, data string) ([]byte, error) {
	c.Log.Debugln("Data to send:", data)

	apiRoot := mxc.apiRoot + endpoint
	if mxc.token != "" {
		apiRoot += fmt.Sprintf("?access_token=%s", mxc.token)
	}

	c.Log.Debugln("Request URL:", apiRoot)

	req, _ := http.NewRequest("POST", apiRoot, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to perform POST request to Matrix as '" +
			mxc.username + "' (conn " + mxc.connName + "): " + err.Error())
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.Status == "200 OK" {
		// Return body.
		return body, nil
	}

	return nil, errors.New("Status: " + resp.Status + ", body: " + string(body))
}

// nolint
func (mxc *MatrixConnection) doPutRequest(endpoint string, data string) ([]byte, error) {
	c.Log.Debugln("Data to send:", data)

	apiRoot := mxc.apiRoot + endpoint
	if mxc.token != "" {
		apiRoot += fmt.Sprintf("?access_token=%s", mxc.token)
	}

	c.Log.Debugln("Request URL:", apiRoot)

	req, _ := http.NewRequest("PUT", apiRoot, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to perform PUT request to Matrix as '" +
			mxc.username + "' (conn " + mxc.connName + "): " + err.Error())
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.Status == "200 OK" {
		// Return body.
		return body, nil
	}

	return nil, errors.New("Status: " + resp.Status + ", body: " + string(body))
}

func (mxc *MatrixConnection) generateTnxID() string {
	// Random tnxid - 16 chars.
	length := 16

	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)

	// Making sure that we have only letters and numbers and resulted
	// string will be exactly requested length.
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = mxc.generateTnxIDSecureBytes(bufferSize)
		}

		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}

	return string(result)
}

func (mxc *MatrixConnection) generateTnxIDSecureBytes(length int) []byte {
	// Get random bytes.
	var randomBytes = make([]byte, length)

	_, err := crand.Read(randomBytes)
	if err != nil {
		c.Log.Fatalln("Unable to generate random bytes for transaction ID!")
	}

	return randomBytes
}

func (mxc *MatrixConnection) Initialize(connName string, apiRoot string, user string, password string, roomID string) {
	mxc.connName = connName
	mxc.apiRoot = apiRoot
	mxc.username = user
	mxc.password = password
	mxc.roomID = roomID
	mxc.token = ""

	c.Log.Debugln("Trying to connect to", mxc.connName, "("+apiRoot+")")

	loginStr := fmt.Sprintf(`{"type": "m.login.password", "user": "%s", "password": "%s"}`, mxc.username, mxc.password)

	c.Log.Debugln("Login string:", loginStr)

	reply, err := mxc.doPostRequest("/login", loginStr)
	if err != nil {
		c.Log.Fatalf("Failed to login to Matrix with user '%s' (conn %s): '%s'", mxc.username, mxc.connName, err.Error())
	}

	// Parse received JSON and get access token.
	data := make(map[string]interface{})

	err1 := json.Unmarshal(reply, &data)
	if err1 != nil {
		c.Log.Fatalf("Failed to parse received JSON from Matrix for user '%s' (conn %s): %s (data was: %s)",
			mxc.username, mxc.connName, err1.Error(), reply)
	}

	mxc.token = data["access_token"].(string)
	mxc.deviceID = data["deviceID"].(string)

	c.Log.Debugf("Login successful for conn '%s', access token is '%s', our deviceID is '%s'",
		mxc.connName, mxc.token, mxc.deviceID)

	// We should check if we're already in room and, if not, join it.
	// We will do this by simply trying to join. We don't care about reply
	// here.
	_, err2 := mxc.doPostRequest("/rooms/"+mxc.roomID+"/join", "{}")
	if err2 != nil {
		c.Log.Fatalf("Failed to join room: %s", err2.Error())
	}
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
			messageToSend = strings.Replace(messageToSend, link[0], `<a href="`+link[1]+`">`+link[2]+`</a>`, -1)
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
	c.Log.Debugf("Sending message to connection '%s': '%s'", mxc.connName, message)

	// We should send notices as it is preferred behavior for bots and
	// appservices.
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

	reply, err := mxc.doPutRequest("/rooms/"+mxc.roomID+"/send/m.room.message/"+mxc.generateTnxID(), msgStr)
	if err != nil {
		c.Log.Fatalf("Failed to send message to room '%s' (conn: '%s'): %s", mxc.roomID, mxc.connName, err.Error())
	}

	c.Log.Debugf("Message sent, reply: %s", string(reply))
}

func (mxc *MatrixConnection) Shutdown() {
	c.Log.Infof("Shutting down connection '%s'...", mxc.connName)

	_, err := mxc.doPostRequest("/logout", "{}")
	if err != nil {
		c.Log.Errorf("Error occurred while trying to log out from Matrix (conn %s): %s", mxc.connName, err.Error())
	}

	mxc.token = ""
	c.Log.Infof("Connection '%s' successfully shutted down", mxc.connName)
}
