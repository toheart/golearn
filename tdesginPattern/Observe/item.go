package Observe

/**
@file:
@author: levi.Tang
@time: 2024/10/9 16:07
@description:  被观察主体
**/

type Item struct {
	observerList []Observer
	name         string
	instock      bool
}

func newItem(name string) *Item {
	return nil
}
