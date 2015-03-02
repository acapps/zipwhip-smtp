package sending

import (
	"fmt"
	"github.com/acapps/smtp-test/server/request"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

func SessionKey(request request.SendRequest) {

	for i := 0; i < len(request.Recipients); i++ {
		getRequest := fmt.Sprintf("https://api.zipwhip.com/message/send?session=%s&destinationAddress=%s&body=%s", request.Key, request.Recipients[i], url.QueryEscape(string(request.Body)))
		go log.Debugln(getRequest)

		_, err := http.Get(getRequest)
		if err != nil {
			go log.Debugf("An error was encountored while sending a message: %s", err)
		}
	}
}

func VendorKey(request request.SendRequest) {

	for i := 0; i < len(request.Recipients); i++ {
		getRequest := fmt.Sprintf("https://vendor.zipwhip.com/message/send?vendorKey=%s&sourceAddress=%s&destinationAddress=%s&body=%s", request.Key, url.QueryEscape(request.Sender.Address), request.Recipients[i], url.QueryEscape(string(request.Body)))
		go log.Debugln(getRequest)

		_, err := http.Get(getRequest)
		if err != nil {
			go log.Debugf("An error was encountored while sending a message: %s", err)
		}
	}
}
