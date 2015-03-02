package parsing

import (
    "strings"

    "github.com/acapps/smtp-test/server/request"
    "fmt"
    "bytes"
    "net/mail"
    "regexp")

const (
    CONTENT_ENCODING = "content-transfer-encoding"
    CONTENT_TYPE = "content-type"
    REPLY_TO = "reply-to"
    SUBJECT = "subject"
    ZIPWHIP_AUTH = "zipwhip-auth"
    FROM = "from"
)

type parseFunc func([]byte, *request.SendRequest) error

var ZipwhipAuthParseTable = map[string]parseFunc {
    ZIPWHIP_AUTH: zipwhipAuth,
}

var FallbackAuthParseTable = map[string]parseFunc {
    SUBJECT: subject,
}

var DefaultParseTable = map[string]parseFunc {
    REPLY_TO: replyTo,
    FROM: from,
}

const SESSION_FORMAT = "^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}:[0-9]$"
var SESSION_MATCHER *regexp.Regexp

func init() {
    SESSION_MATCHER = regexp.MustCompile(SESSION_FORMAT)
}

func Headers(sendRequest *request.SendRequest) error {

    err := authentication(sendRequest)
    if err != nil {
        return fmt.Errorf("Authentication failed: %s", err)
    }

    for key, value := range sendRequest.Headers {

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

func authentication(sendRequest *request.SendRequest) error {

    if _, ok := sendRequest.Headers[ZIPWHIP_AUTH]; ok {

        var parseFunction parseFunc = ZipwhipAuthParseTable[ZIPWHIP_AUTH]

        err := parseFunction(sendRequest.Headers[ZIPWHIP_AUTH], sendRequest)
        if err != nil {
            return err
        }
    } else if _, ok = sendRequest.Headers[SUBJECT]; ok {

        var parseFunction parseFunc = FallbackAuthParseTable[SUBJECT]

        err := parseFunction(sendRequest.Headers[SUBJECT], sendRequest)
        if err != nil {
            return err
        }
    } else {
        return fmt.Errorf("No authorization mechanism found. %+s", sendRequest.Headers)
    }

    err := sendingStrategy(sendRequest)
    if err != nil {
        return fmt.Errorf("authenication failed due to error determining sending strategy. %s", err)
    }

    return nil
}

func subject(subjectLine []byte, sendRequest *request.SendRequest) error {

    subjectLine = bytes.Replace(subjectLine, []byte("\""), []byte(""), -1)

    sendRequest.Key = subjectLine
    sendRequest.Strategy = request.VENDOR

    return nil
}

func zipwhipAuth(vendorKey []byte, sendRequest *request.SendRequest) error {

    sendRequest.Key = vendorKey
    sendRequest.Strategy = request.VENDOR

    return nil
}

func sendingStrategy(sendRequest *request.SendRequest) error {

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

    nameAddress := bytes.Replace(fromHeader, []byte("\""), []byte(""), -1)

    address, err := mail.ParseAddress(string(nameAddress))
    if err != nil {
        return fmt.Errorf("Error parsing Address: %s because of error: %s", address, err)
    }

    address.Address = strings.Split(address.Address, "@")[0]

    sendRequest.Sender = *address
    return nil
}