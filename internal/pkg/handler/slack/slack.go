package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
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
	clusterName = kingpin.Flag(
		"cluster-name",
		"Name of k8s cluster where kwatchman is running, providing context into the notification").Default(
		"Undefined cluster").Envar("KW_CLUSTERNAME").Short('c').String()
	webhookURL = kingpin.Flag(
		"slack-webhook",
		"The slack webhook url (Required)").Envar("KW_SLACK_WEBHOOK").Short('w').Required().String()
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

type cliArgs struct {
	clusterName string
	webhookURL  string
}

func newCLI() *cliArgs {
	return &cliArgs{
		clusterName: *clusterName,
		webhookURL:  *webhookURL,
	}
}

type slackHandler struct {
	opts *cliArgs
}

func NewSlackHandler() handler.Handler {
	return &slackHandler{
		opts: newCLI(),
	}
}

func (h *slackHandler) Run(ctx context.Context, input handler.Input) (handler.Output, error) {
	title := fmt.Sprintf("%s %s\n%s", strings.ToUpper(input.Evt.Kind), input.ResourceKind, input.Evt.Key)

	// From Aug-2018 Slack requires text field to be under 4000 characters
	// https://api.slack.com/changelog/2018-04-truncating-really-long-messages
	var text string
	if len(input.Payload) == 0 {
		text = ""
	} else {
		text = fmt.Sprintf("```%s```", truncateString(string(input.Payload), 3994))
	}

	// https://api.slack.com/docs/message-attachments
	attachment := slack.Attachment{
		Title:      title,
		Color:      EventColour[input.Evt.Kind],
		Fallback:   title,
		AuthorName: "snebel29/kwatchman",
		AuthorLink: "https://github.com/snebel29/kwatchman",
		Text:       text,
		Footer:     h.opts.clusterName,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := &slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	output := handler.Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     true,
	}

	err := slack.PostWebhook(h.opts.webhookURL, msg)
	if err != nil {
		err = errors.Wrap(err, "PostWebhook: ")
		output.RunNext = false
		return output, err
	}
	return output, nil
}
