/***************************
@File        : game.go
@Time        : 2021/01/19 14:40:30
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 游戏逻辑
****************************/

package Handler

import (
	"crypto/rand"
	"draw/Mydb"
	"encoding/json"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Player struct {
	OpenID string
	Ws     *websocket.Conn
	Ready  string
	Score  int
	Status string
	Ok     string
}

type Room struct {
	People      int
	Public      string
	User        []Player
	Owner       string
	Word        string
	Draw        string
	Status      bool
	GuessPeople int
	ID          string
}

type Mes struct {
	RoomID  string
	Message string
	User    string
	Data    string
}
type GameType struct {
	Type   string
	RoomID string
	Search string
}

var Wordctrl = Mydb.NewGuessAnswerCtrl()
var Userctrl = Mydb.NewUserCtrl()

// 是否新建房间的数据解析
var Game GameType

// 开始游戏- 返回房间号
func GameStart(mes []byte, ws *websocket.Conn) string {
	Game.Search = "false"
	var room Room
	err := json.Unmarshal(mes, &Game)
	if err != nil {
		log.Println("解析room:", err.Error())
	}
	if Game.Type == "false" {
		ro, ok := SearchRoom(Game.RoomID)
		if Game.Search == "true" {
			if !ok {
				return StrToJSON("RoomError", "房间", "{'message':'加入房间失败，房间不存在或房间正在游戏中'}")
			}
		}
		room = ro
	} else {
		room = NewRoom()
	}
	message := Init(ws, room)
	return StrToJSON("ok", "房间号", "{'message':'"+message+"'}")
}

// 断线重连
func IsIn(user string, ws *websocket.Conn) Player {
	for _, ro := range PlayRoom {
		for _, item := range ro.User {
			if item.OpenID == user {
				log.Println("欢迎重新链接", item.OpenID)
				item.Ws = ws
				return item
			}
		}
	}
	return Player{}
}

// 查找房间
func SearchRoom(roomID string) (Room, bool) {
	for l, item := range PlayRoom {
		if Game.Search == "true" {
			log.Println("携带房间号搜索--------------------------")
			if l == roomID && item.Status == true && item.People > 0 {
				return PlayRoom[l], true
			}
		} else {
			log.Println("系统匹配进入房间----------------------")
			if item.People > 0 && item.Public == "true" && item.Status == true {
				return PlayRoom[l], true
			}
		}
	}
	if Game.Search == "true" {
		return Room{}, false
	} else {
		return NewRoom(), false
	}
}

// 新建房间
func NewRoom() Room {
	log.Println("新建房间---------------------------")
	Id := ""
	for {
		number := makeRoomID()
		if _, ok := PlayRoom[number]; !ok {
			Id = number
			break
		}
	}
	ro := Room{
		People: 6,
		Public: "true",
		Status: true,
		ID:     Id,
	}
	return ro
}

// 初始化房间
func Init(ws *websocket.Conn, room Room) string {
	client_palyer[ws] = client_user[ws]
	// delete(client_user, ws)
	player := Player{
		OpenID: client_palyer[ws],
		Ws:     ws,
		Ready:  "true",
		Status: "true",
	}
	add := true
	for l, u := range room.User {
		if u.OpenID == player.OpenID {
			log.Println("房间存在此用户----------------------------")
			room.User[l] = player
			add = false
		}
	}
	if add == true {
		room.User = append(room.User, player)
		room.People = room.People - 1
	}
	if room.People == 5 {
		room.Owner = room.User[0].OpenID
	}
	UpdatePlayRoom(room)
	room = GetRoom(room)
	RoomUser(room)
	log.Println("加入后:-----------", client_palyer)
	return GetRoomID(room)
}

// 获取房间号
func GetRoomID(room Room) string {
	for l, ro := range PlayRoom {
		if ro.Owner == room.Owner {
			return l
		}
	}
	log.Println("no room --------------------------------------------")
	return ""
}

// 发送房间成员和状态
func RoomUser(room Room) {
	str := "{'status':'room','mes':'房间成员信息','data':["
	for l, item := range room.User {
		thisUser := Mydb.User{
			OpenID: item.OpenID,
		}
		user, _ := Userctrl.GetUser(thisUser)
		if l == len(room.User)-1 {
			if item.OpenID == room.Owner {
				str = str + "{'user':'" + item.OpenID + "','nickName':'" + user.NickName + "','avatarUrl':'" + user.AvatarURL + "','orther':'" + user.Orther + "','ready':'" + item.Ready + "','homeowner':'" + item.OpenID + "','online':'" + item.Status + "'}"
			} else {
				str = str + "{'user':'" + item.OpenID + "','nickName':'" + user.NickName + "','avatarUrl':'" + user.AvatarURL + "','orther':'" + user.Orther + "','ready':'" + item.Ready + "','online':'" + item.Status + "'}"
			}
		} else {
			if item.OpenID == room.Owner {
				str = str + "{'user':'" + item.OpenID + "','nickName':'" + user.NickName + "','avatarUrl':'" + user.AvatarURL + "','orther':'" + user.Orther + "','ready':'" + item.Ready + "','homeowner':'" + item.OpenID + "','online':'" + item.Status + "'},"
			} else {
				str = str + "{'user':'" + item.OpenID + "','nickName':'" + user.NickName + "','avatarUrl':'" + user.AvatarURL + "','orther':'" + user.Orther + "','ready':'" + item.Ready + "','online':'" + item.Status + "'},"
			}
		}
	}
	str = str + "]}"
	str = strings.Replace(str, "'", "\"", -1)
	ServerRoom(room, str)
}

// 房间内信息
func ServerRoom(room Room, mes string) {
	room = GetRoom(room)
	// log.Println(len(room.User), "发送消息几个人")
	for _, item := range room.User {
		// log.Println(item.Status)
		if item.Status != "false" {
			Send(item.Ws, mes)
		}
	}
}

// 用户离线
func OutLine(ws *websocket.Conn) {
	for _, ro := range PlayRoom {
		for l, item := range ro.User {
			if item.Ws == ws {
				item.Status = "false"
				log.Println("给他掉线", item.OpenID)
			}
			ro.User[l] = item
			UpdatePlayRoom(ro)
		}
	}
}

// 更新房间到房间列表
func UpdatePlayRoom(room Room) {
	b := false
	for l, item := range PlayRoom {
		if item.ID == room.ID {
			PlayRoom[l] = room
			b = true
		}
	}
	if b == false {
		PlayRoom[room.ID] = room
	}
	// log.Println("几个房间----------", len(PlayRoom))
}

// 生成房间号
func makeRoomID() string {
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(10000)))
	return strconv.Itoa(int(result.Int64()))
}

