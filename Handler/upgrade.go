/***************************
@File        : upgrade.go
@Time        : 2021/01/11 13:36:51
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : 用户升级增加金币操作
****************************/

package Handler

import (
	"draw/Mydb"
	"encoding/json"
	"log"
)

// 用户点赞操作
func Upgrade(mes []byte) string {
	// ctrlUser := Mydb.NewUserCtrl()
	// var user Mydb.User
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "点赞失败,数据无法解析")
	}
	User := Mydb.User{
		OpenID: user.OpenID,
	}
	thisUser, has := ctrlUser.GetUser(User)
	// log.Println(thisUser)
	if has {
		thisUser.Like = thisUser.Like + 1
		ctrlUser.Update(thisUser)
	} else {
		return ToMes("error", "点赞失败，找不到用户")
	}
	return UserToString("ok", thisUser, "点赞成功")
}

// 清空数据
func ClearData(mes []byte) string {
	err := json.Unmarshal(mes, &user)
	if err != nil {
		log.Println("数据问题:", err.Error())
		return ToMes("error", "清空失败,数据无法解析")
	}
	User := Mydb.User{
		OpenID: user.OpenID,
	}
	thisUser, has := ctrlUser.GetUser(User)
	// log.Println(thisUser)
	if has {
		thisUser.Number = 0
		thisUser.Score = 0
		thisUser.Like = 0
		ctrlUser.Update(thisUser)
	} else {
		return ToMes("error", "清空失败，找不到用户")
	}
	return UserToString("ok", thisUser, "清空成功")
}
