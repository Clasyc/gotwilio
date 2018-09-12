package gotwilio

import (
	"time"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"net/http"
)

type CopilotService struct {
	AccountSID            string        `json:"account_sid"`
	SID                   string        `json:"sid"`
	DateCreated           time.Time     `json:"date_created"`
	DateUpdated           time.Time     `json:"date_updated"`
	FriendlyName          string        `json:"friendly_name"`
	InboundRequestURL     string        `json:"inbound_request_url"`
	InboundMethod         string        `json:"inbound_method"`
	FallbackURL           string        `json:"fallback_url"`
	FallbackMethod        string        `json:"fallback_method"`
	StatusCallback        string        `json:"status_callback"`
	StickySender          bool          `json:"sticky_sender"`
	SmartEncoding         bool          `json:"smart_encoding"`
	MMSConverter          bool          `json:"mms_converter"`
	FallbackToLongCode    bool          `json:"fallback_to_long_code"`
	ScanMessageContent    string        `json:"scan_message_content"`
	AreaCodeGeomatch      bool          `json:"area_code_geomatch"`
	ValidityPeriod        int           `json:"validity_period"`
	SynchronousValidation bool          `json:"synchronous_validation"`
	Links                 Links         `json:"links"`
	URL                   string        `json:"url"`
	PhoneNumbers          []PhoneNumber `json:"phone_numbers"`
	AlphaSender           *AlphaSender
}

type CopilotServiceList struct {
	Meta Meta `json:"meta"`
	Services []CopilotService `json:"services"`
}

type Links struct {
	PhoneNumbers string `json:"phone_numbers"`
	ShortCodes   string `json:"short_codes"`
	AlphaSenders string `json:"alpha_senders"`
}

type PhoneNumber struct {
	SID            string       `json:"sid"`
	AccountSID     string       `json:"account_sid"`
	ServiceSID     string       `json:"service_sid"`
	DateCreated    time.Time    `json:"date_created"`
	DateUpdated    time.Time    `json:"date_updated"`
	PhoneNumber    string       `json:"phone_number"`
	CountryCode    string       `json:"country_code"`
	Capabilities   []string `json:"capabilities"`
	URL            string       `json:"url"`
	PhoneNumberSID string       `json:"phone_number_sid"`
}

type CreateServiceOptions struct {
	FriendlyName string `json:"FriendlyName"`
    AreaCodeGeomatch bool `json:"AreCodeGeomatch"`
    FallbackMethod string `json:"FallbackMethod"`
    FallbackToLongCode bool `json:"FallbackToLongCode"`
    FallbackUrl string `json:"FallbackUrl"`
    InboundMethod string `json:"InboundMethod"`
    InboundRequestUrl string `json:"InboundRequestUrl"`
    MmsConverter bool `json:"MmsConverter"`
    SmartEncoding bool `json:"SmartEncoding"`
    StatusCallback string `json:"StatusCallback"`
    StickySender bool `json:"StickySender"`
    ValidityPeriod int `json:"ValidityPeriod"`
}

type AlphaSenderList struct {
	Meta Meta `json:"meta"`
	AlphaSenders []AlphaSender `json:"alpha_senders"`
}

