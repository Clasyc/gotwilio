package gotwilio

import (
	"encoding/json"
)

type IncomingPhoneNumbers struct {
	Page                 int                   `json:"page"`
	PageSize             int                   `json:"page_size"`
	IncomingPhoneNumbers []IncomingPhoneNumber `json:"incoming_phone_numbers"`
}

type IncomingPhoneNumber struct {
	AccountSid           string            `json:"account_sid"`
	Sid                  string            `json:"sid"`
	PhoneNumber          string            `json:"phone_number"`
	FriendlyName         string            `json:"friendly_name"`
	DateCreated          string            `json:"date_created"`
	AddressRequirements  string            `json:"address_requirements"`
	APIVersion           string            `json:"api_version"`
	Beta                 bool              `json:"beta"`
	Capabilities         *NumberCapability `json:"capabilities"`
	DateUpdated          string            `json:"date_updated"`
	EmergencyAddressSid  string            `json:"emergency_address_sid"`
	EmergencyStatus      string            `json:"emergency_status"`
	SMSApplicationSid    string            `json:"sms_application_sid"`
	SMSFallbackMethod    string            `json:"sms_fallback_method"`
	SMSFallbackURL       string            `json:"sms_fallback_url"`
	SMSMethod            string            `json:"sms_method"`
	SMSURL               string            `json:"sms_url"`
	StatusCallback       string            `json:"status_callback"`
	StatusCallbackMethod string            `json:"status_callback_method"`
	TrunkSid             string            `json:"trunk_sid"`
	URI                  string            `json:"uri"`
	VoiceApplicationSid  string            `json:"voice_application_sid"`
	VoiceCallerIDLookup  bool              `json:"voice_caller_id_lookup"`
	VoiceFallbackMethod  string            `json:"voice_fallback_method"`
	VoiceFallbackURL     string            `json:"voice_fallback_url"`
	VoiceMethod          string            `json:"voice_method"`
	VoiceURL             string            `json:"voice_url"`
}

type NumberCapability struct {
	MMS   bool `json:"mms"`
	SMS   bool `json:"sms"`
	Voice bool `json:"voice"`
}

func (twilio *Twilio) GetPhoneNumbers() (incomingPhoneNumbers *IncomingPhoneNumbers, exception *Exception, err error) {
	servicesUrl := twilio.BaseUrl + "/Accounts/" + twilio.AccountSid + "/IncomingPhoneNumbers.json"
	servicesUrl = setPageSize(servicesUrl, 100)

	out, exception, err := twilio.getResponseBody(servicesUrl)

	if exception != nil || err != nil {
		return nil, exception, err
	}

	incomingPhoneNumbers = new(IncomingPhoneNumbers)
	err = json.Unmarshal(out, incomingPhoneNumbers)

	return incomingPhoneNumbers, exception, err
}
