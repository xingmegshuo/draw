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
	"time"

	"golang.org/x/net/websocket"
)

type GameType struct {
	New string
}

type Player struct {
	OpenID string
	Ws     *websocket.Conn
	Ready  string
	Score  int
	Status int
}

type Room struct {
	People int
	Public int
	User   []Player
	Owner  string
}

type Mes struct {
	RoomID  string
	Message string
	User    string
}

// 房间
var PlayRoom = make(map[string]Room)

// 是否新建房间的数据解析
var Game GameType

// 开始游戏
func GameStart(mes []byte, ws *websocket.Conn) string {
	var room Room
	err := json.Unmarshal(mes, &Game)
	if err != nil {
		log.Println("数据问题:", err.Error())
	}
	if Game.New == "true" {
		room = SearchRoom()
	} else {
		room = NewRoom()
	}
	message := Init(ws, room)
	return ToMes("ok", message)
}

// 查找房间
func SearchRoom() Room {
	index := "-1"
	for l, item := range PlayRoom {
		if item.People > 0 && item.Public == 0 {
			index = l
		}
	}
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
		Public: 0,
	}
	return ro
}

// 初始化房间
func Init(ws *websocket.Conn, room Room) string {
	userID := client_user[ws]
	player := Player{
		OpenID: client_user[ws],
		Ws:     ws,
		Ready: "false",
	}
	room.User[len(room.User)] = player
	room.People = room.People - 1
	client_user[ws] = client_palyer[ws]
	delete(client_user, ws)
	if room.People == 6 {
		room.Owner = userID
		PlayRoom[room.Owner] = room
	} else {
		for l, ro := range PlayRoom {
			if ro.Owner == room.Owner {
				PlayRoom[l] = room
			}
		}
	}
	ServerRoom(room, "房间公告:"+userID+"进入房间")
	RoomUser(room)
	return room.Owner
}

// 发送房间成员和状态
func RoomUser(room Room){
	str := "{'status':'room','mes':'房间成员信息','data':{"
	for l, item : = range room.User{
		if l == len(room.User)&& l > 0 {
			str=str+"'"+item.OpenID+"':'"+item.Ready+"',"
		}else{
			str=str+"'"+item.OpenID+"':'"+item.Ready+"'}"
		}
	}
	str = strings.Replace(str, "'", "\"", -1)
	ServerRoom(room,str)
}


// 房间内信息
func ServerRoom(room Room, mes string) {
	for _, item := range room.User {
		SendMES(item.Ws, mes)
	}
}

// 用户离线
func OutLine(ws *websocket.Conn) {
	for l, ro := range PlayRoom {
		for n, item := range ro.User {
			if item.Ws == ws {
				item.Status = 000
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
