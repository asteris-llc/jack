package main

import (
	"github.com/BrianHicks/jack"
)

type Adder struct{}

func (_ *Adder) Calculate(a, b float64) float64 {
	return a + b
}

func (_ *Adder) Usage() string {
	return "add b to a"
}

func main() {
	server := jack.NewServer(&Adder{})
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
