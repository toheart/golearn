package Lttlcache

import (
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"testing"
	"time"
)

/**
@file:
@author: levi.Tang
@time: 2024/9/19 17:43
@description:
**/

func Test_TTLCache(t *testing.T) {
	cache := ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](5 * time.Second),
	)
	cache.Set("first", "value1", ttlcache.DefaultTTL)
	cache.Set("second", "value2", ttlcache.NoTTL)
	cache.Set("third", "value1", ttlcache.DefaultTTL)
	go cache.Start()
	go func() {
		for {
			fmt.Println("found len:", cache.Len(), cache.Get("first").Value(), cache.Get("first").IsExpired())

			time.Sleep(2 * time.Second)
		}
	}()

	time.Sleep(2 * time.Minute)
}
