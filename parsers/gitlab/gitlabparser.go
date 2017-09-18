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
package gitlabparser

import (
    // stdlib
    "fmt"
    "regexp"
    "strconv"
    "strings"

    // local
    "lab.pztrn.name/pztrn/opensaps/slack/message"
)

type GitlabParser struct {}

func (gp GitlabParser) Initialize() {
    c.Log.Infoln("Initializing Gitlab parser...")
}

func (gp GitlabParser) parseCommit(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{repo}] {user} pushed to {branch}, {compare_changes}. Commits:{newline}{repeatables}"
    data["user"] = strings.TrimSpace(strings.Split(message.Text, "pushed to")[0])

    // Parse links.
    links_data := gp.parseIssueCommentLink(message.Text)
    data["branch"] = links_data[0][1]
    data["branch_url"] = links_data[0][0]
    data["repo"] = links_data[1][1]
    data["repo_url"] = links_data[1][0]
    data["compare_changes"] = "compare changes"
    data["compare_changes_url"] = links_data[2][0]

    // Parse commits.
    data["repeatable_message"] = "{commit}: {commit_text}"
    data["repeatables"] = "commit,commit_text"
    idx := 0
    for i := range message.Attachments {
        commit_data := gp.parseIssueCommentLink(message.Attachments[i].Text)
        data["repeatable_item_commit" + strconv.Itoa(idx)] = commit_data[0][1]
        data["repeatable_item_commit" + strconv.Itoa(idx) + "_url"] = commit_data[0][0]
        data["repeatable_item_commit_text" + strconv.Itoa(idx)] = strings.Split(message.Attachments[i].Text, ">: ")[1]

        idx += 1
    }
    data["repeatables_count"] = strconv.Itoa(idx)

    return data
}

func (gp GitlabParser) parseCommitLinks(data string) [][]string {
    r := regexp.MustCompile("((htt[?p|ps]://[a-zA-Z0-9./-]+)\\|([a-zA-Z0-9./ _-]+))")

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

func (gp GitlabParser) parseIssueClosed(text string) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] {user} closed issue {issue}"

    // User name comes after "closed by" words.
    data["user"] = strings.Split(text, "closed by ")[1]

    // Parse links.
    // Same as for parseIssueComment because this regexp returns
    // needed data.
    links_data := gp.parseIssueCommentLink(text)
    data["project"] = links_data[0][1]
    data["project_url"] = links_data[0][0]
    data["issue"] = links_data[1][1]
    data["issue_url"] = links_data[1][0]

    return data
}

func (gp GitlabParser) parseIssueComment(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] {user} {commented_on_issue} ({issue_name}):{newline}{repeatables}"
    data["user"] = strings.TrimSpace(strings.Split(message.Text, " <")[0])

    // Parse links in main message.
    links_data := gp.parseIssueCommentLink(message.Text)
    data["commented_on_issue"] = links_data[0][1]
    data["commented_on_issue_url"] = links_data[0][0]
    data["project"] = links_data[1][1]
    data["project_url"] = links_data[1][0]
    data["issue_name"] = strings.Split(message.Text, links_data[1][1] + ">: ")[1]

    // Parse attachments, which contains comments.
    data["repeatable_message"] = "{comment}{newline}"
    data["repeatables"] = "comment"
    idx := 0
    for i := range message.Attachments {
        data["repeatable_item_comment" + strconv.Itoa(idx)] = message.Attachments[i].Text
        idx += 1
    }
    data["repeatables_count"] = strconv.Itoa(idx)

    return data
}

func (gp GitlabParser) parseIssueCommentLink(data string) [][]string {
    r := regexp.MustCompile("((htt[?p|ps]://[a-zA-Z0-9.#!*/ _-]+)\\|([a-zA-Z0-9.#!*/ _-]+))")

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

func (gp GitlabParser) parseIssueOpened(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] {user} opened an issue: {issue}{newline}{issue_text}"

    links_data := gp.parseIssueCommentLink(message.Text)
    data["project"] = links_data[0][1]
    data["project_url"] = links_data[0][0]
    data["user"] = strings.Split(message.Text, "Issue opened by ")[1]
    if len(message.Attachments) > 0 {
        data["issue"] = message.Attachments[0].Title

        // Generate valid issue URL.
        issue_number_raw := strings.Fields(message.Attachments[0].Title)[0]
        // Remove "#" and compose URL.
        issue_number := strings.Replace(issue_number_raw, "#", "", 1)
        data["issue_url"] = links_data[0][0] + "/issues/" + issue_number
        data["issue_text"] = message.Attachments[0].Text
    } else {
        // Issue was reopened.
        data["message"] = strings.Replace(data["message"], ": {issue}{newline}{issue_text}", "", 1)
        data["message"] = strings.Replace(data["message"], "opened", "reopened", 1)
    }

    return data
}

func (gp GitlabParser) parseMergeRequestClosed(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] {user} closed merge request: {merge_request}"
    data["user"] = strings.Split(message.Text, " closed <")[0]

    links_data := gp.parseIssueCommentLink(message.Text)
    data["project"] = links_data[1][1]
    data["project_url"] = links_data[1][0]
    data["merge_request"] = links_data[0][1]
    data["merge_request_url"] = links_data[0][0]

    return data
}

