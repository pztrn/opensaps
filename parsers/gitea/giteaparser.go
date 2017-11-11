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

type GiteaParser struct{}

func (gp GiteaParser) Initialize() {
	c.Log.Infoln("Initializing Gitea parser...")
}

func (gp GiteaParser) parseBranchNew(message slackmessage.SlackMessage) map[string]string {
	data := make(map[string]string)
	data["message"] = "[Repo: {repo} | Branch: {branch}] {header_message}."

	// Parse header.
	// [0] is repo, [1] is branch.
	header_data := gp.cutLinks(message.Text)
	data["repo"] = header_data[0][1]
	data["repo_url"] = header_data[0][0]
	data["branch"] = header_data[1][1]
	data["branch_url"] = header_data[1][0]

	header_msg := strings.Split(message.Text, "] ")[1]
	data["header_message"] = header_msg

	return data
}

func (gp GiteaParser) cutLinks(data string) [][]string {
	c.Log.Debugln("Passed:", data)

	r := regexp.MustCompile("((https??://[a-zA-Z0-9.#!*/ _-]+)\\|([a-zA-Z0-9.#!*/ _+-]+))")

	found := r.FindAllStringSubmatch(data, -1)

	// [i][0] - link
	// [i][1] - string for link
	var result [][]string
	for i := range found {
		res := make([]string, 0, 2)
		res = append(res, found[i][2])
		res = append(res, found[i][3])
		result = append(result, res)
	}

	c.Log.Debugln("Links cutted:", result)
	return result
}

func (gp GiteaParser) parseCommitNew(message slackmessage.SlackMessage) map[string]string {
	data := make(map[string]string)
	data["message"] = "[Repo: {repo} | Branch: {branch}] {header_message}{newline}{repeatables}"

	// Parse header.
	// [0] is repo, [1] is branch.
	header_data := gp.cutLinks(message.Text)
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
		attachment_link := gp.cutLinks(message.Attachments[i].Text)
		data["repeatable_item_commit"+strconv.Itoa(idx)] = attachment_link[0][1]
		data["repeatable_item_commit"+strconv.Itoa(idx)+"_url"] = attachment_link[0][0]
		data["repeatable_item_message"+strconv.Itoa(idx)] = strings.Split(message.Attachments[i].Text, ">: ")[1]

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
	} else if strings.Contains(message.Text, "Pull request closed") {
		data = gp.parsePullRequestClosed(message)
	} else if strings.Contains(message.Text, "Pull request edited") {
		data = gp.parsePullRequestEdited(message)
	} else if strings.Contains(message.Text, "Pull request labels updated") {
		data = gp.parsePullRequestLabelsUpdated(message)
	} else if strings.Contains(message.Text, "Pull request merged") {
		data = gp.parsePullRequestMerged(message)
	} else if strings.Contains(message.Text, "Pull request re-opened") {
		data = gp.parsePullRequestReopened(message)
	} else if strings.Contains(message.Text, "Pull request submitted") {
		data = gp.parsePullRequestSubmitted(message)
	} else if strings.Contains(message.Text, "branch created by") {
		data = gp.parseBranchNew(message)
	} else {
		return map[string]string{"message": "Unknown message type:<br />" + fmt.Sprintf("%+v", message)}
	}

	c.Log.Debugln("Message:", fmt.Sprintf("%+x", data))

	return data
}

func (gp GiteaParser) parsePullRequestClosed(message slackmessage.SlackMessage) map[string]string {
	data := make(map[string]string)
	data["message"] = "[Repo: {repo}] Pull request {pr} by {user} was closed."

	// Parse header.
	// [0] is repo, [1] is branch.
	header_data := gp.cutLinks(message.Text)
	data["pr"] = header_data[0][1]
	data["pr_url"] = header_data[0][0]
	data["user"] = header_data[1][1]
	data["user_url"] = header_data[1][0]

	repo := strings.Split(message.Text, "Pull request")[0]
	repo = strings.TrimLeft(repo, "[")
	repo = strings.TrimRight(repo, "] ")
	data["repo"] = repo

	return data
}

