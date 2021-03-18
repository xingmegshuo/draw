package main

import (
	"log"
	"strings"
)

func main() {
	b := "我要吃饺子"
	c := "饺子"
	for _, i := range c {
		b = strings.Replace(b, string(i), "*", -1)
		log.Println(i, b)
	}
	log.Println(b)
}
