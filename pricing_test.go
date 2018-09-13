package gotwilio

import (
	"fmt"
	"testing"
)

var paramsLive map[string]string

func init() {
	paramsLive = make(map[string]string)

	// Only LIVE credentials possible, because of 20008 error
	paramsLive["SID"] = ""
	paramsLive["TOKEN"] = ""
}

func TestGetPricing(t *testing.T) {
	countryISO := "EE"
	twilio := NewTwilioClient(paramsLive["SID"], paramsLive["TOKEN"])
	_, exc, err := twilio.GetPricing(countryISO)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestTwilioGetCountriesPricing(t *testing.T) {
	twilio := NewTwilioClient(paramsLive["SID"], paramsLive["TOKEN"])
	c, exc, err := twilio.GetCountriesPricing()

	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

	fmt.Println(c)
}

func TestGetPhoneNumbers(t *testing.T) {
	twilio := NewTwilioClient(paramsLive["SID"], paramsLive["TOKEN"])
	c, exc, err := twilio.GetPhoneNumbers()

	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

	fmt.Println(c)
}