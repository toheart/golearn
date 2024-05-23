package main

import "time"

/*
*
@file:
@author: levi.Tang
@time: 2024/4/8 15:06
@description:
*
*/
var sum int

func main() {

	go func() {
		time.Sleep(time.Second)
		panic("there are panic")
	}()

	select {}
}