// 用户准备
func Ready(room Room, user string, status string) {
	for l, item := range room.User {
		if item.OpenID == user {
			if status == "false" {
				item.Ready = "false"
				str := "{'status':'system','mes':'系统消息','data':{'message':'" + "房间公告:" + GetNickName(user) + "取消准备'}}"
				str = strings.Replace(str, "'", "\"", -1)
				ServerRoom(room, str)
			}
			if status == "true" {
				item.Ready = "true"
				str := "{'status':'system','mes':'系统消息','data':{'message':'" + "房间公告:" + GetNickName(user) + "已经准备'}}"
				str = strings.Replace(str, "'", "\"", -1)
				ServerRoom(room, str)
			}
			room.User[l] = item
		}
	}
	UpdatePlayRoom(room)
	RoomUser(room)
}

// 退出房间
func Leave(room Room, user string) {
	for _, item := range room.User {
		if item.OpenID == user {
			delete(client_palyer, item.Ws)
		}
	}
	if len(room.User) <= 1 {
		for l, ro := range PlayRoom {
			if ro.ID == room.ID {
				delete(PlayRoom, l)
				log.Println("删除房间--------------------")
			}
		}
	} else {
		change_owner := false
		for {
			a := -1
			for l, item := range room.User {
				if item.OpenID == user {
					a = l
					if room.Owner == user {
						change_owner = true
					}
				}
			}
			if a != -1 {
				room.User = append(room.User[:a], room.User[a+1:]...)
			}
			if a == -1 {
				break
			}
			if change_owner == true {
				oldOwner := room.Owner
				room.Owner = room.User[0].OpenID
				newOwner := room.Owner
				changeOwner(oldOwner, newOwner)
			}
		}
		room.People = 6 - len(room.User)
		UpdatePlayRoom(room)
		room = GetRoom(room)
		log.Println("退出房间----------------------房间人数", len(room.User))
		RoomUser(room)
		ServerRoom(room, StrToJSON("system", "系统消息", "{'message':'房间公告:"+user+"退出房间'}"))
	}
	log.Println("退出后:", client_palyer)
}

