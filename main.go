package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"bitbucket.org/chrj/smtpd"
	"code.google.com/p/gcfg"
	"github.com/acapps/zipwhip-smtp/parsing"
	"github.com/acapps/zipwhip-smtp/request"
	"github.com/acapps/zipwhip-smtp/sending"
	log "github.com/sirupsen/logrus"
)

/* GLOBAL */
type Config struct {
	Server struct {
		Port     string
		Address  string
		LogLevel int
	}
	MailServer struct {
		IpFilter string
	}
	Zipwhip struct {
		SessionKey string
		VendorKey  string
	}
}

var config Config

const (
	_ = iota
	Open
	Closed
)

const (
	_ = iota
	Subject
	Vendor
	Session
)

var sendingStrategy = Subject
var clientFilter = Open

func init() { // Init will run with unit Tests.

	configFile := flag.String("configFile", "testing.config", "config file")
	sendingStrategy = *flag.Int("sendingStrategy", 0, "1 = subject field, 2 = vendorKey, or 3 = sessionKey")
	flag.Parse()

	err := gcfg.ReadFileInto(&config, *configFile)
	if err != nil {
		log.Panicf("Could not read Config File.")
	}

	log.SetLevel(log.Level(config.Server.LogLevel))

	if len(config.MailServer.IpFilter) > 0 {
		clientFilter = Closed
	}

	// We would prefer to use Vendor Send.
	if len(config.Zipwhip.VendorKey) > 0 {
		sendingStrategy = Vendor
		return
	}

	if len(config.Zipwhip.SessionKey) > 0 {
		sendingStrategy = Session
	}
}

func main() {
	var server *smtpd.Server

	server = &smtpd.Server{

		HeloChecker: func(peer smtpd.Peer, name string) error {

			return nil
		},

		Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {

			go log.Debugf("New connection: %+s, %+s", peer, env)
			go parseRequest(peer, env)
			return nil
		},
	}

	go log.Warnf("Server is listening on: %s:%s", config.Server.Address, config.Server.Port)
	go log.Warnf("%+v", config)

	serverConfiguration := func() string {
		if clientFilter == Open {
			return "Server is set to Open."
		}
		return fmt.Sprintf("Server is set to IP Locked.\n%+v", config.MailServer.IpFilter)
	}
	go log.Warn(serverConfiguration())

	server.ListenAndServe(config.Server.Address + ":" + config.Server.Port)
}

func parseRequest(peer smtpd.Peer, env smtpd.Envelope) {

	if clientFilter == Closed {
		if !strings.HasPrefix(peer.Addr.String(), config.MailServer.IpFilter+":") {
			go log.Debugf("Connection was refused due to the IP Filter: %s", peer.Addr)
			return
		}
	}

	request := request.NewSendRequest()

	err := parseMessage(env.Data, request)
	if err != nil {
        go log.Debugf("Error parsing message: %s", err)
		return
	}

	err = parsing.Recipients(env.Recipients, request)
	if err != nil {
        go log.Debugf("Error parsing recipients: %s", err)
		return
	}

    go log.Debugf("%+s", request)

	sendMessages(request)
}

func parseMessage(body []byte, sendRequest *request.SendRequest) error {

	const (
		HEADERS = 0
		BODY    = 1
	)
	// Headers and Body separated by '\n\n'
	// All other instances are assumed as part of the body.
	headersAndBody := bytes.Split(body, []byte("\n\n"))

	if len(headersAndBody) < 2 {
		go log.Debugf("Not enough segments, %d", len(headersAndBody))
		return fmt.Errorf("Improperly formatted message.")
	}

	err := sendRequest.AddBody(headersAndBody[BODY])
	if err != nil {
		return err
	}

	err = sendRequest.AddHeaders(headersAndBody[HEADERS])
	if err != nil {
		return err
	}

	err = parsing.Headers(sendRequest) // Break all headers into their own element
	if err != nil {
		go log.Warnf("Error occurred while parsing the header: %s", err)
		return err
	}

    err = parsing.Authentication(sendRequest)
    if err != nil {
        return fmt.Errorf("Authentication failed: %s", err)
    }

    err = parsing.SendingStrategy(sendRequest)
    if err != nil {
        return fmt.Errorf("Unable to determine sendingStrategy: %s", err)
    }

	return nil
}

// Send the message to each recipient.
// Messages will be truncated automatically.
// No status is returned.
func sendMessages(sendRequest *request.SendRequest) {

	switch sendRequest.Strategy {
	case request.SESSION:
		sending.SessionKey(*sendRequest)
	case request.VENDOR:
		sending.VendorKey(*sendRequest)
	default:
		go log.Warnf("SendMessages came across a default scenario, when it shouldn't have %s", sendingStrategy)
	}
}
