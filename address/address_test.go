package address

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    "fmt")

const (
    VALID = "+18448982526@smtp.zipwhip.com"
    VARIED_CASE = "+18448982526@SMTP.Zipwhip.com"
    BAD_HOST = "+18448982526@zipwhip.com"
    MISSING_PLUS = "18448982526@smtp.zipwhip.com"
    MISSING_ONE = "+8448982526@smtp.zipwhip.com"
    DIGIT_SHORT = "+1844898252@smtp.zipwhip.com"
)

var Addresses map[string]bool

func setup() {
    Addresses = make(map[string]bool)

    Addresses[VALID] = true
    Addresses[VARIED_CASE] = true

    Addresses[BAD_HOST] = false
    Addresses[MISSING_PLUS] = false
    Addresses[MISSING_ONE] = false
    Addresses[DIGIT_SHORT] = false
}

func Test_extractSender(t *testing.T) {

    setup()

    Convey("Attempt various Addresses.", t, func() {

        for k,v := range Addresses {

            _, err := extractSender([]byte(k))
            if err != nil {
                Convey(fmt.Sprintf("This Address, %s, should be invalid", k), func() {
                    So(v, ShouldBeFalse)
                })
                continue
            }
            Convey(fmt.Sprintf("This Address, %s, should be valid", k), func() {
                So(v, ShouldBeTrue)
            })
        }
    })
}

func Test_IsValidZipwhipAddress(t *testing.T) {

    setup()

    Convey("Attempt various Addresses.", t, func() {

        for k,v := range Addresses {

            err := IsValidZipwhipAddress([]byte(k))
            if err != nil {
                Convey(fmt.Sprintf("This Address, %s, should be invalid", k), func() {
                    So(v, ShouldBeFalse)
                })
                continue
            }
            Convey(fmt.Sprintf("This Address, %s, should be valid", k), func() {
                So(v, ShouldBeTrue)
            })
        }
    })
}
