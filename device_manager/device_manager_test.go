package devices

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	config "github.com/a-castellano/AlarmManager/config_reader"
)

type RoundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *RoundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

func TestGetToken(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	device := config.TuyaDevice{Name: "Test", DeviceType: "99AST", Host: "https://test.windmaker.net", ClientID: "test", Secret: "test", DeviceID: "test"}

	token, tokenError := GetToken(client, device)

	if tokenError != nil {
		t.Errorf("Token retrievement should not fail. Error was %s", tokenError)
	}

	if token != "testtoken" {
		t.Errorf("Retrived token should be testtoken, not %s.", token)
	}

}
