package tqueue

import (
	"fmt"
	"sync"
	"testing"
)

/**
@file:
@author: levi.Tang
@time: 2024/9/28 11:42
@description:
**/

type customQueue struct {
	queue []string
	lock  sync.RWMutex
}

func (c *customQueue) Enqueue(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.queue = append(c.queue, name)
}

func (c *customQueue) Dequeue() error {
	if len(c.queue) > 0 {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.queue = c.queue[1:]
		return nil
	}
	return fmt.Errorf("Pop Error: Queue is empty")
}

func (c *customQueue) Front() (string, error) {
	if len(c.queue) > 0 {
		c.lock.Lock()
		defer c.lock.Unlock()
		return c.queue[0], nil
	}
	return "", fmt.Errorf("Peep Error: Queue is empty")
}

func (c *customQueue) Size() int {
	return len(c.queue)
}

func (c *customQueue) Empty() bool {
	return len(c.queue) == 0
}

func TestDequeue(t *testing.T) {
	c := &customQueue{
		queue: make([]string, 0),
	}
	fmt.Printf("Enqueue: A\n")
	c.Enqueue("A")
	if err := c.Dequeue(); err != nil {
		t.Error(err)
	}
	fmt.Println(c.queue)
}
