package main

import (
	"log"
	"time"
)

func main() {
	a := 30
	for i := 0; i < a; i++ {
		time.Sleep(time.Second * 1)
		log.Println(i)
		if i == 20 {
			a = 25
		}
	}

}
