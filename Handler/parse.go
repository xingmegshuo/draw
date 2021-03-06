/***************************
@File        : parse.go
@Time        : 2020/12/21 15:56:57
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 解析数据根据不同key来进行区分
****************************/

package Handler

import (
	"encoding/json"
	"log"

	"golang.org/x/net/websocket"
)

// 数据格式
type Data struct {
	Name   string
	Values string
}

var client_user = make(map[*websocket.Conn]string)
var client_palyer = make(map[*websocket.Conn]string)

// 房间
var PlayRoom = make(map[string]Room)

// var client_buddy = make(map[*websocket.Conn]string)

// 解析key
func ParseData(con string, ws *websocket.Conn) {
	// log.Println("开始解析", con)
	var data Data
	// 清除连接
	// Clear(clientMap)
	oldData := []byte(con)
	err := json.Unmarshal(oldData, &data)
	if err != nil {
		log.Println(err)
		Send(ws, "解析数据失败,请查看数据格式")
	}
	// log.Printf("类型%T", data.Values)
	info := []byte(data.Values)
	// log.Println(info, "-------", data.Values)
	switch data.Name {
	case "login":
		// log.Println("登录操作:")
		mes := Login(info, ws)
		Send(ws, mes)
	case "like":
		// log.Println("点赞")
		mes := Upgrade(info)
		Send(ws, mes)
	case "clear":
		mes := ClearData(info)
		Send(ws, mes)
	// case "record":
	// 	// log.Println("获取最近战绩")
	// 	mes := GetRecord(info)
	// 	Send(ws, mes)
	// case "recordRate":
	// 	// log.Println("获取全部战斗")
	// 	mes := GetRecordAll(info)
	// 	Send(ws, mes)
	// case "buddy":
	// 	// log.Println("获取好友列表")
	// 	mes := GetBuddy(info)
	// 	Send(ws, mes)
	// case "newbuddy":
	// 	// log.Println("获取好友申请")
	// 	mes := GetNewBuddy(info)
	// 	Send(ws, mes)
	// case "agreebuddy":
	// 	// log.Println("同意好友申请")
	// 	mes := AgreeBuddy(info)
	// 	Send(ws, mes)
	// case "rcombuddy":
	// 	// log.Println("获取好友推荐")
	// 	mes := RecomBuddy(info)
	// 	Send(ws, mes)
	// case "addbuddy":
	// 	// log.Println("添加好友申请")
	// 	mes := AddBuddy(info)
	// 	Send(ws, mes)
	// case "delbuddy":
	// 	// log.Println("删除好友")
	// 	mes := DeleteBuddy(info)
	// 	Send(ws, mes)
	case "chat":
		// log.Println("好友聊天")
		mes := Chat(info)
		Send(ws, mes)
	case "getUser":
		// log.Println("获取用户信息")
		mes := GetUserMes(info)
		Send(ws, mes)
	case "room":
		// log.Println("开始游戏")
		mes := GameStart(info, ws)
		Send(ws, mes)
	case "gaming":
		RoomSocket(info)
	}
}

// 关闭连接时
func CloseUser(ws *websocket.Conn) {
	RemoveRoom(ws)
	// delete(client_user, ws)
}

// 数据返回
func Send(ws *websocket.Conn, mes string) {
	if err := websocket.Message.Send(ws, mes); err != nil {
		// log.Println("客户端丢失", err.Error())
		CloseUser(ws)
		ws.Close()
	}
}

// 移除房间的人
func RemoveRoom(ws *websocket.Conn) {
	for _, ro := range PlayRoom {
		for _, a := range ro.User {
			if a.Ws == ws {
				Leave(ro, a.OpenID)
				log.Println("断线退出")
			}
		}
	}
}
