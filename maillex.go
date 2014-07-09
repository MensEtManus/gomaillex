package main

import (
	"fmt"
//	"regexp"
	"os"
	
)

const usageMsg string = "usage: gomaillex maillog[filename] -option\n"

func usage() {
	fmt.Printf(usageMsg)

	os.Exit(1)
}

func main() {
	args := os.Args
	fmt.Printf("\n")
	if len(args) < 2 {
		usage()
	}
	for i := 0; i < len(args); i++ {
		fmt.Printf("%s\n", args[i])
	}
}