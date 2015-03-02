package parsing

import (
	"testing"

	"github.com/acapps/zipwhip-smtp/request"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	smtpdRecipients = "+18448982526@zipwhip.com"
)

const (
	VendorHeader        = "Zipwhip-Auth: vendorKey\n"
	ComplexVendorHeader = "Zipwhip-Auth: 8ef1211f-d9f2-4c81-906f-7d27da5a32f8\n"
	VendorSubjectHeader = "Subject: \"vendorKey\"\n"
	SessionHeader       = "Subject: \"8ef1211f-d9f2-4c81-906f-7d27da5a32f8:309626613\"\n"
	ReplyToHeader       = "Reply-To: \"Alan\" <alan@zipwhip.com>\n"
	FromHeader          = "Reply-To: \"Desk\" <+18448982526zipwhip.com>\n"
	BadVendorHeader     = "Zipwhip-Auth:\n"
	BadSessionHeader    = "Subject:\n"
)

const (
	VendorKey        = "vendorKey"
	ComplexVendorKey = "8ef1211f-d9f2-4c81-906f-7d27da5a32f8"
	SessionKey       = "8ef1211f-d9f2-4c81-906f-7d27da5a32f8:309626613"

	ReplyTo = "alan@zipwhip.com"
	From    = "+18448982526"
)

func Test_AuthorizationVendor(t *testing.T) {

	var Headers = buildHeader([]byte(VendorHeader))
	sr := request.NewSendRequest()

	Convey("Parse Headers to setup tests should not return error.", t, func() {
		err := sr.AddHeaders(Headers)
		So(err, ShouldBeNil)
	})

	Convey("Parse headers through Authorization.", t, func() {
		err := authentication(sr)
		So(err, ShouldBeNil)
	})

	Convey("VendorKey should now be assigned to the Key of the request.", t, func() {
		So(string(sr.Key), ShouldEqual, VendorKey)
	})
}

func Test_AuthorizationVendorInSubject(t *testing.T) {

	var Headers = buildHeader([]byte(VendorSubjectHeader))
	sr := request.NewSendRequest()

	Convey("Parse Headers to setup tests should not return error.", t, func() {
		err := sr.AddHeaders(Headers)
		So(err, ShouldBeNil)
	})

	Convey("Parse headers through Authorization.", t, func() {
		err := authentication(sr)
		So(err, ShouldBeNil)
	})

	Convey("VendorKey should now be assigned to the Key of the request.", t, func() {
		So(string(sr.Key), ShouldEqual, VendorKey)
	})
}

func Test_AuthorizationSession(t *testing.T) {

	var Headers = buildHeader([]byte(SessionHeader))
	sr := request.NewSendRequest()

	Convey("Parse Headers to setup tests should not return error.", t, func() {
		err := sr.AddHeaders(Headers)
		So(err, ShouldBeNil)
	})

	Convey("Parse headers through Authorization.", t, func() {
		err := authentication(sr)
		So(err, ShouldBeNil)
	})

	Convey("SessionKey should now be assigned to the Key of the request.", t, func() {
		So(string(sr.Key), ShouldEqual, SessionKey)
	})
}

func Test_AuthorizationComplexVendor(t *testing.T) {

	var Headers = buildHeader([]byte(ComplexVendorHeader))
	sr := request.NewSendRequest()

	Convey("Parse Headers to setup tests should not return error.", t, func() {
		err := sr.AddHeaders(Headers)
		So(err, ShouldBeNil)
	})

	Convey("Parse headers through Authorization.", t, func() {
		err := authentication(sr)
		So(err, ShouldBeNil)
	})

	Convey("ComplexVendorKey should now be assigned to the Key of the request.", t, func() {
		So(string(sr.Key), ShouldEqual, ComplexVendorKey)
	})
}

func Test_AuthorizationFallback(t *testing.T) {

	var Headers = buildHeader([]byte(BadVendorHeader), []byte(SessionHeader))
	sr := request.NewSendRequest()

	Convey("Parse Headers to setup tests should not return error.", t, func() {
		err := sr.AddHeaders(Headers)
		So(err, ShouldBeNil)
	})

	Convey("Parse headers through Authorization.", t, func() {
		err := authentication(sr)
		So(err, ShouldBeNil)
	})

	Convey("SessionKey should now be assigned to the Key of the request.", t, func() {
		So(string(sr.Key), ShouldEqual, SessionKey)
	})
}

func Test_AuthorizationFailure(t *testing.T) {

	var Headers = buildHeader([]byte(BadVendorHeader), []byte(BadSessionHeader))
	sr := request.NewSendRequest()

	Convey("Parse Headers to setup tests should not return error.", t, func() {
		err := sr.AddHeaders(Headers)
		So(err, ShouldBeNil)
	})

	Convey("Parse headers through Authorization.", t, func() {
		err := authentication(sr)
		So(err, ShouldNotBeNil)
	})
}

func Test_HeaderParsing(t *testing.T) {

	var headers = buildHeader([]byte(VendorHeader), []byte(SessionHeader), []byte(ReplyToHeader))
	sr := request.NewSendRequest()

	Convey("Parse Headers to setup tests should not return error.", t, func() {
		err := sr.AddHeaders(headers)
		So(err, ShouldBeNil)
	})

	Convey("Parse headers through parser's Header.", t, func() {
		err := Headers(sr)
		So(err, ShouldBeNil)
	})
}

func Test_RecipientParsing(t *testing.T) {

	sr := request.NewSendRequest()
	var headers []string
	headers = append(headers, smtpdRecipients)

	Convey("Parse Recipients.", t, func() {
		err := Recipients(headers, sr)
		So(err, ShouldBeNil)
	})
}

func Test_ReplyTo(t *testing.T) {

	var Headers = buildHeader([]byte(ReplyToHeader))
	sr := request.NewSendRequest()

	Convey("Parse Headers to setup tests should not return error.", t, func() {
		err := sr.AddHeaders(Headers)
		So(err, ShouldBeNil)
	})

	Convey("Parse the ReplyTo Header, parsing should not return an error.", t, func() {
		err := replyTo(sr.Headers["reply-to"], sr)
		t.Logf("%+s", Headers)
		So(err, ShouldBeNil)
	})

	Convey("ReplyTo should now be assigned to the ReplyTo of the request.", t, func() {
		So(string(sr.ReplyTo.Address), ShouldEqual, ReplyTo)
	})
}

func buildHeader(input ...[]byte) (output []byte) {

	for _, v := range input {
		output = append(output, v...)
	}

	return output
}
