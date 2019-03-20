package gotwilio

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var paramsL map[string]string

func init() {
	paramsL = make(map[string]string)

	// Only LIVE credentials possible, because of 20008 error
	paramsL["SID"] = ""
	paramsL["TOKEN"] = ""
}

func TestCreateCopilotService(t *testing.T) {
	twilio := NewTwilioClient(paramsL["SID"], paramsL["TOKEN"])

	for i := 0; i < 55; i++ {
		_, exc, err := twilio.CreateCopilotService(&CreateCopilotServiceOptions{
			FriendlyName: "ServiceTest_" + strconv.Itoa(i),
		})

		if err != nil {
			t.Fatal(err)
		}

		if exc != nil {
			t.Fatal(exc)
		}
	}
}

func TestGetCopilotService(t *testing.T) {
	twilio := NewTwilioClient(paramsL["SID"], paramsL["TOKEN"])
	c, exc, err := twilio.GetCopilotServices()

	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

	for _, s := range c.Services {
		if strings.HasPrefix(s.FriendlyName, "ServiceTest_") {
			fmt.Println("Getting service services")

			se, exc, err := twilio.GetCopilotService(s.SID)

			if err != nil {
				t.Fatal(err)
			}

			if exc != nil {
				t.Fatal(exc)
			}

			fmt.Println("Adding alpha sender id")

			_, exc, err = se.AddAlphaSenderID("testas", twilio)

			if err != nil {
				t.Fatal(err)
			}

			if exc != nil {
				t.Fatal(exc)
			}
		}
	}
}

func TestGetCopilotServices(t *testing.T) {
	twilio := NewTwilioClient(paramsL["SID"], paramsL["TOKEN"])
	c, exc, err := twilio.GetCopilotServices()

	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

	fmt.Println(len(c.Services))

	for _, s := range c.Services {
		if strings.HasPrefix(s.FriendlyName, "ServiceTest_") {
			fmt.Println("Deleting services")

			exc, err = twilio.DeleteCopilotService(s.SID)

			if err != nil {
				t.Fatal(err)
			}

			if exc != nil {
				t.Fatal(exc)
			}
		}
	}
}

func TestAddPhoneNumber(t *testing.T) {
	twilio := NewTwilioClient(paramsL["SID"], paramsL["TOKEN"])

	se, exc, err := twilio.GetCopilotService("MG209524e859f49f2302a0d3ff87876d3d")

	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

	fmt.Println("Adding phone number")

	_, exc, err = se.AddPhoneNumber("PN86ae1a9276d37f32a0cfc21a2331258d", twilio)

	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestDeletePhoneNumber(t *testing.T) {
	twilio := NewTwilioClient(paramsL["SID"], paramsL["TOKEN"])

	se, exc, err := twilio.GetCopilotService("MG209524e859f49f2302a0d3ff87876d3d")

	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

	fmt.Println("Deleting phone number")

	exc, err = se.DeletePhoneNumber("PN86ae1a9276d37f32a0cfc21a2331258d", twilio)

	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}
}
