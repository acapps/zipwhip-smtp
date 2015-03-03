package address

import (
	"bytes"
	"fmt"
	"net/mail"
	"regexp"
)

type ZipwhipAddress struct {
	Sender  []byte
	Address mail.Address
}

const (
	EMAIL_FORMAT = `\+1[0-9]{10}@smtp.zipwhip.com`
	PHONE_FORMAT = `\+1[0-9]{10}`
)

var (
	EMAIL_MATCHER *regexp.Regexp
	PHONE_MATCHER *regexp.Regexp
)

func init() {
	EMAIL_MATCHER = regexp.MustCompile(EMAIL_FORMAT)
	PHONE_MATCHER = regexp.MustCompile(PHONE_FORMAT)
}

func NewZipwhipAddress() *ZipwhipAddress {

	return new(ZipwhipAddress)
}

func IsValidZipwhipAddress(emailAddress []byte) error {

    emailAddress = bytes.ToLower(emailAddress)
	if EMAIL_MATCHER.Match(emailAddress) {

		return nil
	}

	sender, err := extractSender(emailAddress)
	if err != nil {

		fmt.Errorf("Phone number portion was invalid, %s", sender)
	}

	return fmt.Errorf("Host portion of email address was invalid.")
}

func extractSender(emailAddress []byte) ([]byte, error) {

	const (
		PHONE_NUMBER = iota
	)

    if !EMAIL_MATCHER.Match(bytes.ToLower(emailAddress)) {
        return nil, fmt.Errorf("Email address is invalid.")
    }

	addressParts := bytes.Split(emailAddress, []byte("@"))
	if PHONE_MATCHER.Match(addressParts[PHONE_NUMBER]) {

		return []byte(addressParts[PHONE_NUMBER]), nil
	}

	return addressParts[PHONE_NUMBER], fmt.Errorf("Invalid phone number format")
}

func (za *ZipwhipAddress) Parse(emailAddress []byte) error {

	err := IsValidZipwhipAddress(emailAddress)
	if err != nil {

		return fmt.Errorf("Parse failed due to: %s", err)
	}

	za.Address.Address = string(emailAddress)

	sender, err := extractSender(emailAddress)
	if err != nil {
		return fmt.Errorf("Failed to extract sender, %s", err)
	}

	za.Sender = sender

	return nil
}
