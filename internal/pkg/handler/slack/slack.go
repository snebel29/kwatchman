package slack

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
	"strings"
	"time"
)

var (
	// Event names here shoul match definition within kooper generic.go
	// Info is not a valid k8sEvent but potentially used for internal kwatchman notifications
	EventColour = map[string]string{
		"Add":    "#1ADA00",
		"Update": "#F39C12",
		"Delete": "#FF0000",
		"Info":   "#0000FF",
	}
	slackWebhookUrl = kingpin.Flag(
		"slack-webhook",
		"The slack webhook url").Envar("KW_SLACK_WEBHOOK").Short('w').Required().String()
)

func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func truncateString(text string, limit int) string {
	runes := []rune(text)
	limit = minInt(len(runes), limit)
	return string(runes[:limit])
}

func MsgToSlack(kind, objKey, clusterName, resourceKind, payload string) error {
	title := fmt.Sprintf("%s %s %s", strings.ToUpper(kind), resourceKind, objKey)

	// From Aug-2018 Slack requires text field to be under 4000 characters
	// https://api.slack.com/changelog/2018-04-truncating-really-long-messages
	if kind == "Add" || kind == "Delete" {
		payload = ""
	} else {
		payload = fmt.Sprintf("```%s```", truncateString(payload, 3994))
	}

	// https://api.slack.com/docs/message-attachments
	attachment := slack.Attachment{
		Title:      title,
		Color:      EventColour[kind],
		Fallback:   title,
		AuthorName: "snebel29/kwatchman",
		AuthorLink: "https://github.com/snebel29/kwatchman",
		Text:       payload,
		Footer:     clusterName,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(*slackWebhookUrl, &msg)
	if err != nil {
		return errors.Wrap(err, "PostWebhook: ")
	}
	return nil
}
