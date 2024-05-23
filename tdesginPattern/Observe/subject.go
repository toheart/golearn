package Observe

/**
@file:
@author: levi.Tang
@time: 2024/10/9 16:05
@description:
**/

type Subject interface {
	register(observer Observer)
	deregister(observer Observer)

	notifyAll()
}
