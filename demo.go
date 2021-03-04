/***************************
@File        : demo.go
@Time        : 2021/03/04 17:51:07
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : str for
****************************/

package main

import "fmt"

func main() {
	for _, i := range "老虎" {
		fmt.Println(i)
		if i == '老' {
			fmt.Println("真不错")
		}
	}
}
