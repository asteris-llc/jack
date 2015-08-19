package main

import (
	"errors"
	"fmt"
	"github.com/BrianHicks/jack"
	"os"
	"strconv"
)

// ---- wrapper ----

type Math struct {
	*jack.Client
}

func NewMath(plugin string) (*Math, error) {
	command, err := jack.Load("math-" + plugin)
	if err != nil {
		return nil, err
	}

	err = command.Start()
	if err != nil {
		return nil, err
	}

	return &Math{command}, nil
}

func (c *Math) Calculate(a, b float64) (float64, error) {
	results, err := c.Call("Calculate", a, b)
	if err != nil {
		return 0, err
	}

	result, ok := results.([]interface{})[0].(float64)
	if !ok {
		return 0, errors.New("bad result from Calculate")
	}

	return result, nil
}

func (c *Math) Usage() (string, error) {
	results, err := c.Call("Usage")
	if err != nil {
		return "", err
	}

	result, ok := results.([]interface{})[0].(string)
	if !ok {
		return "", errors.New("bad result from Usage")
	}

	return result, nil
}

// ------ CLI ------

func usage() {
	fmt.Println("usage: math {cmd} [args...]")
	fmt.Println("try `math add 1 2` or `math help add`")
	os.Exit(2)
}

func help(plugin string) {
	cmd, err := NewMath(plugin)
	if err != nil {
		panic(err)
	}
	defer cmd.Stop()

	output, err := cmd.Usage()
	if err != nil {
		panic(err)
	}

	fmt.Printf("help for %s:\n%s\n", plugin, output)
	os.Exit(2)
}

func calculate(plugin string, a, b float64) {
	cmd, err := NewMath(plugin)
	if err != nil {
		panic(err)
	}
	defer cmd.Stop()

	result, err := cmd.Calculate(a, b)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", result)
	os.Exit(0)
}

func main() {
	operation := os.Args[1]
	args := os.Args[2:]

	switch operation {
	case "help":
		help(args[0])
	case "usage":
		usage()
	default:
		a, _ := strconv.ParseFloat(args[0], 64)
		b, _ := strconv.ParseFloat(args[1], 64)
		calculate(operation, a, b)
	}

	fmt.Println("test")
}
