/***************************
@File        : login.go
@Time        : 2020/12/21 14:33:25
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 解析登录
****************************/

package Handler

import (
	"draw/Mydb"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

var ctrlUser = Mydb.NewUserCtrl()

// var ctrlBack = Mydb.NewBackCtrl()
var ctrlBuddy = Mydb.NewBuddyCtrl()
var ctrlRecord = Mydb.NewRecordCtrl()
var user Mydb.User
var backpack Mydb.Backpack
var record Mydb.Record
var buddy Mydb.Buddy

// 用户登录函数处理
func Login(mes []byte, ws *websocket.Conn) string {
	// ctrlUser := Mydb.NewUserCtrl()
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "登录操作失败")
	}
	// log.Println("登录发送的数据", user)
	if len(user.OpenID) > 0 && len(user.NickName) > 0 && len(user.AvatarURL) > 0 {
		// log.Println("hhhh")
		thisUser := Mydb.User{
			OpenID: user.OpenID,
		}
		U, has := ctrlUser.GetUser(thisUser)
		// log.Println("aaaa")
		if has {
			mes := UserToString("ok", U, "登录成功")
			// log.Println("bbbb")
			for k, v := range client_user {
				if v == U.OpenID {
					delete(client_user, k)
				}
			}
			client_user[ws] = U.OpenID
			return mes
		} else {
			user.Score = 0
			user.Like = 0
			user.Number = 0
			user.Id = 0
			// log.Println(user)
			ctrlUser.Insert(user)
			// InitBack(user)
			u, _ := ctrlUser.GetUser(user)
			client_user[ws] = u.OpenID
			mes := UserToString("ok", u, "登录成功")
			// log.Println("cccc")
			return mes
		}
	} else {
		return ToMes("error", "请检查发送的数据是否完整")
	}
}

// 转换内容
func UserToString(status string, user Mydb.User, mes string) string {
	str := "{'status':'" + status + "','mes':'" + mes + "','data':{'openID':'" + user.OpenID + "','nickName':'" + user.NickName + "','avatarUrl':'" + user.AvatarURL + "','Number':'" + strconv.Itoa(user.Number) + "','Score':'" + strconv.Itoa(user.Score) + "','like':'" + strconv.Itoa(user.Like) + "','orther':'" + user.Orther + "','id':'" + strconv.Itoa(int(user.Id)) + "'}}"
	str = strings.Replace(str, "'", "\"", -1)
	return str
}

// 不携带数据
func ToMes(status string, mes string) string {
	str := "{'status':'" + status + "','mes':'" + mes + "'}"
	str = strings.Replace(str, "'", "\"", -1)
	return str
}

// 初始化个人仓库的内容
// func InitBack(user Mydb.User) {
// 	user, has := ctrlUser.GetUser(user)
// 	if has {

// 		back := Mydb.Backpack{
// 			Name:     "基础帽子",
// 			Property: 0,
// 			User:     int(user.Id),
// 		}
// 		back1 := Mydb.Backpack{
// 			Name:     "基础下装",
// 			Property: 0,
// 			User:     int(user.Id),
// 		}
// 		back2 := Mydb.Backpack{
// 			Name:     "基础下装",
// 			Property: 0,
// 			User:     int(user.Id),
// 		}
// 		ctrlBack.Insert(back1)
// 		ctrlBack.Insert(back2)
// 		ctrlBack.Insert(back)
// 	}

// }

// 返回用户
func GetUserMes(mes []byte) string {
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "获取信息失败")
	}
	thisUser := Mydb.User{
		OpenID: user.OpenID,
	}
	U, _ := ctrlUser.GetUser(thisUser)
	str := UserToString("ok", U, "获取成功")
	return str
}
