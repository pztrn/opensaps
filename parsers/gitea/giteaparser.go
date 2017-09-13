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
package giteaparser

import (
    // stdlib
    "fmt"
    "regexp"
    "strconv"
    "strings"

    // local
    "lab.pztrn.name/pztrn/opensaps/slack/message"
)

type GiteaParser struct {}

func (gp GiteaParser) Initialize() {
    c.Log.Infoln("Initializing Gitea parser...")
}

func (gp GiteaParser) cutCommitLink(data string) [][]string {
    var links [][]string
    c.Log.Debugln("Passed:", data)

    r := regexp.MustCompile("(http[?s]://[a-zA-Z1-9./-]+)|([a-zA-Z1-9]+)>")

    found := r.FindAllStringSubmatch(data, -1)

    var result []string
    for i := range found {
        if i%2 == 0 {
            result = make([]string, 0, 2)
            result = append(result, found[i][1])
        } else {
            result = append(result, found[i][2])
            links = append(links, result)
        }
    }

    c.Log.Debugln("Links cutted:", links)
    return links
}

func (gp GiteaParser) cutHeaderLinks(data string) [][]string {
    var links [][]string
    c.Log.Debugln("Passed:", data)

    r := regexp.MustCompile("<(http[?s]://[a-zA-Z0-9./-]+)|([a-zA-Z0-9_-]+)>")

    found := r.FindAllStringSubmatch(data, -1)

    var result []string
    for i := range found {
        if i%2 == 0 {
            result = make([]string, 0, 2)
            result = append(result, found[i][1])
        } else {
            result = append(result, found[i][2])
            links = append(links, result)
        }
    }

    c.Log.Debugln("Links cutted:", links)
    return links
}

func (gp GiteaParser) parseCommitNew(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[Repo: {repo} | Branch: {branch}] {header_message}{newline}{repeatables}"

    // Parse header.
    // [0] is repo, [1] is branch.
    header_data := gp.cutHeaderLinks(message.Text)
    data["repo"] = header_data[0][1]
    data["repo_url"] = header_data[0][0]
    data["branch"] = header_data[1][1]
    data["branch_url"] = header_data[1][0]

    header_msg := strings.Split(message.Text, "] ")[1]
    data["header_message"] = header_msg

    // Parse commits.
    data["repeatable_message"] = "{commit}: {message}{newline}"
    data["repeatables"] = "commit,message"
    idx := 0
    for i := range message.Attachments {
        attachment_link := gp.cutCommitLink(message.Attachments[i].Text)
        data["repeatable_item_commit" + strconv.Itoa(idx)] = attachment_link[0][1]
        data["repeatable_item_commit" + strconv.Itoa(idx) + "_url"] = attachment_link[0][0]
        data["repeatable_item_message" + strconv.Itoa(idx)] = strings.Split(message.Attachments[i].Text, ">: ")[1]

        idx += 1
    }
    data["repeatables_count"] = strconv.Itoa(idx)

    return data
}

func (gp GiteaParser) ParseMessage(message slackmessage.SlackMessage) map[string]string {
    c.Log.Debugln("Parsing Gitea message...")

    var data map[string]string
    if strings.Contains(message.Text, "new commit") && strings.Contains(message.Text, "pushed by ") {
        data = gp.parseCommitNew(message)
    } else {
        return map[string]string{"message": "Unknown message type:<br />" + fmt.Sprintf("%+v", message)}
    }

    c.Log.Debugln("Message:", fmt.Sprintf("%+x", data))

    return data
}
