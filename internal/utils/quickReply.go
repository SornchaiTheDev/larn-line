package utils

import (
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

func CreateQuickReply(messages []string) *messaging_api.QuickReply {

	items := make([]messaging_api.QuickReplyItem, 0)

	for _, message := range messages {
		message = strings.ReplaceAll(message, "\"", "")
		runes := []rune(message)
		msgLen := len(runes)

		if msgLen > 300 {
			msgLen = 300
		}

		label := message
		if len(runes) > 17 {
			label = string(runes[:17]) + "..."
		}

		items = append(items, messaging_api.QuickReplyItem{
			Action: &messaging_api.MessageAction{
				Label: label,
				Text:  string(runes[:msgLen]),
			},
		},
		)
	}

	return &messaging_api.QuickReply{
		Items: items,
	}
}