type AlphaSender struct {
	Sid string `json:"sid"`
	AccountSid string `json:"account_sid"`
	ServiceSid string `json:"service_sid"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	AlphaSender string `json:"alpha_sender"`
	Capabilities []string `json:"capabilities"`
	Url string `json:"url"`
 }

 func (l *CopilotService) fetchAlphaSender(twilio *Twilio) (alphaSender *AlphaSender, exception *Exception, err error) {
	 out, exception, err := twilio.getResponseBody(l.Links.AlphaSenders)

	 if exception != nil || err != nil {
		 return nil, exception, err
	 }

	 alphaSenderList := new(AlphaSenderList)
	 err = json.Unmarshal(out, alphaSenderList)

	 if len(alphaSenderList.AlphaSenders) == 1 {
		 l.AlphaSender = &alphaSenderList.AlphaSenders[0]
	 }

	 return l.AlphaSender, exception, err
 }

func (twilio *Twilio) GetService(sid string) (copilotResponse *CopilotService, exception *Exception, err error) {
	servicesUrl := twilio.MessagingUrl + "/Services/" + sid

	out, exception, err := twilio.getResponseBody(servicesUrl)

	if exception != nil || err != nil {
		return nil, exception, err
	}

	copilotResponse = new(CopilotService)
	err = json.Unmarshal(out, copilotResponse)

	return copilotResponse, exception, err
}

func (twilio *Twilio) GetServicesWithAlphaSenders() (copilotServiceList *CopilotServiceList, exception *Exception, err error) {
	servicesUrl := twilio.MessagingUrl + "/Services"

	copilotServiceList, exception, err = twilio.getServices(servicesUrl)

	if err != nil || exception != nil {
		return copilotServiceList, exception, err
	}

	serviceList := new(CopilotServiceList)

	for k := range copilotServiceList.Services {
		_, exception, err := copilotServiceList.Services[k].fetchAlphaSender(twilio)

		if err != nil || exception != nil {
			return copilotServiceList, exception, err
		}
	}

	serviceList.Services = append(serviceList.Services, copilotServiceList.Services...)

	for len(copilotServiceList.Services) == copilotServiceList.Meta.PageSize {
		copilotServiceList, exception, err = copilotServiceList.nextPage(twilio)
		list := copilotServiceList
		s := list.Services

		if err != nil || exception != nil {
			return copilotServiceList, exception, err
		}

		for k := range s {
			_, exception, err := s[k].fetchAlphaSender(twilio)

			if err != nil || exception != nil {
				return copilotServiceList, exception, err
			}
		}

		if copilotServiceList.Meta.URL != copilotServiceList.Meta.FirstPageURL {
			serviceList.Services = append(serviceList.Services, s...)
		}
	}

	return serviceList, exception, err
}

func (twilio *Twilio) GetServices() (copilotServiceList *CopilotServiceList, exception *Exception, err error) {
	servicesUrl := twilio.MessagingUrl + "/Services"

	copilotServiceList, exception, err = twilio.getServices(servicesUrl)

	if err != nil || exception != nil {
		return copilotServiceList, exception, err
	}

	serviceList := new(CopilotServiceList)
	serviceList.Services = append(serviceList.Services, copilotServiceList.Services...)

	for len(copilotServiceList.Services) == copilotServiceList.Meta.PageSize {
		copilotServiceList, exception, err = copilotServiceList.nextPage(twilio)
		list := copilotServiceList
		s := list.Services

		if copilotServiceList.Meta.URL != copilotServiceList.Meta.FirstPageURL {
			serviceList.Services = append(serviceList.Services, s...)
		}
	}

	return serviceList, exception, err
}

func (twilio *Twilio) getServices(servicesUrl string) (copilotServiceList *CopilotServiceList, exception *Exception, err error) {
	out, exception, err := twilio.getResponseBody(servicesUrl)

	if exception != nil || err != nil {
		return nil, exception, err
	}

	copilotServiceList = new(CopilotServiceList)
	err = json.Unmarshal(out, copilotServiceList)

	return copilotServiceList, exception, err
}

func (twilio *Twilio) AddPhoneNumber() {

}

func (twilio *Twilio) CreateService(options *CreateServiceOptions) (copilotResponse *CopilotService, exception *Exception, err error) {
	q := url.Values{}

	setUrlValues(options, &q)

	twilioUrl := twilio.MessagingUrl + "/Services"
	res, err := twilio.post(q, twilioUrl)

	if err != nil {
		return copilotResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return copilotResponse, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return copilotResponse, exception, err
	}

	copilotResponse = new(CopilotService)
	err = json.Unmarshal(responseBody, copilotResponse)
	return copilotResponse, exception, err
}

func (twilio *Twilio) DeleteService(sid string) (exception *Exception, err error) {
	servicesUrl := twilio.MessagingUrl + "/Services/" + sid

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

func (c *CopilotService) DeletePhoneNumber(sid string, twilio *Twilio) (exception *Exception, err error) {
	servicesUrl := twilio.MessagingUrl + "/Services/" + c.SID + "/PhoneNumbers/" + sid

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

func (c *CopilotService) AddAlphaSenderID(id string, twilio *Twilio) (alphaSender *AlphaSender, exception *Exception, err error) {
	q := url.Values{}

	q.Set("AlphaSender", id)

	twilioUrl := twilio.MessagingUrl + "/Services/" + c.SID + "/AlphaSenders"
	res, err := twilio.post(q, twilioUrl)

	if err != nil {
		return alphaSender, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return alphaSender, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return alphaSender, exception, err
	}

	alphaSender = new(AlphaSender)
	err = json.Unmarshal(responseBody, alphaSender)

	return alphaSender, nil, err
}

func (c *CopilotService) AddPhoneNumber(sid string, twilio *Twilio) (phoneNumber *PhoneNumber, exception *Exception, err error) {
	q := url.Values{}

	q.Set("PhoneNumberSid", sid)

	twilioUrl := twilio.MessagingUrl + "/Services/" + c.SID + "/PhoneNumbers"
	res, err := twilio.post(q, twilioUrl)

	if err != nil {
		return phoneNumber, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return phoneNumber, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return phoneNumber, exception, err
	}

	phoneNumber = new(PhoneNumber)
	err = json.Unmarshal(responseBody, phoneNumber)

	return phoneNumber, nil, err
}

func setUrlValues(options interface{}, q *url.Values) {
	v := reflect.ValueOf(options).Elem()

	for i := 0; i < v.NumField(); i++ {
		varName := v.Type().Field(i).Name
		varType := v.Type().Field(i).Type.String()
		varValue := v.Field(i).Interface()

		if varValue != nil {
			switch va := varType; va {
			case "string":
				if varValue.(string) != "" {
					val := varValue.(string)
					q.Set(varName, val)
				}
			case "int":
				if varValue.(int) != 0 {
					val := strconv.Itoa(varValue.(int))
					q.Set(varName, val)
				}
			case "bool":
				val := strconv.FormatBool(varValue.(bool))
				q.Set(varName, val)
			}
		}
	}
}

func (l *CopilotServiceList) nextPage(twilio *Twilio) (copilotServiceList *CopilotServiceList, exception *Exception, err error) {
		return twilio.getServices(l.Meta.NextPageURL.(string))
}

func (l *CopilotServiceList) previousPage(twilio *Twilio) (copilotServiceList *CopilotServiceList, exception *Exception, err error) {
	return twilio.getServices(l.Meta.PreviousPageURL.(string))
}