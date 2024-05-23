package Flyweight

/**
@file:
@author: levi.Tang
@time: 2024/10/8 11:45
@description:
**/

type CounterTerroristDress struct {
	color string
}

func (c *CounterTerroristDress) getColor() string {
	return c.color
}

func newCounterTerroristDress() *CounterTerroristDress {
	return &CounterTerroristDress{color: "blue"}
}
