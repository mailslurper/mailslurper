package mailslurper_test

import (
	"net"
	"os"
	"sync"
	"time"

	"github.com/mailslurper/mailslurper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/adampresley/webframework/logging2"
)

var killChannel chan bool
var wg *sync.WaitGroup
var database mailslurper.IStorage
var smtpServer net.Listener

func TestMailslurper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mailslurper Suite")
}

var _ = BeforeSuite(func() {
	var err error

	killChannel = make(chan bool, 1)
	wg = &sync.WaitGroup{}

	logger := logging2.LogFactory(logging2.LOG_FORMAT_SIMPLE, "MailSlurper", logging2.ERROR)
	logger.EnableColors()

	databaseConnection := &mailslurper.ConnectionInformation{
		Filename: "./temp.db",
	}

	if database, err = mailslurper.ConnectToStorage(mailslurper.STORAGE_SQLITE, databaseConnection, logger); err != nil {
		logger.Errorf("Error connecting to storage")
		os.Exit(-1)
	}

	pool := mailslurper.NewServerPool(logger, 10)
	config := &mailslurper.Configuration{
		SMTPAddress: "localhost",
		SMTPPort:    2500,
	}

	if smtpServer, err = mailslurper.SetupSMTPServerListener(config, logger); err != nil {
		logger.Errorf("There was a problem starting the SMTP listener: %s", err.Error())
		os.Exit(0)
	}

	receivers := []mailslurper.IMailItemReceiver{
		mailslurper.NewDatabaseReceiver(database, logger),
	}

	/*
	 * Start the SMTP dispatcher
	 */
	go mailslurper.Dispatch(pool, smtpServer, receivers, logger, killChannel, wg)
})

var _ = AfterSuite(func() {
	killChannel <- true
	wg.Wait()
	database.Disconnect()
	mailslurper.CloseSMTPServerListener(smtpServer)
})

func DeleteAllMail() {
	database.DeleteMailsAfterDate("")
}

func Pause() {
	time.Sleep(100 * time.Millisecond)
}
