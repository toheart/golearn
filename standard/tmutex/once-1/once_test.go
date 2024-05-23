package once_1

import (
	"fmt"
	"sync"
	"testing"
)

/**
@file:
@author: levi.Tang
@time: 2024/9/3 9:19
@description:
**/

var doOnce sync.Once
var doSomething *DoSomething

type DoSomething struct {
}

func NewDoSomething() *DoSomething {
	doOnce.Do(func() {
		fmt.Println("init doSomething..")
		doSomething = &DoSomething{}
	})
	fmt.Println("outer doOnce..")
	return doSomething
}

func TestSimpleOnce(t *testing.T) {
	first := NewDoSomething()

	second := NewDoSomething()

	fmt.Printf("first:%p second:%p \n", first, second)
}
