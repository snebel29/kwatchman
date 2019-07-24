package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
	"strconv"
	"strings"
	"time"
)

func init() {
	registry.Register(registry.HANDLER, "slack", NewSlackHandler)
}

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

type slackHandler struct {
	config      config.Handler
	EventColour map[string]string
}

// NewSlackHandler return the slack handler
func NewSlackHandler(c config.Handler) handler.Handler {
	return &slackHandler{
		config: c,
		EventColour: map[string]string{
			"Add":    "#1ADA00",
			"Update": "#F39C12",
			"Delete": "#FF0000",
		},
	}
}

func (h *slackHandler) Run(ctx context.Context, evt *handler.Event) error {
	title := fmt.Sprintf("%s %s\n%s", strings.ToUpper(evt.K8sEvt.Kind), evt.ResourceKind, evt.K8sEvt.Key)

	// From Aug-2018 Slack requires text field to be under 4000 characters
	// https://api.slack.com/changelog/2018-04-truncating-really-long-messages
	var text string
	if len(evt.Payload) == 0 {
		text = ""
	} else {
		text = fmt.Sprintf("```%s```", truncateString(string(evt.Payload), 3994))
	}

	// https://api.slack.com/docs/message-attachments
	attachment := slack.Attachment{
		Title:      title,
		Color:      h.EventColour[evt.K8sEvt.Kind],
		Fallback:   title,
		AuthorName: "snebel29/kwatchman",
		AuthorLink: "https://github.com/snebel29/kwatchman",
		Text:       text,
		Footer:     h.config.ClusterName,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := &slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(h.config.WebhookURL, msg)
	if err != nil {
		evt.RunNext = false
		return errors.Wrap(err, "PostWebhook: ")
	}
	return nil
}
