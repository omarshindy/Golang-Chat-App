package entities

type Chat struct {
	Base
	ApplicationToken string      `gorm:"column:application_token" json:"-"`
	Application      Application `gorm:"references:Token"`
	Number           int32
	MessagesCount    int32 `gorm:"column:messages_count"`
}
