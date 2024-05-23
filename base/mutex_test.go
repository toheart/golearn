package base

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
*
@file:
@author: levi.Tang
@time: 2024/7/30 12:01
@description:
*
*/
type Donation struct {
	cond    *sync.Cond
	balance int
}

func TestCond(t *testing.T) {
	donation := &Donation{
		cond: sync.NewCond(&sync.Mutex{}),
	}

	// Listener goroutines
	f := func(goal int) {
		donation.cond.L.Lock()
		for donation.balance < goal {
			donation.cond.Wait()
			fmt.Printf("goal[%d] weak up \n", goal)
		}
		fmt.Printf("%d$ goal reached\n", donation.balance)
		donation.cond.L.Unlock()
	}
	go f(10)
	go f(15)
	// Updater goroutine
	go func() {
		for {
			time.Sleep(time.Second)
			donation.cond.L.Lock()
			donation.balance++
			fmt.Printf("main: %d \n", donation.balance)
			donation.cond.L.Unlock()
			donation.cond.Broadcast()
		}
	}()
	select {}
}

// 没有消费者, 也正常运行, 所以可以在vm内部初始化个cond
func TestCondWithoutConsumer(t *testing.T) {
	donation := &Donation{
		cond: sync.NewCond(&sync.Mutex{}),
	}

	// Updater goroutine
	go func() {
		for {
			time.Sleep(time.Second)
			donation.cond.L.Lock()
			donation.balance++
			fmt.Printf("main: %d \n", donation.balance)
			donation.cond.L.Unlock()
			donation.cond.Broadcast()
		}
	}()
	select {}
}

// 没有消费者, 也正常运行, 所以可以在vm内部初始化个cond
func TestCondLock(t *testing.T) {
	donation := &Donation{
		cond: sync.NewCond(&sync.Mutex{}),
	}
	// Listener goroutines
	f := func(goal int) {
		donation.cond.L.Lock()
		for donation.balance < goal {
			donation.cond.Wait()
			fmt.Printf("goal[%d] weak up \n", goal)
		}
		donation.cond.L.Unlock()
		fmt.Printf("%d$ goal reached\n", donation.balance)
	}
	go f(10)
	// Updater goroutine
	go func() {
		for {
			time.Sleep(time.Second)
			donation.cond.L.Lock()
			donation.balance++
			fmt.Printf("main: %d \n", donation.balance)
			donation.cond.L.Unlock()
			donation.cond.Broadcast()
		}
	}()
	select {}
}
