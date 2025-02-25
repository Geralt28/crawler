package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, World!")
	args := os.Args[1:]
	dl := len(args)
	if dl == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if dl > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	base_url := args[0]
	fmt.Println("starting crawl of:", base_url)
}
