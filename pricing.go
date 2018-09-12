package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type PricingResponse struct {
	Country           string             `json:"country"`
	ISOCountry        string             `json:"iso_country"`
	PriceUnit         string             `json:"price_unit"`
	OutboundSMSPrices []OutboundSMSPrice `json:"outbound_sms_prices"`
	InboundSMSPrice   []Price            `json:"inbound_sms_prices"`
	Url               string             `json:"uri"`
}

type OutboundSMSPrice struct {
	MCC     string  `json:"mcc"`
	MNC     string  `json:"mnc"`
	Carrier string  `json:"carrier"`
	Prices  []Price `json:"prices"`
}

type Price struct {
	NumberType   string  `json:"number_type"`
	BasePrice    float64 `json:"base_price,string"`
	CurrentPrice float64 `json:"current_price,string"`
}

type PricingList struct {
	Meta Meta `json:"meta"`
	Countries []Country `json:"countries"`
}

type Country struct {
	Country string `json:"country"`
	IsoCountry string `json:"iso_country"`
	Url string `json:"url"`
}

func (twilio *Twilio) GetPricing(countryISO string) (pricingResponse *PricingResponse, exception *Exception, err error) {
	pricingUrl := twilio.PricingUrl + "/Messaging/Countries/" + strings.ToUpper(countryISO)

	res, err := twilio.get(pricingUrl)
	if err != nil {
		return pricingResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return pricingResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return pricingResponse, exception, err
	}

	pricingResponse = new(PricingResponse)
	err = json.Unmarshal(responseBody, pricingResponse)

	return pricingResponse, exception, err
}

func (twilio *Twilio) GetCountriesPricing() (pricingList *PricingList, exception *Exception, err error) {
	pricingUrl := twilio.PricingUrl + "/Messaging/Countries"
	pricingUrl = setPageSize(pricingUrl, 250)

	res, err := twilio.get(pricingUrl)
	if err != nil {
		return pricingList, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return pricingList, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return pricingList, exception, err
	}

	pricingList = new(PricingList)
	err = json.Unmarshal(responseBody, pricingList)

	return pricingList, exception, err
}