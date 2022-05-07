package entities

type Message struct {
	Base
	MessageBody string `gorm:"column:message_body" `
	Number      int    `gorm:"column:number" json:"-"`
	ChatID      int32  `gorm:"column:chat_id" json:"-"`
	Chat        Chat   `json:"-"`
}
