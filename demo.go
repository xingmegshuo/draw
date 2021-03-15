/***************************
@File        : demo.go
@Time        : 2021/03/04 17:51:07
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : str for
****************************/

package main

import (
	"log"
	"time"
)

var a []int

func main() {
	a = []int{1, 2, 3, 4}
	// log.Println(a)
	go change()
	consloe()
}

// 循环输出
func consloe() {
	for _, i := range a {
		log.Println(i)
	}
}

// 输出中修改值
func change() {
	for {
		for l, b := range a {
			if b == 2 {
				a[l] = 5
				log.Println("修改")
				// log.Println(a)
				time.Sleep(time.Second * 1)
				break
			}
		}
	}
}
