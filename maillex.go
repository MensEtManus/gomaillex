package main

import (
	"fmt"
//	"regexp"
	"os"
	"io/ioutil"
	"strings"
)

const usageMsg string = "usage: gomaillex maillog[filename] -option\n" + 
						"-option: -s [summary]\n"

func usage() {
	fmt.Printf(usageMsg)
	os.Exit(1)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Parse the selected log file 
// And obtain the information needed
func parser(file string) {
    dat, err := ioutil.ReadFile(file)
    check(err)
 //   fmt.Print(string(dat))
    fmt.Print("\n")

    var data []string = strings.Split(string(dat), "\n")
    fmt.Printf("num of lines: %d\n", len(data))
}

func main() {
	args := os.Args
	fmt.Printf("\n")

	if len(args) < 2 {
		usage()
	}
	var file string = args[1]
	parser(file)

	for i := 0; i < len(args); i++ {
		fmt.Printf("%s\n", args[i])
	}
}