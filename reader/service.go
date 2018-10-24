package reader

import "C"
import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/ibm-messaging/mq-golang/ibmmq"
)

func NewService() {
	var qMgrName string
	var err error
	var qMgr ibmmq.MQQueueManager
	var rc int

	if len(os.Args) != 4 {
		fmt.Println("clientconn <qmgrname> <channelname> <conname>")
		fmt.Println("")
		fmt.Println("For example")
		fmt.Println("  clientconn QMGR1 \"SYSTEM.DEF.SVRCONN\" \"myhost.example.com(1414)\"")
		fmt.Println("All parameters are required.")
		os.Exit(1)
	}

	// Which queue manager do we want to connect to
	qMgrName = os.Args[1]

	// Allocate the MQCNO and MQCD structures needed for the
	// MQCONNX call.
	cno := ibmmq.NewMQCNO()
	cd := ibmmq.NewMQCD()

	// Fill in the required fields in the
	// MQCD channel definition structure
	cd.ChannelName = os.Args[2]
	cd.ConnectionName = os.Args[3]

	// Reference the CD structure from the CNO
	// and indicate that we want to use the client
	// connection method.
	cno.ClientConn = cd
	cno.Options = C.MQCNO_CLIENT_BINDING

	// Also fill in the userid and password if the MQSAMP_USER_ID
	// environment variable is set. This is the same as the C
	// sample programs such as amqsput.
	userId := os.Getenv("MQSAMP_USER_ID")
	if userId != "" {
		scanner := bufio.NewScanner(os.Stdin)
		csp := ibmmq.NewMQCSP()
		csp.AuthenticationType = C.MQLONG(C.MQCSP_AUTH_USER_ID_AND_PWD)
		csp.UserId = userId

		fmt.Printf("Enter password for qmgr %s: \n", qMgrName)
		scanner.Scan()
		csp.Password = scanner.Text()

		// And make the CNO refer to the CSP structure
		cno.SecurityParms = csp
	}

	// And connect. Wait a short time before
	// disconnecting.
	qMgr, err = ibmmq.Connx(qMgrName, cno)
	if err == nil {
		fmt.Printf("Connection to %s succeeded.\n", qMgrName)
		d, _ := time.ParseDuration("5s")
		time.Sleep(d)
		qMgr.Disc()
		rc = 0
	} else {
		fmt.Printf("Connection to %s failed.\n", qMgrName)
		fmt.Println(err)
		rc = int(err.(*ibmmq.MQReturn).MQCC)
	}

	fmt.Println("Done.")
	os.Exit(rc)

}