// 修改房主
func changeOwner(old string, new string) {
	for l, ro := range PlayRoom {
		if ro.Owner == old {
			ro.Owner = new
		}
		PlayRoom[l] = ro
	}
}

// 房间内消息
func RoomSocket(mes []byte) {
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
	switch Msg.Message {
	case "ready":
		go Ready(room, Msg.User, Msg.Data)
	case "send":
		go ServerRoom(room, StrToJSON("room", "房间转发信息", "{'message':':"+Msg.Data+"'}"))
	case "leave":
		go Leave(room, Msg.User)
	case "start":
		go Start(room, Msg.User)
	case "word":

		go Word(room, Msg.User)
	case "choose":
		// log.Println("选词----------------", room.Owner, Msg.Data)
		str := strings.Replace(Msg.Data, "\"", "", -1)
		go Choose(room, str)
	case "guess":
		// log.Println("猜词---------------", room.Word)
		str := strings.Replace(Msg.Data, "\"", "", -1)
		go Guess(room, Msg.User, str)
	}
}

// 猜词语
func Guess(room Room, user string, word string) {
	room = GetRoom(room)
	str := ""
	if room.Status == true {
		str = word
	} else {
		str = word
		if len(room.Word) > 0 {
			for _, i := range room.Word {
				str = strings.Replace(str, string(i), "*", -1)
			}
		}
		add := false
		for l, item := range room.User {
			if item.OpenID == user && room.Draw != user && room.GuessPeople > 0 {
				if len(word) <= 12 {
					if word == room.Word && item.Ok != "true" {
						item.Score = item.Score + room.GuessPeople*2
						room.GuessPeople = room.GuessPeople - 1
						add = true
						a := "{'status':'room','mes':'答对加分','data':{'message':" + "[{'user':'" + user + "','score':'" + strconv.Itoa(item.Score) + "'}]" + "}}"
						a = strings.Replace(a, "'", "\"", -1)
						ServerRoom(room, a)
						item.Ok = "true"
					} else {
						if item.Ok != "true" {
							a := "{'status':'system','mes':'答错了','data':{'message':'房间公告:" + GetNickName(item.OpenID) + "回答错误','user':'" + item.OpenID + "'}}"
							a = strings.Replace(a, "'", "\"", -1)
							ServerRoom(room, a)
						}
					}
				}

			}
			room.User[l] = item
		}
		for l, item := range room.User {
			if item.OpenID == room.Draw && add == true {
				item.Score = item.Score + 2
			}
			room.User[l] = item
		}
	}

	data := "{'status':'system','mes':'猜答案和交流','data':{'message':'" + str + "','user':'" + user + "'}}"
	data = strings.Replace(data, "'", "\"", -1)
	UpdatePlayRoom(room)
	ServerRoom(room, data)
}

// 格式化返回数据
func StrToJSON(status string, mes string, message string) string {
	str := "{'status':'" + status + "','mes':'" + mes + "','data':" + message + "}"
	str = strings.Replace(str, "'", "\"", -1)
	return str
}

