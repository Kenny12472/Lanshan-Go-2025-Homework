package main

import (
	"fmt"
	"time"

	"sixth/utils"
)

func main() {
	token, expire, err := utils.MakeToken("testuser", time.Now().Add(10*time.Minute))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("access token:", token)
	fmt.Println("expire:", expire)

	username, exp, err := utils.ParseToken(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("username:", username)
	fmt.Println("token expire:", exp)
}
