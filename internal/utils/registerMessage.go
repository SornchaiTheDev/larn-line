package utils

import (
	"log"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

func CreateRegisterMessage() *messaging_api.FlexMessage {
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
          "uri": "https://line.sornchaithedev.com"
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
		QuickReply: CreateQuickReply([]string{"เพิ่มขนาดตัวอักษร", "ตั้งค่าการแจ้งเตือนให้มีเสียงดังขึ้น", "วิธีถ่ายภาพหน้าจอ", "จะส่งรูปภาพทางไลน์", "วิธีตั้งนาฬิกาปลุก", "เชื่อม WiFi กับโทรศัพท์", "ลบแอปพลิเคชัน", "เปิดใช้งานโหมดประหยัดแบตเตอรี่"}),
	}
}
