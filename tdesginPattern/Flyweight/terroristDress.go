package Flyweight

/**
@file:
@author: levi.Tang
@time: 2024/10/8 11:44
@description:
**/

type TerroristDress struct {
	color string
}

func (t *TerroristDress) getColor() string {
	return t.color
}

func newTerroristDress() *TerroristDress {
	return &TerroristDress{color: "red"}
}
