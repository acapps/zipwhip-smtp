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

var Addresses map[string]string

func setup() {
    Addresses = make(map[string]string)

    Addresses[VALID] = "+18448982526"
    Addresses[VARIED_CASE] = "+18448982526"

    Addresses[BAD_HOST] = "bad"
    Addresses[MISSING_PLUS] = "bad"
    Addresses[MISSING_ONE] = "bad"
    Addresses[DIGIT_SHORT] = "bad"
}

func Test_extractSender(t *testing.T) {

    setup()

    Convey("Attempt various Addresses.", t, func() {

        for k,v := range Addresses {

            err := IsValidZipwhipAddress([]byte(k))
            if err != nil {

                Convey(fmt.Sprintf("This Address, %s, should be invalid", k), func() {
                    So(v, ShouldEqual, "bad")
                })

            } else {

                Convey(fmt.Sprintf("This Address, %s, should be valid", k), func() {
                    result, err := extractSender([]byte(k))
                    So(err, ShouldBeNil)
                    So(string(result), ShouldEqual, v)
                })
            }
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
                    So(v, ShouldEqual, "bad")
                })
                continue
            }
            Convey(fmt.Sprintf("This Address, %s, should be valid", k), func() {
                So(err, ShouldBeNil)
            })
        }
    })
}
