package Bridge

import (
	"fmt"
	"testing"
)

/**
@file:
@author: levi.Tang
@time: 2024/10/3 17:34
@description:
先关注谁是适用方：
* 抽象层
* 实施层
**/

type Computer interface {
	Print()
	SetPrinter(Printer)
}

type Printer interface {
	PrintFile()
}

// 精确抽象层
// mac
type Mac struct {
	printer Printer
}

func (m *Mac) Print() {
	fmt.Println("print request for mac")
	m.printer.PrintFile()
}

func (m *Mac) SetPrinter(printer Printer) {
	m.printer = printer
}

// windows
type Windows struct {
	printer Printer
}

func (w *Windows) Print() {
	fmt.Println("print request for mac")
	w.printer.PrintFile()
}

func (w *Windows) SetPrinter(printer Printer) {
	w.printer = printer
}

// 实施层, 具体干活的对象

type Epson struct {
}

func (p *Epson) PrintFile() {
	fmt.Println("Printing by a EPSON Printer")
}

type Hp struct {
}

func (p *Hp) PrintFile() {
	fmt.Println("Printing by a HP Printer")
}

func Test_Printer(t *testing.T) {

	hpPrinter := &Hp{}
	epsonPrinter := &Epson{}

	macComputer := &Mac{}

	macComputer.SetPrinter(hpPrinter)
	macComputer.Print()
	fmt.Println()

	macComputer.SetPrinter(epsonPrinter)
	macComputer.Print()
	fmt.Println()

	winComputer := &Windows{}

	winComputer.SetPrinter(hpPrinter)
	winComputer.Print()
	fmt.Println()

	winComputer.SetPrinter(epsonPrinter)
	winComputer.Print()
	fmt.Println()
}
