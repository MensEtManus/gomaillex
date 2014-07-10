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
const cleanupRegex string = "postfix/cleanup\\[\\d+\\]"
const qmgrRegex string = "postfix/qmgr\\[\\d+\\]"
const sendmailRegex string = "sendmail\\[\\d+\\]"
const fromRegex string = "from=<([^>]+){1}>"
const toRex string = "to=<([^>]+){1}>"
const sizeRegex string = "size=([0-9]+){1}"
const statusRegex string = "status=([^ ]+){1}"
const dateRegex string = "^(\\w{3}[^a-zA-Z]+)"
const sendmailIDRegex string = "([a-f0-9]{14}){1})$"
const idRegex string = "([a-f0-9]{11}){1})$"

// Regular Expression Variables
var postfixsmtp = regexp.MustCompile(smtpRegex)
var postfixsmtpd = regexp.MustCompile(smtpdRegex)
var postfixcleanup = regexp.MustCompile(cleanupRegex)
var postfixqmgr = regexp.MustCompile(qmgrRegex)
var sendmail = regexp.MustCompile(sendmailRegex)
var sender = regexp.MustCompile(fromRegex)
var receiver = regexp.MustCompile(toRex)
var mailSize = regexp.MustCompile(sizeRegex)
var sendStatus = regexp.MustCompile(statusRegex) 
var dateInfo = regexp.MustCompile(dateRegex)
var sendmailID = regexp.MustCompile(sendmailIDRegex)
var postfixID = regexp.MustCompile(idRegex)

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
    fmt.Print("\n")

    var data []string = strings.Split(string(dat), "\n")
    fmt.Printf("num of lines: %d\n", len(data))
    return data
}

// Parse the selected log file 
// And obtain the information needed
func parse(data []string) {
   fmt.Printf("Start parsing....\n")
   var smtp, smtpd, cleanup, qmgr, sdmData string
   var from, to, size, status, date string
   for i := 0; i < len(data); i++ {
   	/*	var line []string = strings.Split(data[i], " ")
   		for j := 0; j < len(line); j++ {

   		}
   	*/	
   		// find matching strings from the  line 
   		smtp = postfixsmtp.FindString(data[i])
   		smtpd = postfixsmtpd.FindString(data[i])
   		cleanup = postfixcleanup.FindString(data[i])
   		qmgr = postfixqmgr.FindString(data[i])
   		sdmData = sendmail.FindString(data[i]) 
   		from = sender.FindString(data[i])
   		to = receiver.FindString(data[i])
   		size = mailSize.FindString(data[i])
   		status = sendStatus.FindString(data[i])
   		date = dateInfo.FindString(data[i])

   		// detailed info about destination, delay, relay and status etc
   		if  smtp != "" {
   		//	fmt.Printf("\n")
   		//	fmt.Printf(data[i] + "\n")
   		}
   		// the Host/IP Address of the client connected to the SMTP daemon
   		if smtpd != "" {
   		//	fmt.Printf("\n")
   		}
   		// the msg id of the currently processed email
   		if cleanup != "" {
   		//	fmt.Printf(data[i] + "\n")
   		}
   		// the time an email was removed from que or its size, sender and number of recipients
   		if qmgr != "" {
   		//	fmt.Printf("\n")
   		}
   		// sendmail service data
   		if sdmData != "" {
   			fmt.Printf("sendmail   ")
   			fmt.Printf(date + "---")
   			if from != "" {
   				fmt.Printf(from + " ")
   				fmt.Printf(size + "\n")
   			}
   			if to != "" {
   				fmt.Printf(to + " ")
   				fmt.Printf(status + "\n")
   			}
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