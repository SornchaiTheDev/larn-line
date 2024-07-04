package utils

import "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

func CreateQuickReply(messages []string) *messaging_api.QuickReply {

	items := make([]messaging_api.QuickReplyItem, 0)

	for _, message := range messages {
		items = append(items, messaging_api.QuickReplyItem{
			Action: &messaging_api.MessageAction{
				Label: string([]rune(message)[:20]),
				Text:  message,
			},
		},
		)
	}

	return &messaging_api.QuickReply{
		Items: items,
	}
}
