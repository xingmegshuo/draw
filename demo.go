/***************************
@File        : demo.go
@Time        : 2021/03/04 17:51:07
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : str for
****************************/

package main

import (
	"draw/Mydb"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ctrl := Mydb.NewGuessAnswerCtrl()
	SearchWord := Mydb.GuessAnswer{
		Id: 0,
	}
	words := ctrl.GetAnswer(SearchWord)

	rand.Seed(time.Now().Unix())
	j := rand.Intn(len(words))
	fmt.Println(words[j].Answer)
}
