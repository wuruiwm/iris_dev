package admin

import (
	"github.com/kataras/iris/v12"
	"iris_dev/app/response"
	"iris_dev/app/websocket"
)

func ChatIndex(c iris.Context){
	_ = c.View("admin/chat/index.html")
}

func TestSendMessage(c iris.Context){
	msg := c.PostValue("msg")
	chatRoomId := c.PostValueIntDefault("id",0)
	if chatRoomId <= 0 {
		response.Error(c,"聊天室id错误")
		return
	}
	chatClientList := websocket.GetChatClientList()
	for _,v := range chatClientList{
		if v.ChatRoomId == chatRoomId{
			v.SendMessage(msg)
		}
	}
	response.Success(c,"发送成功",nil)
}