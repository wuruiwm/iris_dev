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

var mutex sync.Mutex //对chatClientList操作时 上锁
var chatClientList = make(map[string]*Client) //客户端连接集合
var handleClientListNum = 0 //因为map的delete不会真正的释放内存 所以在操作一定次数后 重新初始化一次map 释放内存
var initClientListNum = 10000 //多少次以后 重新初始化map 释放内存

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
		websocket.OnNativeMessage: onMessage,
	})
	//连接时
	ws.OnConnect = onConnect
	//关闭连接时
	ws.OnDisconnect = onDisConnect
	//升级到websocket协议失败时
	ws.OnUpgradeError = onUpgradeError
	return ws
}

func onMessage(nsConn *websocket.NSConn, msg websocket.Message)error{
	log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())
	content := string(msg.Body)
	client,ok := chatClientList[nsConn.Conn.ID()]//获取聊天室id
	if !ok{
		return errors.New("获取当前连接失败")
	}
	//将消息存入mysql
	_ = model.ChatRoomMessageSave(client.ChatRoomId,content)
	//发送消息
	client.SendMessage(content)
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
		v.SendMessage(content)
	}
	return nil
}

func onUpgradeError(err error){
	log.Printf("Upgrade Error: %v", err)
}

func onDisConnect(c *websocket.Conn){
	deleteChatClient(c)//删除客户端连接
	log.Printf("[%s] Disconnected from server", c.ID())
}

func onConnect(c *websocket.Conn)error{
	//心跳包
	go heartbeat(c)
	ctx := websocket.GetContext(c)//获取content
	id := ctx.Params().GetIntDefault("id",0)//获取聊天室id
	addChatClient(c,id)//保存客户端连接
	log.Printf("[%s] Connected to server!", c.ID())
	return nil
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
	handleClientListNum++
	if handleClientListNum < initClientListNum{
		return
	}
	handleClientListNum = 0
	var mapTmp = make(map[string]*Client)
	for k,v := range chatClientList{
		mapTmp[k] = v
	}
	chatClientList = mapTmp
}

func GetChatClientList()map[string]*Client{
	return chatClientList
}

func (client *Client)SendMessage(content string){
	data := websocket.Message{
		Body:[]byte(content),
		IsNative:true,
	}
	client.Conn.Write(data)
}