// 给所有人准备
func ReadyAll(room Room) {
	for l, item := range room.User {
		if item.Ready == "false" {
			item.Ready = "true"
		}
		room.User[l] = item
	}
	UpdatePlayRoom(room)
	RoomUser(room)
}

// 给所有人取消准备
func UnReadyAll(room Room) {
	for l, item := range room.User {
		if item.Ready == "true" {
			item.Ready = "false"
		}
		room.User[l] = item
	}
	UpdatePlayRoom(room)
	RoomUser(room)
}

// 是否有人取消准备
func IsStart(room Room) bool {
	room = GetRoom(room)
	if len(room.User) < 2 {
		return false
	}
	for _, item := range room.User {
		if item.Ready == "false" {
			return false
		}
	}
	return true
}

// 倒计时
func IsStartUnderTime(count int, room Room, mes string) bool {
	ServerRoom(room, StrToJSON("room", "房间状态", mes))
	for i := 0; i < count; i++ {
		ServerRoom(room, StrToJSON("time", "系统时间提示", "{'message':'房间公告: 倒计时还有"+strconv.Itoa(count-i)+"秒'}"))
		ServerRoom(room, StrToJSON("room", "倒计时", "{'message':'"+strconv.Itoa(count-i)+"'}"))
		b := IsStart(room)
		// log.Println(b, "是否有人取消准备或退出")
		if b == false {
			// room = GetRoom(room)
			ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'StartCountdownStop'}"))
			ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'GameError'}"))
			ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 用户取消准备,游戏未能开始'}"))
			return false
		}
		time.Sleep(time.Second * 1)
	}
	return true
}

// 开始游戏
func Start(room Room, user string) {
	a := 0
	for _, item := range room.User {
		if item.Ready == "true" && item.OpenID != room.Owner {
			a = a + 1
		}
	}
	if room.Owner == user && a >= 1 && room.People <= 4 {
		ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 游戏五秒后开始'}"))
		status := IsStartUnderTime(5, room, "{'message':'GameCountdown'}")
		room = GetRoom(room)
		if status {
			ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'StartCountdownStop'}"))
			ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'GameSuccess'}"))
			ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 游戏正式开始'}"))
			log.Println("开始游戏---------------修改房间状态不可加入")
			room.Status = false
			UpdatePlayRoom(room)
			OneGame(room)
		}
	} else {
		ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'CountdownError'}"))
		ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'GameError'}"))
		ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 房间准备人数不足,请玩家准备'}"))
	}
}

// 随机生成四个词
func Word(room Room, user string) {
	str := "{'status':'room','mes':'词语','data':["
	words := make(map[int]string)
	a := 1
	for {
		w := GetWord()
		b := false
		for _, value := range words {
			if w == value {
				b = true
			}
		}
		if b == false {
			words[a] = w
			a = a + 1
		}
		if len(words) == 4 {
			break
		}
	}
	str = str + "'" + words[1] + "'," + "'" + words[2] + "'," + "'" + words[3] + "'," + "'" + words[4] + "'"
	str = str + "]}"
	str = strings.Replace(str, "'", "\"", -1)
	room = GetRoom(room)
	log.Println("发送四个词语调用------------------,获取词语", len(room.User), room.Draw)
	for _, item := range room.User {
		if item.OpenID == user && room.Draw == user {
			// log.Println("发送四个词语--------给谁发送", item.OpenID, str)
			Send(item.Ws, str)
		}
	}
}

// 从数据库中获取词语
func GetWord() string {
	SearchWord := Mydb.GuessAnswer{
		Id: 0,
	}
	words := Wordctrl.GetAnswer(SearchWord)
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
	return words[int(result.Int64())].Answer
}

// 选词
func Choose(room Room, word string) {
	room.Word = word
	ServerRoom(room, StrToJSON("system", "选择的词语", "{'message':'"+word+"'}"))
	ServerRoom(room, StrToJSON("room", "选词完毕状态", "{'message':'ok'}"))
	// log.Println(room.Word, "-----------------词语")
	UpdatePlayRoom(room)
}

