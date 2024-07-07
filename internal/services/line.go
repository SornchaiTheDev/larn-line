package services

import (
	"context"
	"errors"
	"fmt"
	"larn-line/internal/constants"
	"larn-line/internal/models"
	"larn-line/internal/utils"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LineService struct {
	bot           *messaging_api.MessagingApiAPI
	channelSecret string
	channelToken  string
	firestore     *firestore.Client
	quickReplies  *messaging_api.QuickReply
}

func NewLineService(channelSecret string, channelToken string) (*LineService, error) {

	bot, err := messaging_api.NewMessagingApiAPI(
		channelToken,
	)

	if err != nil {
		return nil, err
	}

	firestore, err := NewFirestore()
	if err != nil {
		log.Fatal(err)
	}

	quickReply := utils.CreateQuickReply([]string{"เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"})

	return &LineService{
		bot,
		channelSecret,
		channelToken,
		firestore,
		quickReply,
	}, nil
}

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
			switch s := e.Source.(type) {
			case webhook.UserSource:
				app.bot.ShowLoadingAnimation(&messaging_api.ShowLoadingAnimationRequest{
					ChatId:         s.UserId,
					LoadingSeconds: 60,
				})
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					switch message.Text {
					case constants.NEWS_CHECK:
						app.sendNewsTut(e.ReplyToken)
					case constants.READ_MORE:
						app.sendTmpMessages(s.UserId, e.ReplyToken)
					default:
						app.handleLarnMessage(s.UserId, message.Text, e.ReplyToken)
					}
				default:
					log.Printf("Unsupported message content: %T\n", e.Message)
				}
			default:
				log.Printf("Unsupported message content: %T\n", e.Message)
			}
		case webhook.FollowEvent:
			switch s := e.Source.(type) {
			case webhook.UserSource:
				app.bot.ShowLoadingAnimation(&messaging_api.ShowLoadingAnimationRequest{
					ChatId:         s.UserId,
					LoadingSeconds: 60,
				})

				app.createUserIfNotExist(s.UserId)
			}

			if _, err := app.bot.ReplyMessage(
				&messaging_api.ReplyMessageRequest{
					ReplyToken: e.ReplyToken,
					Messages: []messaging_api.MessageInterface{
						&messaging_api.TextMessage{
							Text:       constants.WELCOME_MESSAGE,
							QuickReply: app.quickReplies,
						},

						&messaging_api.TextMessage{
							Text:       constants.EXAMPLE_MESSAGE_1,
							QuickReply: app.quickReplies,
						},
						&messaging_api.TextMessage{
							Text:       constants.EXAMPLE_MESSAGE_2,
							QuickReply: app.quickReplies,
						},
						&messaging_api.TextMessage{
							Text:       constants.EXAMPLE_MESSAGE_3,
							QuickReply: app.quickReplies,
						},
					},
				},
			); err != nil {
				log.Fatal(err)
			}

		case webhook.UnfollowEvent:
			ctx := context.Background()

			switch s := e.Source.(type) {
			case webhook.UserSource:
				app.firestore.Collection("users").Doc(s.UserId).Delete(ctx)
				utils.DeleteCollection(app.firestore, fmt.Sprintf("users/%s/messages", s.UserId))
				utils.DeleteCollection(app.firestore, fmt.Sprintf("users/%s/tmp_messages", s.UserId))
			}

		default:
			log.Printf("Unsupported message: %T\n", event)
		}
	}

}

func (app *LineService) sendNewsTut(replyToken string) {
	messages := []messaging_api.MessageInterface{
		&messaging_api.TextMessage{
			Text:       constants.NEWS_CHECK_MESSAGE,
			QuickReply: app.quickReplies,
		},
		&messaging_api.VideoMessage{
			OriginalContentUrl: "https://storage.googleapis.com/smooth-brain-bucket/ShareChat.mov",
			PreviewImageUrl:    "https://storage.googleapis.com/smooth-brain-bucket/Untitled%20design.png",
			QuickReply:         app.quickReplies,
		},
	}

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

func (app *LineService) createUserIfNotExist(userId string) {
	ctx := context.Background()
	userDoc := app.firestore.Collection("users").Doc(userId)

	_, err := userDoc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			updateUserAgent(userDoc, ctx, &map[string]any{
				"currentAgent": nil,
			})
		}
	}

}

