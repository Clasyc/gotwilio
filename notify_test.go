package gotwilio

import (
	"testing"
	"os"
)

var notifyTestParams map[string]string

func init() {
	notifyTestParams = make(map[string]string)

	// Only LIVE credentials possible, because of 20008 error
	notifyTestParams["SID"] = os.Getenv("TEST_SID")
	notifyTestParams["TOKEN"] = os.Getenv("TEST_TOKEN")
}

// service CURL operations needs to run sequentially, create -> update -> delete
func TestNotify(t *testing.T) {
	t.Run("CreateService", func(t *testing.T) {
		twilio := NewTwilioClient(notifyTestParams["SID"], notifyTestParams["TOKEN"])

		t.Log("Service succesfully created!")

		s, exc, err := twilio.CreateNotifyService(&CreateNotifyServiceOptions{
			FriendlyName: "NotifyServiceTest_02",
			MessagingServiceSid: "MGa98ebb11e8ae98a7f699cad20f5ec347",
		})

		if err != nil {
			t.Fatal(err)
		}

		if exc != nil {
			t.Fatal(exc)
		}

		notifyTestParams["SERVICE_TO_UPDATE"] = s.SID
	})

	t.Run("UpdateService", func(t *testing.T) {
		twilio := NewTwilioClient(notifyTestParams["SID"], notifyTestParams["TOKEN"])

		t.Log("Service succesfully updated!")

		_, exc, err := twilio.UpdateNotifyService(notifyTestParams["SERVICE_TO_UPDATE"], &CreateNotifyServiceOptions{
			FriendlyName: "NotifyServiceTest_Update",
		})

		if err != nil {
			t.Fatal(err)
		}

		if exc != nil {
			t.Fatal(exc)
		}
	})

	t.Run("DeleteService", func(t *testing.T) {
		twilio := NewTwilioClient(notifyTestParams["SID"], notifyTestParams["TOKEN"])

		t.Log("Service succesfully deleted!")

		exc, err := twilio.DeleteNotifyService(notifyTestParams["SERVICE_TO_UPDATE"])

		if err != nil {
			t.Fatal(err)
		}

		if exc != nil {
			t.Fatal(exc)
		}
	})
}