// 更新房间
func GetRoom(room Room) Room {
	var ro Room
	for _, item := range PlayRoom {
		if item.ID == room.ID {
			ro = item
		}
	}
	return ro
}

// 选词倒计时
func ChooseWordUnderTime(count int, room Room, mes string) bool {
	ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'"+mes+"'}"))
	for i := 0; i < count; i++ {
		ServerRoom(room, StrToJSON("time", "系统时间提示", "{'message':'房间公告: 倒计时还有"+strconv.Itoa(count-i)+"秒'}"))
		ServerRoom(room, StrToJSON("room", "倒计时", "{'message':'"+strconv.Itoa(count-i)+"'}"))
		// log.Println(count - i)
		ro := GetRoom(room)
		if len(ro.User) < 2 {
			return false
		}
		if ro.Word != "" {
			log.Println("选词结束------------", len(ro.User))
			ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'ChooseCountdownStop'}"))
			ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 选词完毕'}"))
			return true
		}
		time.Sleep(time.Second * 1)
	}
	return false
}

// 游戏流程
func OneGame(room Room) {
	// go listen(len(room.User), room)
	for _, ro := range PlayRoom {
		if ro.Owner == room.Owner {
			for i, item := range ro.User {
				ro = GetRoom(ro)
				ro.GuessPeople = len(ro.User) - 1
				ro.Draw = item.OpenID
				UpdatePlayRoom(ro)
				ServerRoom(ro, StrToJSON("room", "画家", "{'message':'"+item.OpenID+"'}"))
				ServerRoom(ro, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 第"+strconv.Itoa(i+1)+"回合,画师为"+GetNickName(item.OpenID)+",请他开始选词'}"))
				log.Println(ro.Word, "是否携带词语-------------------------")
				w := ChooseWordUnderTime(10, ro, "ChooseWordCountdown")
				ro = GetRoom(ro)
				if len(ro.User) < 2 {
					// GameOver(ro)
					break
				}
				if w == false {
					log.Println("选词结束------------", len(ro.User))
					ServerRoom(ro, StrToJSON("room", "房间状态", "{'message':'ChooseCountdownStop'}"))
					Choose(ro, GetWord())
					ServerRoom(ro, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 选词完毕'}"))
				}
				time.Sleep(time.Second * 1)
				ok := RoundTime(90, ro)
				ro = GetRoom(ro)
				if ok == true {
					ServerRoom(ro, StrToJSON("room", "房间状态", "{'message':'DrawCountdownStop'}"))
					RoundOver(ro)
				}
			}
			GameOver(ro)
		}
	}
}

// 监听游戏中的退出
func listen(l int, ro Room) {
	for {
		ro = GetRoom(ro)
		if len(ro.User) != l {
			log.Println("游戏过程中发生了改变----------------------------------")
			break
		}
	}
}

// 回合结束
func RoundOver(room Room) {
	ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 画家已画完'}"))
	ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 本轮回合结束,正确答案"+room.Word+"'}"))
	ServerRoom(room, StrToJSON("room", "正确答案", "{'message':'"+room.Word+"'}"))
	for l, item := range room.User {
		item.Ok = ""
		room.User[l] = item
	}
	room.Word = ""
	room.GuessPeople = len(room.User) - 1
	UpdatePlayRoom(room)
	// log.Println("回合结束正确答案:", room.Word)
	ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 点赞开始'}"))
	ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'RoundOver'}"))
	UnderTime(5, room, "RoundCountdown")
	ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'RoundCountdownStop'}"))
	ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 点赞结束'}"))
}

// 写入数据
func Settlement(room Room) {
	for _, item := range room.User {
		user := Mydb.User{
			OpenID: item.OpenID,
		}
		thisUser, has := Userctrl.GetUser(user)
		if has {
			thisUser.Number = thisUser.Number + 1
			thisUser.Score = thisUser.Score + item.Score
			Userctrl.Update(thisUser)
		}
	}
}

