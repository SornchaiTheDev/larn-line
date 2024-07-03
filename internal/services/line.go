package services

import (
	"errors"
	"larn-line/internal/constants"
	"larn-line/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LineService struct {
	bot           *messaging_api.MessagingApiAPI
	channelSecret string
	channelToken  string
}

func NewLineService(channelSecret string, channelToken string) (*LineService, error) {

	bot, err := messaging_api.NewMessagingApiAPI(
		channelToken,
	)

	if err != nil {
		return nil, err
	}

	return &LineService{
		bot,
		channelSecret,
		channelToken,
	}, nil
}

var users []string

func (app *LineService) Callback(c *gin.Context) {

	cb, err := webhook.ParseRequest(app.channelSecret, c.Request)

	if err != nil {
		log.Printf("Cannot parse request: %+v\n", err)
		if errors.Is(err, webhook.ErrInvalidSignature) {
			c.Status(400)
		} else {
			c.Status(500)
		}
		return
	}

	log.Println("Handling events...")

	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				switch message.Text {
				case constants.NEWS_CHECK:
					app.sendMessages(e.ReplyToken, []messaging_api.MessageInterface{
						&messaging_api.TextMessage{
							Text: constants.NEWS_CHECK_MESSAGE,
						},
						&messaging_api.VideoMessage{
							OriginalContentUrl: "https://storage.googleapis.com/smooth-brain-bucket/ShareChat.mov",
							PreviewImageUrl:    "https://storage.googleapis.com/smooth-brain-bucket/Untitled%20design.png",
						},
					})
				default:
					app.handleLarnMessage(message.Text, e.ReplyToken)
				}
				switch s := e.Source.(type) {
				case webhook.UserSource:
					if !utils.Has(users, s.UserId) {
						users = append(users, s.UserId)
					}
				}

			default:
				log.Printf("Unsupported message content: %T\n", e.Message)
			}
		case webhook.FollowEvent:
			if _, err := app.bot.ReplyMessage(
				&messaging_api.ReplyMessageRequest{
					ReplyToken: e.ReplyToken,
					Messages: []messaging_api.MessageInterface{
						&messaging_api.TextMessage{
							Text:       constants.WELCOME_MESSAGE,
							QuickReply: createQuickReply([]string{"เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"}),
						},

						&messaging_api.TextMessage{
							Text:       constants.EXAMPLE_MESSAGE_1,
							QuickReply: createQuickReply([]string{"เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"}),
						},
						&messaging_api.TextMessage{
							Text:       constants.EXAMPLE_MESSAGE_2,
							QuickReply: createQuickReply([]string{"เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"}),
						},
						&messaging_api.TextMessage{
							Text:       constants.EXAMPLE_MESSAGE_3,
							QuickReply: createQuickReply([]string{"เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"}),
						},
						// createRegisterMessage(),
					},
				},
			); err != nil {
				log.Fatal(err)
			}

		default:
			log.Printf("Unsupported message: %T\n", event)
		}
	}

}

func (app *LineService) sendMessages(replyToken string, messages []messaging_api.MessageInterface) {

	if _, err := app.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages:   messages,
		},
	); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent text reply.")
	}
}

func createQuickReply(messages []string) *messaging_api.QuickReply {

	items := make([]messaging_api.QuickReplyItem, 0)

	for _, message := range messages {
		items = append(items, messaging_api.QuickReplyItem{
			Action: &messaging_api.MessageAction{
				Label: string([]rune(message)[:22]),
				Text:  message,
			},
		},
		)
	}

	return &messaging_api.QuickReply{
		Items: items,
	}
}

func createRegisterMessage() *messaging_api.FlexMessage {
	jsonString := `{
  "type": "bubble",
  "hero": {
    "type": "image",
    "size": "full",
    "aspectRatio": "16:9",
    "aspectMode": "cover",
    "url": "https://f.ptcdn.info/055/074/000/qvs14fduxLyta6ie2Uo-o.jpg"
  },
  "body": {
    "type": "box",
    "layout": "vertical",
    "contents": [
      {
        "type": "text",
        "text": "ยินดีต้อนรับ",
        "weight": "bold",
        "size": "xl"
      },
      {
        "type": "text",
        "text": "กดปุ่มลงทะเบียนด้านล่างเพื่อเริ่มต้นใช้งาน",
        "wrap": true
      }
    ]
  },
  "footer": {
    "type": "box",
    "layout": "vertical",
    "spacing": "sm",
    "contents": [
      {
        "type": "button",
        "style": "primary",
        "height": "sm",
        "action": {
          "type": "uri",
          "label": "ลงทะเบียน",
          "uri": "https://line.me/"
        }
      }
    ],
    "flex": 0
  }
}`
	contents, err := messaging_api.UnmarshalFlexContainer([]byte(jsonString))
	if err != nil {
		log.Fatal(err)
	}

	return &messaging_api.FlexMessage{
		AltText:    "ลงทะเบียนเพื่อเริ่มต้นใช้งาน หลานเอง",
		Contents:   contents,
		QuickReply: createQuickReply([]string{"hello"}),
	}
}

func (app *LineService) handleLarnMessage(text string, replyToken string) error {
	res, err := GetLarn(text)

	if err != nil {
		log.Fatal(err)
	}

	if _, err = app.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text:       res.Response,
					QuickReply: createQuickReply([]string{"เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"}),
				},
			},
		},
	); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent text reply.")
	}
	return nil
}

// func handleQuickReply(replyToken string) error {
// 	msg := &messaging_api.TextMessage{
// 		Text: "Select your favorite food category or send me your location!",
// 		QuickReply: &messaging_api.QuickReply{
// 			Items: []messaging_api.QuickReplyItem{
// 				{
// 					Action: &messaging_api.MessageAction{
// 						Label: "Sushi",
// 						Text:  "Sushi",
// 					},
// 				},
// 				{
// 					Action: &messaging_api.MessageAction{
// 						Label: "Tempura",
// 						Text:  "Tempura",
// 					},
// 				},
// 				{
// 					Action: &messaging_api.LocationAction{
// 						Label: "Send location",
// 					},
// 				},
// 				{
// 					Action: &messaging_api.UriAction{
// 						Label: "LINE Developer",
// 						Uri:   "https://developers.line.biz/",
// 					},
// 				},
// 			},
// 		},
// 	}
// 	if _, err := app.bot.ReplyMessage(
// 		&messaging_api.ReplyMessageRequest{
// 			ReplyToken: replyToken,
// 			Messages:   []messaging_api.MessageInterface{msg},
// 		},
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }
