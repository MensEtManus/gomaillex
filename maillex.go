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
const fromRegex string = "from=<([^>]+){1}>"
const toRex string = "to=<([^>]+){1}>"
const sizeRegex string = "size=([0-9]+){1}"
const statusRegex string = "status=([^ ]+){1}"
const dateRegex string = "^(\\w{3}[^a-zA-Z]+)"
const idRegex string = "[A-F0-9]{11}"
const clientRegex string = "client=.*?\\[([0-9.]+)+\\]"

// Regular Expression Variables
var postfixsmtp = regexp.MustCompile(smtpRegex)
var postfixsmtpd = regexp.MustCompile(smtpdRegex)
var postfixcleanup = regexp.MustCompile(cleanupRegex)
var postfixqmgr = regexp.MustCompile(qmgrRegex)
var sender = regexp.MustCompile(fromRegex)
var receiver = regexp.MustCompile(toRex)
var mailSize = regexp.MustCompile(sizeRegex)
var sendStatus = regexp.MustCompile(statusRegex) 
var dateInfo = regexp.MustCompile(dateRegex)
var postfixID = regexp.MustCompile(idRegex)
var clientHostIP = regexp.MustCompile(clientRegex)

/**************************************************
 *
 * Data Structures to hold Important Information
 *
***************************************************/

// email struct to hold info about one email
type Email struct {
    queueID string    // message id processed by Postfix
    sender string 	  // address of the sender
    receiver string   // address of the receiver
    size int          // size of the email
    date string       // date and time 
  	emailType string  // outgoing email or incoming email
}
// email type vars
var outgoing string = "outgoing"
var incoming string = "incoming"

var emailList []Email    // slice of email structs


func usage() {
	fmt.Printf(usageMsg)
	os.Exit(1)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// check if the queue ID already exists in the email list
func hasEmailID(list []Email, id string) bool {
	for _, email := range list {
		if email.queueID == id {
			return true
		}
	}
	return false
}

// find the index of an email in the email list
func findEmailIndex(list []Email, id string) int {
	var index int
	for i, email := range list {
		if email.queueID == id {
			index = i
		}
	}
	return index 
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

   var smtp, smtpd, cleanup, qmgr string
   var from, to, size, status, date string
   var client string    // Hostname and IP address of the clients connected to the SMTP daemon
   var qID string       // queue ID of each email 
   var emailIndex int    // index of the email in the email list

   // Loop through all the lines to obtain info needed
   for i := 0; i < len(data); i++ {

   		// find matching strings from the current line 
   		smtp = postfixsmtp.FindString(data[i])
   		smtpd = postfixsmtpd.FindString(data[i])
   		cleanup = postfixcleanup.FindString(data[i])
   		qmgr = postfixqmgr.FindString(data[i])
   		from = sender.FindString(data[i])
   		to = receiver.FindString(data[i])
   		size = mailSize.FindString(data[i])
   		status = sendStatus.FindString(data[i])
   		date = dateInfo.FindString(data[i])
   		client = clientHostIP.FindString(data[i])
   		qID = postfixID.FindString(data[i])


   		if qID != "" {
   			// if the email is NOT in emailList 
   			if !hasEmailID(emailList, qID) {
   				emailList = append(emailList, Email{queueID: qID})
   			} else {
   				emailIndex = findEmailIndex(emailList, qID)
   			}
   		}
   		
   	
   		// the Host/IP Address of the client connected to the SMTP daemon
   		if smtpd != "" {
   			if client != "" {
   				islocal, _ := regexp.MatchString("client=localhost.*", client)
   				if islocal {
   					emailList[emailIndex].emailType = outgoing
   				} else {
   					emailList[emailIndex].emailType = incoming
   				}
   				
   			}
   			

   		}

   		// the msg id of the currently processed email
   		if cleanup != "" {
   		//	fmt.Printf(data[i] + "\n")
   		}

   		// the time an email was removed from que or its size, sender and number of recipients
   		if qmgr != "" {
   		//	fmt.Printf("\n")
   			if from != "" {
   				fmt.Printf(from + "   ")
   				fmt.Printf(size + "\n")
   			}
   		}

   		// detailed info about destination, delay, relay and status etc
   		if  smtp != "" {

   			if to != "" {
   				fmt.Printf(date + "---")
   				fmt.Printf(to + "   ")
   				fmt.Printf(status[7:] + "\n")
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