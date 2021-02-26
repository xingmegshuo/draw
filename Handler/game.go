/***************************
@File        : game.go
@Time        : 2021/01/19 14:40:30
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 游戏逻辑
****************************/

package Handler

import (
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

type Player struct {
	OpenID string
	Ws     *websocket.Conn
	Ready  string
	Score  int
	Status string
}

type Room struct {
	People int
	Public string
	User   []Player
	Owner  string
}

type Mes struct {
	RoomID  string
	Message string
	User    string
	Data    string
}
type GameType struct {
	Type string
}

// 房间
var PlayRoom = make(map[string]Room)

// 是否新建房间的数据解析
var Game GameType

// 开始游戏- 返回房间号
func GameStart(mes []byte, ws *websocket.Conn) string {
	var room Room
	err := json.Unmarshal(mes, &Game)
	if err != nil {
		log.Println("解析room:", err.Error())
	}
	if Game.Type == "false" {
		room = SearchRoom()
	} else {
		room = NewRoom()
	}
	message := Init(ws, room)
	str := "{'status':'ok','mes':'房间号','data':{'message':'" + message + "'}}"
	str = strings.Replace(str, "'", "\"", -1)
	return str
}

// 查找房间
func SearchRoom() Room {
	index := "-1"
	for l, item := range PlayRoom {
		log.Println(l, "--------房间id")
		if item.People > 0 && item.Public == "true" {
			index = l
		}
	}
	log.Println("查找房间----------", index)
	if index != "-1" {
		return PlayRoom[index]
	} else {
		return NewRoom()
	}
}

// 新建房间
func NewRoom() Room {
	ro := Room{
		People: 6,
		Public: "true",
	}
	return ro
}

// 初始化房间
func Init(ws *websocket.Conn, room Room) string {
	userID := client_user[ws]
	player := Player{
		OpenID: client_user[ws],
		Ws:     ws,
		Ready:  "false",
		Status: "true",
	}
	add := "false"
	for l, item := range room.User {
		if item.OpenID == player.OpenID {
			add = "True"
			room.User[l] = player
		}
	}
	log.Println(player.OpenID, "-----------用户ID")
	if add == "false" {
		room.User = append(room.User, player)
	}
	// room.User[len(room.User)] = player
	client_user[ws] = client_palyer[ws]
	room.People = room.People - 1
	delete(client_user, ws)
	if room.People == 5 {
		room.Owner = userID
		PlayRoom[room.Owner] = room
	} else {
		for l, ro := range PlayRoom {
			if ro.Owner == room.Owner {
				PlayRoom[l] = room
			}
		}
	}
	log.Println("------------房间人员", room.User, "******", len(PlayRoom))
	str := "{'status':'system','mes':'系统消息','data':{'message':'" + "房间公告:" + userID + "进入房间'}}"
	str = strings.Replace(str, "'", "\"", -1)
	ServerRoom(room, str)
	RoomUser(room)
	return room.Owner
}

// 发送房间成员和状态
func RoomUser(room Room) {
	str := "{'status':'room','mes':'房间成员信息','data':["
	for l, item := range room.User {
		if l == len(room.User)-1 {
			str = str + "{'user':'" + item.OpenID + "','ready':'" + item.Ready + "','onLine':'" + item.Status + "'}"
		} else {
			str = str + "{'user':'" + item.OpenID + "','ready':'" + item.Ready + "','onLine':'" + item.Status + "'},"
		}

	}
	str = str + "]}"
	str = strings.Replace(str, "'", "\"", -1)
	ServerRoom(room, str)
}

// 房间内信息
func ServerRoom(room Room, mes string) {
	for _, item := range room.User {
		if item.Status != "false" {
			SendMES(item.Ws, mes)
		}
	}
}

// 用户离线
func OutLine(ws *websocket.Conn) {
	for l, ro := range PlayRoom {
		for n, item := range ro.User {
			if item.Ws == ws {
				item.Status = "false"
				ro.User[n] = item
			}
		}
		PlayRoom[l] = ro
	}
}

// 发消息
func SendMES(ws *websocket.Conn, mes string) {
	if err := websocket.Message.Send(ws, mes); err != nil {
		log.Println("用户离线", err.Error())
		OutLine(ws)
	}
}

// 更新房间到房间列表
func UpdatePlayRoom(room Room) {
	for l, item := range PlayRoom {
		if item.Owner == room.Owner {
			PlayRoom[l] = room
		}
		if item.People == 6 {
			log.Println("删除房间-----------")
			delete(PlayRoom, l)
		}
	}
}

// 用户准备
func Ready(room Room, user string) {
	for l, item := range room.User {
		if item.OpenID == user {
			item.Ready = "true"
			room.User[l] = item
		}
	}
	RoomUser(room)
	UpdatePlayRoom(room)
}

// 退出房间
func Leave(room Room, user string) {
	str := "{'status':'system','mes':'系统消息','data':{'message':'" + "房间公告:" + user + "退出房间'}}"
	str = strings.Replace(str, "'", "\"", -1)
	ServerRoom(room, str)
	for l, item := range room.User {
		if item.OpenID == user {
			room.User = append(room.User[:l], room.User[l+1:]...)
			room.People = room.People + 1
			client_user[item.Ws] = client_palyer[item.Ws]
			delete(client_palyer, item.Ws)
		}
	}
	RoomUser(room)
	UpdatePlayRoom(room)
}

// 房间内消息
func RoomSocket(mes []byte) {
	log.Println("--------------------房间")
	var Msg Mes
	var room Room
	err := json.Unmarshal(mes, &Msg)
	if err != nil {
		log.Println("数据问题:", err.Error())
	}
	for l, item := range PlayRoom {
		if l == Msg.RoomID {
			room = item
		}
	}
	log.Println(Msg.Message, "----------------")
	switch Msg.Message {
	case "ready":
		Ready(room, Msg.User)
	case "send":
		log.Println("发送消息-------------------------")
		str := "{'status':'room','mes':'房间转发信息','data':{'message':'" + Msg.Data + "'}}"
		str = strings.Replace(str, "'", "\"", -1)
		ServerRoom(room, str)
	case "leave":
		log.Println("退出房间----------------")
		Leave(room,Msg.User)
	}
	
}
