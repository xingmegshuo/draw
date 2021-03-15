package main

import (
	"fmt"
	"time"
)

var a []int

func sum(c chan bool) {
	time.Sleep(time.Second * 4)
	c <- true
}
func main() {
	c := make(chan bool)
	a = []int{1, 2, 3, 4, 5}
	// x := false
	go sum(c)
	for i := range c {
		if i == true {
			a = []int{1, 2, 7}
			fmt.Println(a)
			break
		}
	}
	go console()
}

func console() {
	for _, i := range a {
		time.Sleep(time.Second * 1)
		fmt.Println(i)
	}
}
