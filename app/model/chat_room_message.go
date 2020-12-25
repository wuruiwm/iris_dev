package model

import "iris_dev/common"

type ChatRoomMessage struct {
	ID         uint   `gorm:"column:id;not null;primary_key;AUTO_INCREMENT;type:int(11)" json:"id"`
	ChatRoomID int    `gorm:"column:chat_room_id;not null;default:0;comment:'聊天室id';type:int(11)" json:"chat_room_id"`
	Content    string `gorm:"column:content;not null;default:'';comment:'消息内容';type:varchar(1000)" json:"content"`
	CreateTime int    `gorm:"column:create_time;not null;default:0;comment:'创建时间';type:int(11)" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time;not null;default:0;comment:'修改时间';type:int(11)" json:"update_time"`
}

func (*ChatRoomMessage) TableName() string {
	return `chat_room_message`
}

func ChatRoomMessageSave(id int,content string)(err error){
	chatRoomMessage := ChatRoomMessage{
		ChatRoomID: id,
		Content: content,
		CreateTime: common.GetUnixTime(),
		UpdateTime: common.GetUnixTime(),
	}
	return db.Create(&chatRoomMessage).Error
}