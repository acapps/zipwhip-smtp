package request

import (
	"bytes"
	"encoding/base64"
	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"net/mail"
)

type SendingStrategy int

const SESSION SendingStrategy = 1
const VENDOR SendingStrategy = 2

const (
    CONTENT_ENCODING = "content-transfer-encoding"
    CONTENT_TYPE     = "content-type"
    REPLY_TO         = "reply-to"
    SUBJECT          = "subject"
    ZIPWHIP_AUTH     = "zipwhip-auth"
    FROM             = "from"
)

var ParsingTable = map[string]bool {
    FROM: true,
    SUBJECT: true,
    CONTENT_ENCODING: true,
    CONTENT_TYPE: true,
    REPLY_TO: true,
    ZIPWHIP_AUTH: true,
}

type SendRequest struct {
	Key        []byte
	Sender     mail.Address
	Strategy   SendingStrategy
	Recipients [][]byte
	ReplyTo    mail.Address
	Body       []byte
	Headers    map[string][]byte
	Result     ZipwhipResponse
}

func NewSendRequest() *SendRequest {
	return new(SendRequest)
}

func (sr *SendRequest) AddBody(body []byte) error {

	body = bytes.TrimSuffix(body, []byte("\n"))

	if govalidator.IsBase64(string(body)) {
		base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
		length, err := base64.StdEncoding.Decode(base64Text, body)
		if err != nil {
			go log.Warnf("Error decoding what we believe is a base64 encoded body. %s", err)
			return err
		}
		sr.Body = base64Text[:length]
		return nil
	}

	sr.Body = body

	return nil
}

func (sr *SendRequest) AddHeaders(headers []byte) error {

	const (
		KEY   = 0
		VALUE = 1
	)

	headers = bytes.TrimSpace(headers)

	headerArray := bytes.Split(headers, []byte("\n"))

	sr.Headers = make(map[string][]byte, len(ParsingTable))

	for i := 0; i < len(headerArray); i++ {
		header := bytes.SplitN(headerArray[i], []byte(":"), 2)
        header[KEY] = bytes.ToLower(header[KEY])
        go log.Debugf("%s", header)
        if _,ok := ParsingTable[string(header[KEY])]; !ok {
            continue
        }
		if len(header[VALUE]) > 0 {
			sr.Headers[string(bytes.ToLower(header[KEY]))] = bytes.TrimSpace(header[VALUE])
		}
	}

	return nil
}
