package request

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAddBody_Ascii(t *testing.T) {
	var input = []byte("A new message body.")
	var output = []byte("A new message body.")

	Convey("Create an initial SendRequest Object and Ascii message body.", t, func() {
		sr := NewSendRequest()

		Convey("Add the Ascii body to the SendRequest.", func() {
			sr.AddBody(input)

			Convey("SendRequest.Body should be equal to output", func() {
				So(sr.Body, ShouldResemble, output)
			})
		})
	})
}

func TestAddBody_AsciiWithNewline(t *testing.T) {
	var input = []byte("A new message body.\n")
	var output = []byte("A new message body.")

	Convey("Create an initial SendRequest Object and Ascii message body.", t, func() {
		sr := NewSendRequest()

		Convey("Add the Ascii body to the SendRequest.", func() {
			sr.AddBody(input)

			Convey("SendRequest.Body should be equal to output", func() {
				So(sr.Body, ShouldResemble, output)
			})
		})
	})
}

func TestAddBody_Base64(t *testing.T) {
	var expectedOutput = []byte("A new message body.")
	var input = []byte("QSBuZXcgbWVzc2FnZSBib2R5Lg==")

	Convey("Create an initial SendRequest Object", t, func() {
		sr := NewSendRequest()

		Convey("Add the base64 body to the SendRequest.", func() {
			sr.AddBody(input)

			Convey("SendRequest.Body should be equal to expectedOutput", func() {
				So(sr.Body, ShouldResemble, expectedOutput)
			})
		})
	})
}