// 游戏结束
func GameOver(room Room) {
	ServerRoom(room, StrToJSON("system", "系统提示信息", "{'message':'房间公告: 游戏结束'}"))
	ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'GameOver'}"))
	Settlement(room)
	str := "{'status':'room','mes':'游戏结算','data':["
	for l, item := range room.User {
		if l == len(room.User)-1 {
			str = str + "{'user':'" + item.OpenID + "','score':'" + strconv.Itoa(item.Score) + "'}"
			item.Score = 0
		} else {
			str = str + "{'user':'" + item.OpenID + "','score':'" + strconv.Itoa(item.Score) + "'},"
			item.Score = 0
		}
		room.User[l] = item
		if item.Status == "false" {
			Leave(room, item.OpenID)
		}
	}
	str = str + "]}"
	str = strings.Replace(str, "'", "\"", -1)
	ServerRoom(room, str)
	if len(room.User) != 6 {
		room.Status = true
	}
	UnderTime(5, room, "GameOverCountdown")
	ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'GameOverCountdownStop'}"))
	room.Word = ""
	UpdatePlayRoom(room)
	// UnReadyAll(room)
}

// 倒计时
func UnderTime(count int, room Room, mes string) {
	ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'"+mes+"'}"))
	for i := 0; i < count; i++ {
		ServerRoom(room, StrToJSON("time", "系统时间提示", "{'message':'房间公告: 倒计时还有"+strconv.Itoa(count-i)+"秒'}"))
		ServerRoom(room, StrToJSON("room", "倒计时", "{'message':'"+strconv.Itoa(count-i)+"'}"))
		time.Sleep(time.Second * 1)
	}
}

// 回合倒计时
func RoundTime(count int, room Room) bool {
	ServerRoom(room, StrToJSON("room", "房间状态", "{'message':'DrawCountdown'}"))
	b := len(room.User)
	for i := 0; i < count; i++ {
		ServerRoom(room, StrToJSON("time", "系统时间提示", "{'message':'房间公告: 倒计时还有"+strconv.Itoa(count-i)+"秒'}"))
		ServerRoom(room, StrToJSON("room", "倒计时", "{'message':'"+strconv.Itoa(count-i)+"'}"))
		time.Sleep(time.Second * 1)
		room = GetRoom(room)
		if len(room.User) < 2 {
			return true
		}
		if i == 0 {
			ServerRoom(room, StrToJSON("room", "答案提示", "{'message':'答案提示: "+GetWordMess("first", room.Word)+"'}"))
		}
		if i == 30 {
			ServerRoom(room, StrToJSON("room", "答案提示", "{'message':'答案提示: "+GetWordMess("second", room.Word)+"'}"))
		}
		if i == 60 {
			ServerRoom(room, StrToJSON("room", "答案提示", "{'message':'答案提示: "+GetWordMess("third", room.Word)+"'}"))
		}
		if room.GuessPeople == 0 && i+3 < count {
			count = i + 3
			room.GuessPeople = -1
			UpdatePlayRoom(room)
		}
		if b != room.GuessPeople && room.GuessPeople != -1 {
			count = count - 10
			b = room.GuessPeople
		}
	}
	return true
}

// 获取提示
func GetWordMess(num string, word string) string {
	w := Mydb.GuessAnswer{
		Answer: word,
	}
	A, has := Wordctrl.GetAnswerOne(w)
	str := ""
	if has {
		switch num {
		case "first":
			str = A.First
		case "second":
			str = A.Second
		case "third":
			str = A.Third
		}
	}
	return str
}

// 根据OpenID获取nickname
func GetNickName(ID string) string {
	user := Mydb.User{
		OpenID: ID,
	}
	u, has := Userctrl.GetUser(user)
	str := ""
	if has {
		str = u.NickName
	}
	return str
}
