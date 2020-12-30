package websocket

import (
	"errors"
	gorillaWs "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	"iris_dev/app/model"
	"log"
	"net/http"
	"sync"
	"time"
)
var mutex sync.Mutex
var chatClientList = make(map[string]*Client)

type Client struct {
	ChatRoomId int
	Conn *websocket.Conn
	ConnId string
}

var wsUpgrade = gorilla.Upgrader(gorillaWs.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
})

func Chat()*neffos.Server{
	ws := websocket.New(wsUpgrade, websocket.Events{
		//收到消息时
		websocket.OnNativeMessage: onMessage(),
	})
	//连接时
	ws.OnConnect = onConnect()
	//关闭连接时
	ws.OnDisconnect = onDisConnect()
	//升级到websocket协议失败时
	ws.OnUpgradeError = onUpgradeError()
	return ws
}

func onMessage()func(nsConn *websocket.NSConn, msg websocket.Message) error{
	return func(nsConn *websocket.NSConn, msg websocket.Message) error {
		log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())
		content := string(msg.Body)
		client,ok := chatClientList[nsConn.Conn.ID()]//获取聊天室id
		if !ok{
			return errors.New("获取当前连接失败")
		}
		//将消息存入mysql
		_ = model.ChatRoomMessageSave(client.ChatRoomId,content)
		//发送消息
		SendMessage(nsConn.Conn,content)
		//获取所有连接
		for k,v := range chatClientList{
			//不推送消息给自己
			if k == client.ConnId{
				continue
			}
			//只推送给同聊天室的客户端
			if v.ChatRoomId != client.ChatRoomId{
				continue
			}
			SendMessage(v.Conn,content)
		}
		return nil
	}
}

func onUpgradeError()func(err error){
	return func(err error) {
		log.Printf("Upgrade Error: %v", err)
	}
}

func onDisConnect()func(conn *websocket.Conn){
	return func(c *websocket.Conn) {
		deleteChatClient(c)//删除客户端连接
		log.Printf("[%s] Disconnected from server", c.ID())
	}
}

func onConnect()func(conn *websocket.Conn)error{
	return func(c *websocket.Conn) error {
		//心跳包
		go heartbeat(c)
		ctx := websocket.GetContext(c)//获取content
		id := ctx.Params().GetIntDefault("id",0)//获取聊天室id
		addChatClient(c,id)//保存客户端连接
		log.Printf("[%s] Connected to server!", c.ID())
		return nil
	}
}

func heartbeat(c *websocket.Conn){
	for{
		time.Sleep(time.Second*5)
		c.Write(websocket.Message{
			Body:[]byte("heartbeat"),
			IsNative:true,
		})
	}
}

func addChatClient(c *websocket.Conn,chatRoomId int){
	client := &Client{
		Conn:c,
		ConnId:c.ID(),
		ChatRoomId:chatRoomId,
	}
	mutex.Lock()
	defer mutex.Unlock()
	chatClientList[c.ID()] = client
}

func deleteChatClient(c *websocket.Conn){
	mutex.Lock()
	defer mutex.Unlock()
	delete(chatClientList,c.ID())
}

func GetChatClientList()map[string]*Client{
	return chatClientList
}

func SendMessage(c *websocket.Conn,content string){
	data := websocket.Message{
		Body:[]byte(content),
		IsNative:true,
	}
	c.Write(data)
}