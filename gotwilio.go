// Package gotwilio is a library for interacting with http://www.twilio.com/ API.
package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL       = "https://api.twilio.com/2010-04-01"
	videoURL      = "https://video.twilio.com"
	pricingURL    = "https://pricing.twilio.com/v1"
	messagingURL  = "https://messaging.twilio.com/v1"
	notifyURL     = "https://notify.twilio.com/v1"
	clientTimeout = time.Second * 30
)

// The default http.Client that is used if none is specified
var defaultClient = &http.Client{
	Timeout: time.Second * 30,
}

// Twilio stores basic information important for connecting to the
// twilio.com REST api such as AccountSid and AuthToken.
type Twilio struct {
	AccountSid   string
	AuthToken    string
	BaseUrl      string
	VideoUrl     string
	PricingUrl   string
	MessagingUrl string
	NotifyUrl    string
	HTTPClient   *http.Client

	APIKeySid    string
	APIKeySecret string
}

// Exception is a representation of a twilio exception.
type Exception struct {
	Status   int    `json:"status"`    // HTTP specific error code
	Message  string `json:"message"`   // HTTP error message
	Code     int    `json:"code"`      // Twilio specific error code
	MoreInfo string `json:"more_info"` // Additional info from Twilio
}

// Create a new Twilio struct.
func NewTwilioClient(accountSid, authToken string) *Twilio {
	return NewTwilioClientCustomHTTP(accountSid, authToken, nil)
}

// Create a new Twilio client, optionally using a custom http.Client
func NewTwilioClientCustomHTTP(accountSid, authToken string, HTTPClient *http.Client) *Twilio {
	if HTTPClient == nil {
		HTTPClient = defaultClient
	}

	return &Twilio{
		AccountSid:   accountSid,
		AuthToken:    authToken,
		BaseUrl:      baseURL,
		VideoUrl:     videoURL,
		PricingUrl:   pricingURL,
		MessagingUrl: messagingURL,
		NotifyUrl:    notifyURL,
		HTTPClient:   HTTPClient,
	}
}

func (twilio *Twilio) WithAPIKey(apiKeySid string, apiKeySecret string) *Twilio {
	twilio.APIKeySid = apiKeySid
	twilio.APIKeySecret = apiKeySecret
	return twilio
}

func (twilio *Twilio) post(formValues url.Values, twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return twilio.do(req)
}

func (twilio *Twilio) get(twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", twilioUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)

	return twilio.do(req)
}

func (twilio *Twilio) delete(twilioUrl string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", twilioUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)

	return twilio.do(req)
}

func (twilio *Twilio) do(req *http.Request) (*http.Response, error) {
	client := twilio.HTTPClient
	if client == nil {
		client = defaultClient
	}

	return client.Do(req)
}

func (twilio *Twilio) getResponseBody(url string) (responseBody []byte, exception *Exception, err error) {
	res, err := twilio.get(url)
	if err != nil {
		return responseBody, exception, err
	}
	defer res.Body.Close()

	responseBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return responseBody, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return responseBody, exception, err
	}

	return responseBody, exception, err
}

func (twilio *Twilio) DeleteResource(servicesUrl string) (exception *Exception, err error) {
	res, err := twilio.delete(servicesUrl)
	if err != nil {
		return exception, err
	}

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNoContent {
		exc := new(Exception)
		err = json.Unmarshal(respBody, exc)
		return exc, err
	}
	return nil, nil
}