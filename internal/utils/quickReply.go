package utils

import (
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

func CreateQuickReply(messages []string) *messaging_api.QuickReply {

	items := make([]messaging_api.QuickReplyItem, 0)

	for _, message := range messages {
		msgLen := len([]rune(string(message)))

		if msgLen > 300 {
			msgLen = 300
		}

		message = strings.ReplaceAll(message, "\"", "")

		items = append(items, messaging_api.QuickReplyItem{
			Action: &messaging_api.MessageAction{
				Label: string([]rune(message)[:17]) + "...",
				Text:  string([]rune(message)[:msgLen]),
			},
		},
		)
	}

	return &messaging_api.QuickReply{
		Items: items,
	}
}
