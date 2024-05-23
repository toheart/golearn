package Observe

/**
@file:
@author: levi.Tang
@time: 2024/10/9 16:05
@description: 观察者 接口
**/

type Observer interface {
	update(string)
	getID() string
}
