package main

import (
	"github.com/BrianHicks/jack"
)

type Divider struct{}

func (_ *Divider) Calculate(a, b float64) float64 {
	return a / b
}

func (_ *Divider) Usage() string {
	return "divide a by b"
}

func main() {
	server := jack.NewServer(&Divider{})
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
