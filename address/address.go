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
	EMAIL_FORMAT = `^\+1[0-9]{10}@smtp.zipwhip.com$`
	PHONE_FORMAT = `^\+1[0-9]{10}$`
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

    // Leverage, the email package for parsing.
    var email *mail.Address
    emailAddress = bytes.ToLower(emailAddress)

    email, err := mail.ParseAddress(string(emailAddress))
    if err != nil {
        return fmt.Errorf("Error parsing emailAddress. %s", err)
    }

	if EMAIL_MATCHER.Match([]byte(email.Address)) {

		return nil
	}

    return fmt.Errorf("Email address was not in the proper format. %s", email.Address)
}

func extractSender(emailAddress []byte) ([]byte, error) {

	const (
		PHONE_NUMBER = iota
	)

	addressParts := bytes.Split(emailAddress, []byte("@"))
	if PHONE_MATCHER.Match(addressParts[PHONE_NUMBER]) {

		return []byte(addressParts[PHONE_NUMBER]), nil
	}

	return nil, fmt.Errorf("Invalid phone number format")
}

func (za *ZipwhipAddress) Parse(emailAddress []byte) error {

	err := IsValidZipwhipAddress(emailAddress)
	if err != nil {

		return fmt.Errorf("Parse failed due to: %s", err)
	}

	sender, err := extractSender(emailAddress)
	if err != nil {
		return fmt.Errorf("Failed to extract sender, %s", err)
	}

	za.Sender = sender

	return nil
}
