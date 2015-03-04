package parsing

import (
	"bytes"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/acapps/zipwhip-smtp/request"
    "github.com/acapps/zipwhip-smtp/address"
)

const (
	CONTENT_ENCODING = "content-transfer-encoding"
	CONTENT_TYPE     = "content-type"
	REPLY_TO         = "reply-to"
	SUBJECT          = "subject"
	ZIPWHIP_AUTH     = "zipwhip-auth"
	FROM             = "from"
)

type parseFunc func([]byte, *request.SendRequest) error

var DefaultParseTable = map[string]parseFunc{
	REPLY_TO: replyTo,
	FROM:     from,
}

const (
    SESSION_FORMAT = "^[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}:[0-9]"
    SUBJECT_VENDOR_FORMAT = `^[a-zA-Z0-9]{1,}::\+1[0-9]{10}@smtp.zipwhip.com`
    EMAIL_FORMAT = `\+1[0-9]{10}@smtp.zipwhip.com`
)

var (
    SESSION_MATCHER *regexp.Regexp
    SUBJECT_VENDOR_MATCHER *regexp.Regexp
    EMAIL_MATCHER *regexp.Regexp
)


func init() {
	SESSION_MATCHER = regexp.MustCompile(SESSION_FORMAT)
    SUBJECT_VENDOR_MATCHER = regexp.MustCompile(SUBJECT_VENDOR_FORMAT)
    EMAIL_MATCHER = regexp.MustCompile(EMAIL_FORMAT)
}

// Run through the Headers, apply any necessary formatting.
func Headers(sendRequest *request.SendRequest) error {

    for key, value := range sendRequest.Headers {

        sendRequest.Headers[key] = bytes.Trim(value, "\"")

        if _, ok := DefaultParseTable[key]; ok {

            var parseFunction parseFunc = DefaultParseTable[key]

            err := parseFunction(value, sendRequest)
            if err != nil {
                return err
            }
        }
    }

	return nil
}

func Recipients(recipients []string, sendRequest *request.SendRequest) error {

	for _, recipient := range recipients {
		nameAddress := strings.Replace(recipient, "\"", "", -1)

		address, err := mail.ParseAddress(nameAddress)
		if err != nil {
			return fmt.Errorf("Error parsing Address: %s because of error: %s", address, err)
		}

		address.Address = strings.Split(address.Address, "@")[0]

		if len(address.Address) >= 9 {
			sendRequest.Recipients = append(sendRequest.Recipients, []byte(address.Address))
		}
	}

	return nil
}

func SendingStrategy(sendRequest *request.SendRequest) error {

    if SESSION_MATCHER.Match(sendRequest.Key) {
        sendRequest.Strategy = request.SESSION
        return nil
    }

    if len(sendRequest.Key) > 0 {
        sendRequest.Strategy = request.VENDOR
        return nil
    }

    return fmt.Errorf("No sending strategy could be detected. %s", sendRequest.Key)
}

// Authentication is defined by meeting the following criteria
// A. Zipwhip-Auth populated VendorKey and From Address set to a valid Email, phoneNumber@smtp.zipwhip.com
// B. Zipwhip-Auth populated VendorKey and Subject set to a valid Email, phoneNumber@smtp.zipwhip.com
// C. Zipwhip Auth populated SessionKey
// D. Subject populated with VendorKey and From Address set to a valid Email, [vendorKey]::phoneNumber@smtp.zipwhip.com
// E. Subject populated with SessionKey
func Authentication(sendRequest *request.SendRequest) error {

    if _, ok := sendRequest.Headers[ZIPWHIP_AUTH]; ok {

        sendRequest.Key = sendRequest.Headers[ZIPWHIP_AUTH]

        err := authenticationZipwhipAuth(sendRequest)
        if err != nil {
            return fmt.Errorf("An error occurred while runing authentication on ZipwhipAuth. %s", err)
        }

        return nil
    }

    if _, ok := sendRequest.Headers[SUBJECT]; ok {

        sendRequest.Key = sendRequest.Headers[SUBJECT]

        err := authenticationSubject(sendRequest)
        if err != nil {
            return fmt.Errorf("An error occurred while runing authentication on Subject. %s", err)
        }

        return nil
	}

    return fmt.Errorf("No authorization mechanism found. %+s", sendRequest.Headers)
}

func authenticationZipwhipAuth(sendRequest *request.SendRequest) error {

    // C.
    if SESSION_MATCHER.Match(sendRequest.Headers[ZIPWHIP_AUTH]) {

        sendRequest.Strategy = request.SESSION
        return nil
    }

    _, from := sendRequest.Headers[FROM]
    _, subject := sendRequest.Headers[SUBJECT]

    if !from && !subject {
        return fmt.Errorf("Unable to locate appropriate Sender in from or subject.")
    }

    sendRequest.Strategy = request.VENDOR
    var HeaderToParse []byte

    if from {
        HeaderToParse = sendRequest.Headers[FROM]
    } else {
        HeaderToParse = sendRequest.Headers[SUBJECT]
    }

    if EMAIL_MATCHER.Match(HeaderToParse) {

        err := sendRequest.Sender.Parse(HeaderToParse)
        if err != nil {

            return fmt.Errorf("An error occurred while parsing the From Address: %s", err)
        }
        return nil
    }

    return address.IsValidZipwhipAddress(HeaderToParse)
}

func authenticationSubject(sendRequest *request.SendRequest) error {

    // E.
    if SESSION_MATCHER.Match(sendRequest.Headers[SUBJECT]) {

        sendRequest.Strategy = request.SESSION
        return nil
    }

    if SUBJECT_VENDOR_MATCHER.Match(sendRequest.Headers[SUBJECT]) {

        sendRequest.Strategy = request.VENDOR

        subjectParts := bytes.Split(bytes.Trim(sendRequest.Headers[SUBJECT], "\""), []byte("::"))
        sendRequest.Key = subjectParts[0]
        sendRequest.Sender.Parse(subjectParts[1])

        return nil
    }

    return fmt.Errorf("Unable to locate the proper sending strategy, key, and/or sender info.\nSubject: %s", sendRequest.Headers[SUBJECT])
}

func subject(subjectLine []byte, sendRequest *request.SendRequest) error {

	subjectLine = bytes.Replace(subjectLine, []byte("\""), []byte(""), -1)

	sendRequest.Key = subjectLine

	return nil
}

func replyTo(replyHeader []byte, sendRequest *request.SendRequest) error {

	nameAddress := bytes.Replace(replyHeader, []byte("\""), []byte(""), -1)

	address, err := mail.ParseAddress(string(nameAddress))
	if err != nil {
		return fmt.Errorf("Error parsing Address: %s becasue of error: %s", address, err)
	}

	sendRequest.ReplyTo = *address
	return nil
}

func from(fromHeader []byte, sendRequest *request.SendRequest) error {

    fromHeader = bytes.Trim(fromHeader, " <>")
    sendRequest.Headers[FROM] = fromHeader

	return nil
}
