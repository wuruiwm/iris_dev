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
	"time"
)

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
		id := nsConn.Conn.Get("id")//获取聊天室id
		idInt,ok := id.(int)
		if !ok{
			return errors.New("获取聊天室id失败")
		}
		content := string(msg.Body)
		//将消息存入mysql
		_ = model.ChatRoomMessageSave(idInt,content)
		//发送的消息体
		data := websocket.Message{
			Body:[]byte(content),
			IsNative:true,
		}
		nsConn.Conn.Write(data)
		//获取所有连接
		connList := nsConn.Conn.Server().GetConnections()
		for k,v := range connList{
			//不推送消息给自己
			if k == nsConn.Conn.ID(){
				continue
			}
			//只推送给同聊天室的客户端
			if v.Get("id") != id{
				continue
			}
			v.Write(data)
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
		log.Printf("[%s] Disconnected from server", c.ID())
	}
}

func onConnect()func(conn *websocket.Conn)error{
	return func(c *websocket.Conn) error {
		//心跳包
		go heartbeat(c)
		ctx := websocket.GetContext(c)//获取content
		id := ctx.Params().GetIntDefault("id",0)//获取聊天室id
		c.Set("id",id)
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