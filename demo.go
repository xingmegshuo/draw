/***************************
@File        : demo.go
@Time        : 2021/03/04 17:51:07
@AUTHOR      : small_ant
@Email       : xms.chnb@gmail.com
@Desc        : str for
****************************/

package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
)

func main() {
	for i := 0; i < 3; i++ {
		result, _ := rand.Int(rand.Reader, big.NewInt(20))
		// result = int(result)
		fmt.Println(reflect.TypeOf(int(result.Int64())), result)
	}
}
