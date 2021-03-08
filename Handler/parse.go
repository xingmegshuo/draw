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
		log.Println("登录操作:")
		mes := Login(info, ws)
		Send(ws, mes)
	case "upgrade":
		log.Println("账号升级")
		mes := Upgrade(info)
		Send(ws, mes)
	// case "back":
	// 	log.Println("获取背包")
	// 	mes := GetBack(info)
	// 	Send(ws, mes)
	// case "addback":
	// 	log.Println("购买商品，增加背包")
	// 	mes := AddBack(info)
	// 	Send(ws, mes)
	case "record":
		log.Println("获取最近战绩")
		mes := GetRecord(info)
		Send(ws, mes)
	case "recordRate":
		log.Println("获取全部战斗")
		mes := GetRecordAll(info)
		Send(ws, mes)
	case "buddy":
		log.Println("获取好友列表")
		mes := GetBuddy(info)
		Send(ws, mes)
	case "newbuddy":
		log.Println("获取好友申请")
		mes := GetNewBuddy(info)
		Send(ws, mes)
	case "agreebuddy":
		log.Println("同意好友申请")
		mes := AgreeBuddy(info)
		Send(ws, mes)
	case "rcombuddy":
		log.Println("获取好友推荐")
		mes := RecomBuddy(info)
		Send(ws, mes)
	case "addbuddy":
		log.Println("添加好友申请")
		mes := AddBuddy(info)
		Send(ws, mes)
	case "delbuddy":
		log.Println("删除好友")
		mes := DeleteBuddy(info)
		Send(ws, mes)
	case "chat":
		log.Println("好友聊天")
		mes := Chat(info)
		Send(ws, mes)
	case "getUser":
		log.Println("获取用户信息")
		mes := GetUserMes(info)
		Send(ws, mes)
	case "room":
		log.Println("开始游戏")
		mes := GameStart(info, ws)
		Send(ws, mes)
	case "gaming":
		go RoomSocket(info)
	}
}

// 关闭连接时退出用户
func CloseUser(ws *websocket.Conn) {
	delete(client_user, ws)
	delete(client_palyer, ws)
	ws.Close()
	log.Println(len(client_user), len(client_palyer), "现有的链接数量")
	RemoveRoom()
}

// 数据返回
func Send(ws *websocket.Conn, mes string) {
	if err := websocket.Message.Send(ws, mes); err != nil {
		log.Println("客户端丢失", err.Error())
		CloseUser(ws)
	}
}

// // 清除连接
// func Clear(clientMap map[*websocket.Conn]string) {
// 	for ws, _ := range client_palyer {
// 		if _, ok := clientMap[ws]; !ok {
// 			log.Println("清除无效链接")
// 			CloseUser(ws)
// 		}
// 	}
// }

// 移除房间的人
func RemoveRoom() {
	for i, ro := range PlayRoom {
		for {
			l := IsUser(ro)
			if l == -1 {
				break
			}
			ro = DeleteUser(ro, l)
		}
		log.Println("还有几个人", len(ro.User))
		if len(client_palyer) == 0 {
			delete(PlayRoom, i)
		} else {
			if len(ro.User) > 1 {
				// RoomUser(ro)
				// str := "{'status':'system','mes':'系统消息','data':{'message':'" + "房间公告:" + openID + "退出房间'}}"
				// str = strings.Replace(str, "'", "\"", -1)
				// ServerRoom(ro, str)
				PlayRoom[i] = ro
			} else {
				delete(PlayRoom, i)
			}
		}
	}
}

// 移除无效用户
func DeleteUser(room Room, l int) Room {
	if l != -1 {
		room.User = append(room.User[:l], room.User[l+1:]...)
	}
	return room
}

// 是否有无效用户
func IsUser(room Room) int {
	for l, user := range room.User {
		if _, ok := client_palyer[user.Ws]; ok {
			log.Println("正常---", user.OpenID)
		} else {
			log.Println("此用户断开链接", user.OpenID)
			room.People = room.People + 1
			return l
		}
	}
	return -1
}
