package main

import (
	"github.com/BrianHicks/jack"
)

type Subtracter struct{}

func (_ *Subtracter) Calculate(a, b float64) float64 {
	return a - b
}

func (_ *Subtracter) Usage() string {
	return "subtract b from a"
}

func main() {
	server := jack.NewServer(&Subtracter{})
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
