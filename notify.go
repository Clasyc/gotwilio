package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
	"net/url"
)

type NotifyService struct {
	SID                                     string      `json:"sid"`
	AccountSID                              string      `json:"account_sid"`
	FriendlyName                            string      `json:"friendly_name"`
	DateCreated                             time.Time   `json:"date_created"`
	DateUpdated                             time.Time   `json:"date_updated"`
	ApnCredentialSID                        string      `json:"apn_credential_sid"`
	GcmCredentialSID                        string      `json:"gcm_credential_sid"`
	FcmCredentialSID                        string      `json:"fcm_credential_sid"`
	MessageServiceSID                       string      `json:"message_service_sid"`
	FacebookMessengerPageId                 string      `json:"facebook_messenger_page_id"`
	DefaultApnNotificationProtocolVersion   string      `json:"default_apn_notification_protocol_version"`
	DefaultGcmNotificationProtocolVersion   string      `json:"default_gcm_notification_protocol_version"`
	DefaultFcmNotificationProtocolVersion   string      `json:"default_fcm_notification_protocol_version"`
	LogEnabled                              bool        `json:"log_enabled"`
	URL                                     string      `json:"url"`
	Links                                   NotifyLinks `json:"links"`
	AlexaSkillId                            string      `json:"alexa_skill_id"`
	DefaultAlexaNotificationProtocolVersion string      `json:"default_alexa_notification_protocol_version"`
}

type NotifyLinks struct {
	PhoneNumbers string `json:"phone_numbers"`
	ShortCodes   string `json:"short_codes"`
	AlphaSenders string `json:"alpha_senders"`
}

type CreateNotifyServiceOptions struct {
	FriendlyName                          string `json:"friendly_name"`
	ApnCredentialSID                      string `json:"apn_credential_sid"`
	GcmCredentialSID                      string `json:"gcm_credential_sid"`
	MessagingServiceSid                   string `json:"messaging_service_sid"`
	FacebookMessengerPageId               string `json:"facebook_messenger_page_id"`
	DefaultApnNotificationProtocolVersion string `json:"default_apn_notification_protocol_version"`
	DefaultGcmNotificationProtocolVersion string `json:"default_gcm_notification_protocol_version"`
	FcmCredentialSid                      string `json:"fcm_default_sid"`
	DefaultFcmNotificationProtocolVersion string `json:"default_fcm_notification_protocol_version"`
}

//TODO: create fetch and fetch all functions

func (twilio *Twilio) CreateNotifyService(options *CreateNotifyServiceOptions) (notifyService *NotifyService,
	exception *Exception, err error) {
	return twilio.postNotifyService("/Services", options)
}

func (twilio *Twilio) UpdateNotifyService(serviceSID string, options *CreateNotifyServiceOptions) (notifyService *NotifyService,
	exception *Exception, err error) {
	return twilio.postNotifyService(fmt.Sprintf("/Services/%s", serviceSID), options)
}

func (twilio *Twilio) DeleteNotifyService(sid string) (exception *Exception, err error) {
	servicesUrl := twilio.NotifyUrl + "/Services/" + sid

	return twilio.DeleteResource(servicesUrl)
}

func (twilio *Twilio) postNotifyService(path string, options *CreateNotifyServiceOptions) (notifyService *NotifyService,
	exception *Exception, err error) {
	q := url.Values{}

	setUrlValues(options, &q)

	twilioUrl := twilio.NotifyUrl + path
	res, err := twilio.post(q, twilioUrl)

	if err != nil {
		return notifyService, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return notifyService, exception, err
	}

	if !(res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusOK) {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)
		return notifyService, exception, err
	}

	notifyService = new(NotifyService)
	err = json.Unmarshal(responseBody, notifyService)
	return notifyService, exception, err
}
