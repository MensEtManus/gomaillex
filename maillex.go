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
const reasonRegex string = "status=.*?(.*)"

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
var statusReason = regexp.MustCompile(reasonRegex)

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
   	client              []string     // client hostname and IP address for inbound email
    cleanup    			[]string     // store info when inbound email gets cleaned up 
    status     			string       // status of the processed email
    reason              string       // the explanation of the delivery status
    msgID      			string       // message id of the inbound email
  	emailType  			string       // outgoing email or incoming 
}
// email type vars
var outgoing string = "outgoing"
var incoming string = "incoming"

var emailIn  []Email  // global variable for slice of incoming emails
var emailOut []Email  // global variable for slice of outgoing emails



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
	var index int = -1
	for i, email := range list {
		if email.queueID == id {
			index = i
		}
	}
	return index 
}

// print info in Email 
func printEmail(emails []Email) {
	fmt.Println("                           Email List Info")
	fmt.Println("----------------------------------------------------------------------------------------")
	for i, email := range emails {
		fmt.Println(i)
		fmt.Println("From: " + email.sender)
		fmt.Println("To: " + email.receiver)
		fmt.Print("Size: ")
		fmt.Println(email.size)
		fmt.Println("Date: " + email.date)
		fmt.Println("Status: " + email.status)
		fmt.Println("Reason: " + email.reason)
		fmt.Println("Host Name: " + email.client[0])
		fmt.Println("IP address: " + email.client[1])
		fmt.Println("In/Out: " + email.emailType)
		fmt.Println()
	}	
}

// Append new email into Incoming/Outgoing email List
func addEmail(emailList []Email, qID string) {
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
						reason: "",
						msgID: "",
						emailType: ""}
		emailList = append(emailList, email)

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

   var smtp, smtpd, cleanup, qmgr string 
   var from, to, size, status, date string 
   var client string    // Hostname and IP address of the clients connected to the SMTP daemon
   var qID string       // queue ID of each email 
   var mID string       // message id of the inbound email
   var inEmailInx int = -1  // index of the email in the incoming email list
   var outEmailInx int = -1 // index of the email int the outgoing email list
   var statusRsn string     // reason of the email delivery status

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
   		statusRsn = statusReason.FindString(data[i])
   		
   		// find the email index in both incoming and outgoing email lists
   		if qID != "" {
   			inEmailInx = findEmailIndex(emailIn, qID)
   			outEmailInx = findEmailIndex(emailOut, qID)
   		}

   		// the Host/IP Address of the client connected to the SMTP daemon
   		// inbound email
   		if smtpd != "" {
   			
   			if qID != "" {
   			//	addEmail(emailIn, qID)
   				if !hasEmailID(emailIn, qID) {
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
   						emailIn = append(emailIn, email)

   					} 
   				inEmailInx = findEmailIndex(emailIn, qID)
   				
   			}
   			if client != "" {
   				
   				startHost := strings.Index(client, "=") + 1
   				startIP := strings.Index(client, "[") + 1
   				endHost := startIP - 1
   				endIP := len(client) - 1
   				hostname := client[startHost: endHost]
   				IPaddress := client[startIP: endIP]
   				hostIP := []string{hostname, IPaddress}

   				if inEmailInx != -1 {
   					emailIn[inEmailInx].client[0] = hostIP[0]
   					emailIn[inEmailInx].client[1] = hostIP[1]
   					emailIn[inEmailInx].emailType = incoming	
   				}				
   			}
   		}

   		// the msg id of the currently processed email
   		if cleanup != "" {
   			if mID != "" {
   				if inEmailInx != -1 {
   					emailIn[inEmailInx].msgID = mID
   			
   				}
   			}
   		}

   		// the time an email was removed from que or its size, sender and number of recipients
   		if qmgr != "" {
   			// deal with incoming email
   			if inEmailInx != -1 {
   				if from != "" {
   					emailIn[inEmailInx].sender = from[6: (len(from) - 1)]
   				}
   				if size != "" {
   					msgSize, err := strconv.Atoi(size[5: len(size)])
   					if err == nil {
   						emailIn[inEmailInx].size = msgSize
   					}	
   				}
   				if date != "" {
   					emailIn[inEmailInx].date = date
   				}
   			} 
   			
   			// deal with outgoing email queue manager
   			if outEmailInx != -1 {
   				if from != "" {
   					emailOut[outEmailInx].sender = from[6: (len(from) - 1)]
   				}
   				if size != "" {
   					msgSize, err := strconv.Atoi(size[5: len(size)])
   					if err == nil {
   						emailOut[outEmailInx].size = msgSize
   					}	
   				}
   				if date != "" {
   					emailOut[outEmailInx].date = date
   				}
   			}
   			
   		}

   		// detailed info about destination, delay, relay and status etc
   		if  smtp != "" {
   			// adding email to the outgoing email list
   			if qID != "" {
   				if !hasEmailID(emailOut, qID) {
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
   						emailOut = append(emailOut, email)

   				} 
   				outEmailInx = findEmailIndex(emailOut, qID)
   				emailOut[outEmailInx].emailType = outgoing
   			}

   			if to != "" {
   				var endInx = strings.Index(to, ">")
   				var receiver = to[4: endInx]
   				var rsnStart, rsnEnd int
   				var rsn string 
   				if statusRsn != "" {
   					rsnStart = strings.Index(statusRsn, "(") + 1
   					rsnEnd = len(statusRsn) - 1
   					rsn = statusRsn[rsnStart: rsnEnd]
   				}
   				if outEmailInx != -1 {
   					emailOut[outEmailInx].receiver = receiver
   					emailOut[outEmailInx].status = status[7:]
   					emailOut[outEmailInx].reason = rsn
   				}
   								
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
    printEmail(emailOut)

}