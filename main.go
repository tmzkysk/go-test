package main

import (
	"fmt"
	"os"

	"github.com/tmzkysk/go-test/calc"
	"github.com/tmzkysk/go-test/somefunc"
)

func main() {
	os.Exit(run())
}

func run() int {

	fmt.Println(calc.Add(1, 2))
	c := somefunc.Client{&somefunc.ExampleCaller{}}
	fmt.Println(c.Run(5))
	return 0
}
