package main

import (
	"fmt"
	"regexp"
	"os"
	"io/ioutil"
	"strings"
	"strconv"
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
const msgIDRegex string = "message-id=<([^>]+){1}>"

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
var messageID = regexp.MustCompile(msgIDRegex)

/**************************************************
 *
 * Data Structures to hold Important Information
 *
***************************************************/

// email struct to hold info about one email
type Email struct {
    queueID    			string       // message id processed by Postfix
    sender     			string 	     // address of the sender
    receiver   			string       // address of the receiver
    size       			int          // size of the email
    date       			string       // date and time email
   	client              []string       // client hostname and IP address for inbound email
    cleanup    			[]string       // store info when inbound email gets cleaned up 
    status     			string       // status of the processed email
    msgID      			string       // message id of the inbound email
  	emailType  			string       // outgoing email or incoming 
}
// email type vars
var outgoing string = "outgoing"
var incoming string = "incoming"

var emailList []Email  // global variable for slice of emails



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

/*
// Initialization of an Email struct
func (email Email) emailInit(){
	email.queueID = ""
	email.sender = "" 
	email.receiver = ""
	email.size = 0
	email.date = ""
	email.clientHostName = ""
	email.clientIP = ""
	email.cleanup = ""
	email.status = ""
	email.msgID = ""
	email.emailType = ""
}*/

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
   var mID string       // message id of the inbound email
   var emailIndex int = -1   // index of the email in the email list
   var income int = 0
   var out int = 0
   var total int = 0
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
   		mID = messageID.FindString(data[i])

   		if qID != "" {
   			// if the email is NOT in emailList 
   			if !hasEmailID(emailList, qID) {
   				// create a new email and initialize it
   			
   				clientSlice := make([]string, 2)
   				cleanupSlice := make([]string, 2) 

   				email := Email{queueID: qID,
   							   sender: "",
   							   receiver: "",
   							   size: 0,
   							   date: "",
   							   client: clientSlice,
   							   cleanup: cleanupSlice,
   							   status: "",
   							   msgID: "",
   							   emailType: ""}
   				emailList = append(emailList, email)
   			} else {
   				emailIndex = findEmailIndex(emailList, qID)
   			}
   		}
   		
   	
   		// the Host/IP Address of the client connected to the SMTP daemon
   		if smtpd != "" {
   			if client != "" {
   				fmt.Println(client)
   				startHost := strings.Index(client, "=") + 1
   				startIP := strings.Index(client, "[") + 1
   				endHost := startIP - 1
   				endIP := len(client) - 1
   				hostname := client[startHost: endHost]
   				IPaddress := client[startIP: endIP]
   				hostIP := []string{hostname, IPaddress}

   				if emailIndex != -1 {
   					emailList[emailIndex].client = append(emailList[emailIndex].client, hostIP[0], hostIP[1])
   	
   				}				
   			}
   		
   			if emailIndex != -1 {
   				emailList[emailIndex].emailType = incoming	
   			}
   		}

   		// the msg id of the currently processed email
   		if cleanup != "" {
   		//	fmt.Printf(data[i] + "\n")
   			if mID != "" {
   				if emailIndex != -1 {
   					emailList[emailIndex].msgID = mID
   					fmt.Println(mID)
   				}
   			}
   		}

   		// the time an email was removed from que or its size, sender and number of recipients
   		if qmgr != "" {
   			// deal with incoming email
   			if emailList[emailIndex].emailType == incoming {
   				if from != "" {
   					emailList[emailIndex].sender = from[6: (len(from) - 1)]
   				}
   				if size != "" {
   					msgSize, err := strconv.Atoi(size[5: len(size)])
   					if err == nil {
   						emailList[emailIndex].size = msgSize
   					}	
   				}
   				if date != "" {
   					emailList[emailIndex].date = date
   					fmt.Println(date + "---")
   				}
   			}
   			
   			
   		}

   		// detailed info about destination, delay, relay and status etc
   		if  smtp != "" {
   		//	fmt.Println(data[i])
   			if to != "" {
   				
   				fmt.Printf(to + "   ")
   				fmt.Printf(status[7:] + "\n")
   			}
   			
   		}
   }
   fmt.Println("Inbound mails num: ", income)
   fmt.Println("Outgoing emails num: ", out)
   fmt.Println("Total emails: ", total)
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