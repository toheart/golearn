package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sync"
)

/**
@file:
@author: levi.Tang
@time: 2024/6/18 14:55
@description:
**/

func main() {
	key, value := "test", "test"
	c, err := redis.Dial("tcp", "175.178.49.104:6379", redis.DialPassword("ningzaichun"))
	if err != nil {
		panic(err)
	}
	defer c.Close()
	timeoutMs := 60 * 1000
	isset, err := redis.Int(c.Do("SETNX", key, value))
	if err != nil {
		fmt.Errorf("TryLock SETNX failed, err[%s]", err)
		return
	}
	var system sync.Mutex
	system.Lock()
	if isset == 1 {
		if timeoutMs > 0 {
			c.Do("PEXPIRE", key, timeoutMs)
		}
		return
	}

	if timeoutMs > 0 {
		// if key exists but ttl not set, reset ttl
		if ttl, err := redis.Int64(c.Do("PTTL", key)); err == nil && ttl == -1 {
			c.Do("PEXPIRE", key, timeoutMs)
		}
	}
}