func (gp GiteaParser) parsePullRequestEdited(message slackmessage.SlackMessage) map[string]string {
	data := make(map[string]string)
	data["message"] = "[Repo: {repo}] Pull request {pr} was edited by {user}."

	// Parse header.
	// [0] is repo, [1] is branch.
	header_data := gp.cutLinks(message.Text)
	data["pr"] = header_data[0][1]
	data["pr_url"] = header_data[0][0]
	data["user"] = header_data[1][1]
	data["user_url"] = header_data[1][0]

	repo := strings.Split(message.Text, "Pull request")[0]
	repo = strings.TrimLeft(repo, "[")
	repo = strings.TrimRight(repo, "] ")
	data["repo"] = repo

	return data
}

func (gp GiteaParser) parsePullRequestLabelsUpdated(message slackmessage.SlackMessage) map[string]string {
	data := make(map[string]string)
	data["message"] = "[Repo: {repo}] Pull request {pr} by {user} labels was updated."

	// Parse header.
	// [0] is repo, [1] is branch.
	header_data := gp.cutLinks(message.Text)
	data["pr"] = header_data[0][1]
	data["pr_url"] = header_data[0][0]
	data["user"] = header_data[1][1]
	data["user_url"] = header_data[1][0]

	repo := strings.Split(message.Text, "Pull request")[0]
	repo = strings.TrimLeft(repo, "[")
	repo = strings.TrimRight(repo, "] ")
	data["repo"] = repo

	return data
}

func (gp GiteaParser) parsePullRequestMerged(message slackmessage.SlackMessage) map[string]string {
	data := make(map[string]string)
	data["message"] = "[Repo: {repo}] Pull request {pr} by {user} was merged."

	// Parse header.
	// [0] is repo, [1] is branch.
	header_data := gp.cutLinks(message.Text)
	data["pr"] = header_data[0][1]
	data["pr_url"] = header_data[0][0]
	data["user"] = header_data[1][1]
	data["user_url"] = header_data[1][0]

	repo := strings.Split(message.Text, "Pull request")[0]
	repo = strings.TrimLeft(repo, "[")
	repo = strings.TrimRight(repo, "] ")
	data["repo"] = repo

	return data
}

func (gp GiteaParser) parsePullRequestReopened(message slackmessage.SlackMessage) map[string]string {
	data := make(map[string]string)
	data["message"] = "[Repo: {repo}] Pull request {pr} by {user} was re-opened."

	// Parse header.
	// [0] is repo, [1] is branch.
	header_data := gp.cutLinks(message.Text)
	data["pr"] = header_data[0][1]
	data["pr_url"] = header_data[0][0]
	data["user"] = header_data[1][1]
	data["user_url"] = header_data[1][0]

	repo := strings.Split(message.Text, "Pull request")[0]
	repo = strings.TrimLeft(repo, "[")
	repo = strings.TrimRight(repo, "] ")
	data["repo"] = repo

	return data
}

func (gp GiteaParser) parsePullRequestSubmitted(message slackmessage.SlackMessage) map[string]string {
	data := make(map[string]string)
	data["message"] = "[Repo: {repo}] Pull request {pr} submitted by {user}."

	header_data := gp.cutLinks(message.Text)
	data["user"] = header_data[0][1]
	data["user_url"] = header_data[0][0]

	repo := strings.Split(message.Text, "Pull request")[0]
	repo = strings.TrimLeft(repo, "[")
	repo = strings.TrimRight(repo, "] ")
	data["repo"] = repo

	// There is a pull request description in attachments.
	prdata := gp.cutLinks(message.Attachments[0].Title)
	data["pr"] = prdata[0][1]
	data["pr_url"] = prdata[0][0]

	return data
}
