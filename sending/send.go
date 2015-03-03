package sending

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/acapps/zipwhip-smtp/request"
	log "github.com/sirupsen/logrus"
)

func SessionKey(request request.SendRequest) {

	for i := 0; i < len(request.Recipients); i++ {
		getRequest := fmt.Sprintf("https://api.zipwhip.com/message/send?session=%s&contacts=%s&body=%s", request.Key, request.Recipients[i], url.QueryEscape(string(request.Body)))
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

/*

body, err := ioutil.ReadAll(r.Body)
if err != nil {
log.Printf("An error occurred while reading in the Request's Body:\n\t%s", err)
}
go bodyToMessage(&body)

*/