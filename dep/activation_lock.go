package dep

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	activationLockPath       = "device/activationlock"
	activationUnlockEndpoint = "https://deviceservices-external.apple.com/deviceservicesworkers/escrowKeyUnlock"
)

type ActivationLockRequest struct {
	Device string `json:"device"`

	//  If the escrow key is not provided, the device will be locked to the person who created the MDM server in the portal.
	// https://developer.apple.com/documentation/devicemanagement/device_assignment/activation_lock_a_device/creating_and_using_bypass_codes
	// The EscrowKey is a hex-encoded PBKDF2 derivation of the bypass code. See activationlock.BypassCode.
	EscrowKey string `json:"escrow_key"`

	LostMessage string `json:"lost_message"`
}

type ActivationLockResponse struct {
	SerialNumber string `json:"serial_number"`
	Status       string `json:"response_status"`
}

func (c *Client) ActivationLock(alr *ActivationLockRequest) (*ActivationLockResponse, error) {
	req, err := c.newRequest("POST", activationLockPath, &alr)
	if err != nil {
		return nil, errors.Wrap(err, "create activation lock request")
	}

	var response ActivationLockResponse
	err = c.do(req, &response)
	return &response, errors.Wrap(err, "activation lock")
}

type ActivationUnlockRequest struct {
	Querystring string `json:"querystring"`
	Messagebody string `json:"messagebody"`
}

func (c *Client) ActivationUnlock(alur *ActivationUnlockRequest) (interface{}, error) {
	cert, err1 := tls.LoadX509KeyPair("/root/apns/appleapncert.pem", "/root/apns/private.key")
	if err1 != nil {
		log.Fatal(err1)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpclient := &http.Client{Transport: transport}

	URLparts := []string{activationUnlockEndpoint, "?", alur.Querystring}
	var RequestURL = strings.Join(URLparts, "")
	var Messagebody = alur.Messagebody

	var buffer bytes.Buffer
	buffer.WriteString(Messagebody)
	fmt.Println(buffer.String())
	req, err := http.NewRequest("POST", RequestURL, &buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "*/*")

	var response *http.Response
	response, _ = httpclient.Do(req)
	data, err := ioutil.ReadAll(response.Body)
	return data, err
}
