package services

import (
	"context"
	"errors"
	"fmt"
	"larn-line/internal/constants"
	"larn-line/internal/models"
	"larn-line/internal/utils"
	"log"

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

	return &LineService{
		bot,
		channelSecret,
		channelToken,
		firestore,
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
					switch s := e.Source.(type) {
					case webhook.UserSource:
						app.bot.ShowLoadingAnimation(&messaging_api.ShowLoadingAnimationRequest{
							ChatId:         s.UserId,
							LoadingSeconds: 60,
						})

						if !utils.Has(users, s.UserId) {
							users = append(users, s.UserId)
						}

						app.handleLarnMessage(s.UserId, message.Text, e.ReplyToken)
					}
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

				if !utils.Has(users, s.UserId) {
					users = append(users, s.UserId)
				}

				app.createUserIfNotExist(s.UserId)
			}

			quickReply := utils.CreateQuickReply([]string{"เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"})

			if _, err := app.bot.ReplyMessage(
				&messaging_api.ReplyMessageRequest{
					ReplyToken: e.ReplyToken,
					Messages: []messaging_api.MessageInterface{
						&messaging_api.TextMessage{
							Text:       constants.WELCOME_MESSAGE,
							QuickReply: quickReply,
						},

						&messaging_api.TextMessage{
							Text:       constants.EXAMPLE_MESSAGE_1,
							QuickReply: quickReply,
						},
						&messaging_api.TextMessage{
							Text:       constants.EXAMPLE_MESSAGE_2,
							QuickReply: quickReply,
						},
						&messaging_api.TextMessage{
							Text:       constants.EXAMPLE_MESSAGE_3,
							QuickReply: quickReply,
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

		sendMessage(userDoc, ctx, text)
		sendMessage(userDoc, ctx, message.Response)

	}()

	c <- res

	if _, err = app.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text:       res.Response,
					QuickReply: utils.CreateQuickReply([]string{"เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"}),
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

func sendMessage(userDoc *firestore.DocumentRef, ctx context.Context, message string) {
	_, _, err := userDoc.Collection("messages").Add(ctx, map[string]any{
		"from":      "user",
		"message":   message,
		"timestamp": firestore.ServerTimestamp,
	})
	if err != nil {
		log.Fatal(err)
	}
}
