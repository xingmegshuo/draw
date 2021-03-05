/***************************
@File        : Answer.go
@Time        : 2021/03/05 16:13:21
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : db for answer
****************************/

package Mydb

import "log"

// struct
type GuessAnswer struct {
	Id     int64
	Answer string
	First  string
	Second string
	Third  string
}

// 插入单个谜底
func (G GuessAnswer) Insert(a ...interface{}) bool {
	_, err := orm.InsertOne(a[0])
	if err != nil {
		log.Panic(err)
	}
	return true
}

// 获取单个谜底
func (G GuessAnswer) GetAnswer(a ...interface{}) (GuessAnswer, bool) {
	u, ok := a[0].(GuessAnswer)
	if ok != false {
		has, _ := orm.Get(&u)
		return u, has
	} else {
		return GuessAnswer{}, false
	}
}