func (app *LineService) handleLarnMessage(userId string, text string, replyToken string) error {

	ctx := context.Background()

	userDoc := app.firestore.Collection("users").Doc(userId)

	histories := getUserHistory(userDoc, ctx)

	res, err := GetLarn(text, histories)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	c := make(chan *models.Message)

	go func() {

		message := <-c

		user, err := userDoc.Get(ctx)
		if err != nil {
			log.Fatal(err)
		}

		if u := user.Data(); u["currentAgent"] != message.Classification {
			updateUserAgent(userDoc, ctx,
				&map[string]interface{}{
					"currentAgent": message.Classification,
				},
			)
			utils.DeleteCollection(app.firestore, fmt.Sprintf("users/%s/messages", userId))
		}

		sendMessage(userDoc, ctx, text, "user")
		sendMessage(userDoc, ctx, message.Response, "model")

	}()

	c <- res

	splitMessages := strings.Split(res.Response, "% % % % %")

	allMessages := make([]messaging_api.MessageInterface, 0)

	recommends, err := GetRecommend(res.Response)

	quickReply := utils.CreateQuickReply(recommends)

	currentMessage := 0
	for _, message := range splitMessages {

		imgIdx := utils.IndexOf(message, '[')

		if imgIdx == -1 {
			text := strings.TrimSpace(message)
			if len(text) != 0 {
				allMessages = append(allMessages,
					messaging_api.TextMessage{
						Text:       text,
						QuickReply: quickReply,
					},
				)
				currentMessage++
			}
		} else {

			imagePart := message[imgIdx:]

			endOfImg := imgIdx + utils.IndexOf(imagePart, ']')

			image := message[imgIdx+1 : endOfImg]

			text := strings.TrimSpace(message[:imgIdx])
			if len(text) > 0 {

				allMessages = append(allMessages,
					messaging_api.TextMessage{
						Text:       text,
						QuickReply: quickReply,
					},
				)
				currentMessage++
			}

			if strings.HasPrefix(image, "https://") {
				allMessages = append(allMessages,
					messaging_api.ImageMessage{
						OriginalContentUrl: image,
						PreviewImageUrl:    image,
						QuickReply:         quickReply,
					},
				)
				currentMessage++
			}
		}

	}

	messagesLength := len(allMessages)
	var finalMessages []messaging_api.MessageInterface

	if messagesLength <= 5 {
		finalMessages = allMessages
	} else {
		tmpMessages := allMessages[5:]
		saveTmpMessage(userDoc, ctx, tmpMessages)
		quickReply = utils.CreateQuickReply([]string{"อ่านต่อ"})

		for i, message := range allMessages[:5] {
			switch m := message.(type) {
			case messaging_api.TextMessage:
				if i == 4 {
					m.QuickReply = quickReply
				}
				finalMessages = append(finalMessages, m)
			case messaging_api.ImageMessage:
				if i == 4 {
					m.QuickReply = quickReply
				}
				finalMessages = append(finalMessages, m)
			}

		}
	}

	if _, err = app.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages:   finalMessages,
		},
	); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent text reply.")
	}
	return nil
}

func getUserHistory(userDoc *firestore.DocumentRef, ctx context.Context) []models.History {

	iter := userDoc.Collection("messages").Documents(ctx)

	histories := make([]models.History, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var history models.History
		if err := doc.DataTo(&history); err != nil {
			log.Fatal(err)
		}

		histories = append(histories, history)
	}

	return histories
}

func updateUserAgent(userDoc *firestore.DocumentRef, ctx context.Context, body *map[string]any) {
	_, err := userDoc.Set(ctx, body)
	if err != nil {
		log.Fatal(err)
	}
}

func sendMessage(userDoc *firestore.DocumentRef, ctx context.Context, message string, from string) {
	_, _, err := userDoc.Collection("messages").Add(ctx, map[string]any{
		"from":      from,
		"message":   message,
		"timestamp": firestore.ServerTimestamp,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func saveTmpMessage(userDoc *firestore.DocumentRef, ctx context.Context, messages []messaging_api.MessageInterface) {
	var payload map[string]any
	for _, message := range messages {
		switch m := message.(type) {
		case messaging_api.TextMessage:
			payload = map[string]any{
				"type":      "message",
				"text":      m.Text,
				"timestamp": firestore.ServerTimestamp,
			}
		case messaging_api.ImageMessage:
			payload = map[string]any{
				"type":      "image",
				"original":  m.OriginalContentUrl,
				"preview":   m.PreviewImageUrl,
				"timestamp": firestore.ServerTimestamp,
			}
		}

		_, _, err := userDoc.Collection("tmp_messages").Add(ctx, payload)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getTmpMessages(userDoc *firestore.DocumentRef, ctx context.Context) []models.TmpHistory {

	iter := userDoc.Collection("tmp_messages").OrderBy("timestamp", firestore.Asc).Limit(5).Documents(ctx)

	histories := make([]models.TmpHistory, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var history models.TmpHistory

		m := doc.Data()

		switch m["type"].(string) {
		case "message":
			history = models.TmpHistory{
				Message: models.TextMessage{
					Text: m["text"].(string),
				},
			}
		case "image":
			history = models.TmpHistory{
				Message: models.ImageMessage{
					Preview:  m["preview"].(string),
					Original: m["original"].(string),
				},
			}
		}

		_, err = userDoc.Collection("tmp_messages").Doc(doc.Ref.ID).Delete(ctx)
		if err != nil {
			log.Println(err)
		}

		histories = append(histories, history)
	}

	return histories
}

func (app *LineService) sendTmpMessages(userId string, replyToken string) {
	userDoc := app.firestore.Collection("users").Doc(userId)
	ctx := context.Background()

	histories := getTmpMessages(userDoc, ctx)

	var allMessages []messaging_api.MessageInterface
	for _, history := range histories {
		switch h := history.Message.(type) {
		case models.TextMessage:
			allMessages = append(allMessages, &messaging_api.TextMessage{
				Text: h.GetText(),
			})
		case models.ImageMessage:
			allMessages = append(allMessages, &messaging_api.ImageMessage{
				PreviewImageUrl:    h.GetPreview(),
				OriginalContentUrl: h.GetOriginal(),
			})
		}

	}

	if len(allMessages) == 0 {
		allMessages = append(allMessages, &messaging_api.TextMessage{
			Text:       "ไม่มีข้อความให้อ่านต่อแล้วค่ะ 🤗",
			QuickReply: app.quickReplies,
		})
		return
	}

	if _, err := app.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages:   allMessages,
		},
	); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent text reply.")
	}
}