func (gp GitlabParser) parseMergeRequestComment(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] {user} {commented_on_merge_request} ({merge_request_name}):{newline}{repeatables}"
    data["user"] = strings.TrimSpace(strings.Split(message.Text, " <")[0])

    // Parse links in main message.
    links_data := gp.parseIssueCommentLink(message.Text)
    data["commented_on_merge_request"] = links_data[0][1]
    data["commented_on_merge_request_url"] = links_data[0][0]
    data["project"] = links_data[1][1]
    data["project_url"] = links_data[1][0]
    data["merge_request_name"] = strings.Split(message.Text, links_data[1][1] + ">: ")[1]

    // Parse attachments, which contains comments.
    data["repeatable_message"] = "{comment}{newline}"
    data["repeatables"] = "comment"
    idx := 0
    for i := range message.Attachments {
        data["repeatable_item_comment" + strconv.Itoa(idx)] = message.Attachments[i].Text
        idx += 1
    }
    data["repeatables_count"] = strconv.Itoa(idx)

    return data
}

func (gp GitlabParser) parseMergeRequestMerged(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] {user} merged {merge_request}"
    data["user"] = strings.Split(message.Text, " merged <")[0]

    links_data := gp.parseIssueCommentLink(message.Text)
    data["project"] = links_data[1][1]
    data["project_url"] = links_data[1][0]
    data["merge_request"] = links_data[0][1]
    data["merge_request_url"] = links_data[0][0]

    return data

}

func (gp GitlabParser) parseMergeRequestOpened(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] {user} opened new merge request: {merge_request}"
    data["user"] = strings.Split(message.Text, " opened <")[0]

    links_data := gp.parseIssueCommentLink(message.Text)
    data["project"] = links_data[1][1]
    data["project_url"] = links_data[1][0]
    data["merge_request"] = links_data[0][1]
    data["merge_request_url"] = links_data[0][0]

    return data
}

func (gp GitlabParser) parsePipelineMessage(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] Pipeline {pipeline_number} of branch {branch} by {user} {status} in {time}"

    var status string = ""
    if strings.Contains(message.Attachments[0].Text, "failed") {
        status = "failed"
    } else if strings.Contains(message.Attachments[0].Text, "passed") {
        status = "passed"
    }

    data["status"] = status

    user := strings.Split(message.Attachments[0].Text, "> by ")[1]
    data["user"] = strings.Split(user, " " + status + " in")[0]
    data["time"] = strings.Split(message.Attachments[0].Text, " " + status + " in ")[1]

    links_data := gp.parseCommitLinks(message.Attachments[0].Text)
    data["project"] = links_data[0][1]
    data["project_url"] = links_data[0][0]
    data["pipeline_number"] = links_data[1][1]
    data["pipeline_number_url"] = links_data[1][0]
    data["branch"] = links_data[2][1]
    data["branch_url"] = links_data[2][0]

    return data
}

func (gp GitlabParser) parsePushedNewBranch(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] {user} pushed new branch: {branch}"

    links_data := gp.parseIssueCommentLink(message.Text)
    data["branch"] = links_data[0][1]
    data["branch_url"] = links_data[0][0]
    data["project"] = links_data[1][1]
    data["project_url"] = links_data[1][0]

    data["user"] = strings.Split(message.Text, " pushed new branch")[0]

    return data
}

func (gp GitlabParser) parseTagPush(message slackmessage.SlackMessage) map[string]string {
    data := make(map[string]string)
    data["message"] = "[{project}] {user} pushed new tag: {tag}"
    data["user"] = strings.Split(message.Text, " pushed new tag")[0]

    links_data := gp.parseIssueCommentLink(message.Text)
    data["tag"] = links_data[0][1]
    data["tag_url"] = links_data[0][0]
    data["project"] = links_data[1][1]
    data["project_url"] = links_data[1][0]

    return data
}

func (gp GitlabParser) ParseMessage(message slackmessage.SlackMessage) map[string]string {
    c.Log.Debugln("Parsing Gitlab message...")

    var data map[string]string

    if strings.Contains(message.Attachments[0].Text, "Pipeline") && strings.Contains(message.Attachments[0].Text, "of branch") {
        data = gp.parsePipelineMessage(message)
    }

    if strings.Contains(message.Text, "pushed to") {
        data = gp.parseCommit(message)
    } else if strings.Contains(message.Text, "commented on issue") {
        data = gp.parseIssueComment(message)
    } else if strings.Contains(message.Text, "closed by ") {
        data = gp.parseIssueClosed(message.Text)
    } else if strings.Contains(message.Text, "Issue opened by ") {
        data = gp.parseIssueOpened(message)
    } else if strings.Contains(message.Text, "merge_requests") && strings.Contains(message.Text, " closed <") {
        data = gp.parseMergeRequestClosed(message)
    } else if strings.Contains(message.Text, "commented on merge request") {
        data = gp.parseMergeRequestComment(message)
    } else if strings.Contains(message.Text, "merge_requests") && strings.Contains(message.Text, " merged <") {
        data = gp.parseMergeRequestMerged(message)
    } else if strings.Contains(message.Text, "merge_requests") && strings.Contains(message.Text, " opened <") {
        data = gp.parseMergeRequestOpened(message)
    } else if strings.Contains(message.Text, "pushed new branch") {
        data = gp.parsePushedNewBranch(message)
    } else if strings.Contains(message.Text, " pushed new tag ") {
        data = gp.parseTagPush(message)
    } else {
        return map[string]string{"message": "Unknown message type:<br />" + fmt.Sprintf("%+v", message)}
    }

    c.Log.Debugln("Message:", fmt.Sprintf("%+v", data))

    return data
}
