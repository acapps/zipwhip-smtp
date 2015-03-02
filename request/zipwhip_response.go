package request

import (
    "encoding/json"
    "fmt"
)

type ZipwhipResponse struct {
    Fingerprint []byte
    MessageId   []byte
}

type VendorResponse struct {
    Fingerprint string `json:"fingerprint"`
    Root        string `json:"root"`
}

type SessionResponse struct {
    Response struct {
        Fingerprint string `json:"fingerprint"`
        Root        string `json:"root"`
    } `json:"response"`
    Success bool `json:"success"`
}

func (vr *VendorResponse) ParseJson(inputJson *[]byte) error {
    return json.Unmarshal(*inputJson, vr)
}

func (sr *SessionResponse) ParseJson(inputJson *[]byte) error {
    return json.Unmarshal(*inputJson, sr)
}

func (vr *VendorResponse) parseVendorResponse(body *[]byte, sr *SendRequest) error {
    var vr VendorResponse
    err := vr.ParseJson(body)
    if err != nil {
        return fmt.Errorf("An error occurred while parsing Vendor response: %s", body)
    }

    return nil
}

func (sr *SessionResponse) parseSessionResponse(body *[]byte, request *SendRequest) error {
    var sr SessionResponse
    err := sr.ParseJson(body)
    if err != nil {
        return fmt.Errorf("An error occurred while parsing Vendor response: %s", body)
    }

    return nil
}
