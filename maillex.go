package main

import (
	"fmt"
	"regexp"
	"os"
	"io/ioutil"
	"strings"
)

const usageMsg string = "usage: gomaillex maillog[filename] -option\n" + 
						"-option: -s [summary]\n"

// Constants of Regular Expression Patterns
const smtpRegex string = "postfix/smtp\\[\\d+\\]"
const smtpdRegex string = "postfix/smtpd\\[\\d+\\]"

var postfixsmtp = regexp.MustCompile(smtpRegex)
var postfixsmtpd = regexp.MustCompile(smtpdRegex)

func usage() {
	fmt.Printf(usageMsg)
	os.Exit(1)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Read data from the file 
// Split file data by newline
// Return an array of lines of strings
func openFile(file string) []string {
	dat, err := ioutil.ReadFile(file)
    check(err)
//    fmt.Print(string(dat))
    fmt.Print("\n")

    var data []string = strings.Split(string(dat), "\n")
    fmt.Printf("num of lines: %d\n", len(data))
    return data
}

// Parse the selected log file 
// And obtain the information needed
func parse(data []string) {
   fmt.Printf("Start parsing....\n")
   for i := 0; i < len(data); i++ {
   	/*	var line []string = strings.Split(data[i], " ")
   		for j := 0; j < len(line); j++ {

   		}
   	*/	
   		var smtp string = postfixsmtpd.FindString(data[i])
   		if  smtp != "" {
   			fmt.Printf(data[i] + "\n")
   		}
   }

}

func main() {
	args := os.Args
	fmt.Printf("\n")

	if len(args) < 2 {
		usage()
	}
	var file string = args[1]
	var data []string = openFile(file)
	parse(data)

	for i := 0; i < len(args); i++ {
		fmt.Printf("%s\n", args[i])
	}
}