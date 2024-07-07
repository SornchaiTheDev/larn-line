package utils

import (
	"log"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

func CreateCallLarnMessage() *messaging_api.FlexMessage {
	jsonString := `{
  "type": "bubble",
  "body": {
    "type": "box",
    "layout": "vertical",
    "contents": [
      {
        "type": "text",
        "text": "เรียกหลาน",
        "weight": "bold",
        "size": "xl"
      },
      {
        "type": "text",
        "text": "สามารถเพิ่มเพื่อนโดยกดปุ่มข้างล่าง เพื่อคุยกับหลานตัวเป็น ๆ",
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
          "label": "เพิ่มเพื่อน",
          "uri": "https://line.me/ti/p/lP-CvWiMKT"
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
		AltText:    "โทรหาหลาน",
		Contents:   contents,
		QuickReply: CreateQuickReply([]string{"โทรหาหลาน", "เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"}),
	}
}
