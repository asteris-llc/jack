package main

import (
	"github.com/BrianHicks/jack"
)

type Multiplier struct{}

func (_ *Multiplier) Calculate(a, b float64) float64 {
	return a * b
}

func (_ *Multiplier) Usage() string {
	return "multiply a by b"
}

func main() {
	server := jack.NewServer(&Multiplier{})
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
