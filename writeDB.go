/***************************
@File        : writeDB.go
@Time        : 2021/03/05 15:13:25
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : write db
****************************/

package main

import (
	"draw/Mydb"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	ctrl := Mydb.NewGuessAnswerCtrl()
	f, err := excelize.OpenFile("thesaurus.xlsx")
	if err != nil {
		println(err.Error())
		return
	}
	// 获取工作表中指定单元格的值
	// 获取 Sheet1 上所有单元格
	rows := f.GetRows("Sheet1")
	for l, row := range rows {
		if l == 0 {
			continue
		}
		// fmt.Println(row[1], row[2], row[3], row[4])
		answer := Mydb.GuessAnswer{
			Answer: row[1],
			First:  row[2],
			Second: row[3],
			Third:  row[4],
		}
		ctrl.Insert(answer)
	}
}
