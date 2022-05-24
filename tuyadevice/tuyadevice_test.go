package tuyadevice

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
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

	device := TuyaDevice{Name: "Test"}

	tokenError := device.RetrieveToken(client)

	if tokenError != nil {
		t.Errorf("Token retrievement should not fail. Error was %s", tokenError)
	}

	if device.Token != "testtoken" {
		t.Errorf("Retrived token should be testtoken, not %s.", device.Token)
	}

}

func TestGetTokenFailed(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593`))}}}

	device := TuyaDevice{Name: "Test"}

	tokenError := device.RetrieveToken(client)

	if tokenError == nil {
		t.Errorf("Token retrievement should fail.")
	}

}

func TestGetDeviceInfo(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"disarmed"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"}`))}}}

	device := TuyaDevice{Name: "Test"}

	_, deviceInfoErr := device.GetDeviceInfo(client)

	if deviceInfoErr != nil {
		t.Errorf("Device info retrievement should not fail. Error was %s", deviceInfoErr)
	}

}

func TestValidateInvalidDevice(t *testing.T) {

	device := TuyaDevice{Name: "Test"}
	validationErr := device.Validate()
	if validationErr == nil {
		t.Errorf("Device Validation should fail.")
	}

}

func TestValidateValidDevice(t *testing.T) {

	device := TuyaDevice{Name: "Test", Host: "host.io", ClientID: "clientid", Secret: "secret", DeviceID: "deviceid", DeviceType: "alarm"}
	validationErr := device.Validate()
	if validationErr != nil {
		t.Errorf("Device Validation shouldn't fail. Error was %s", validationErr)
	}

}

func TestGetDeviceType(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"disarmed"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"}`))}}}

	device := TuyaDevice{Name: "Test", DeviceType: "TestType"}

	device.GetDeviceInfo(client)

	deviceType := device.GetDeviceType()
	if deviceType == "" {
		t.Errorf("Device type should not be an empty string.")
	}

}

func TestGetDeviceName(t *testing.T) {

	device := TuyaDevice{Name: "Test", DeviceType: "TestType"}

	deviceName := device.GetDeviceName()
	if deviceName != "Test" {
		t.Errorf("Device type should be 'Test', not '%s'.", deviceName)
	}

}

func TestChangeMode(t *testing.T) {

	device := TuyaDevice{Name: "Test", Host: "host.io", ClientID: "clientid", Secret: "secret", DeviceID: "deviceid", DeviceType: "alarm"}

	clientChangueModeResponse := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":true,"success":true,"t":1653184890385,"tid":"18dd6963d97311eca734f2b4cd1fee5a"}`))}}}
	changeModeErr := device.ChangeMode(clientChangueModeResponse, "newmode")
	if changeModeErr != nil {
		t.Errorf("Device mode change shouldn't fail. Error was %s", changeModeErr)
	}

}

func TestChangeModeFailed(t *testing.T) {

	device := TuyaDevice{Name: "Test", Host: "host.io", ClientID: "clientid", Secret: "secret", DeviceID: "deviceid", DeviceType: "alarm"}

	clientChangueModeResponse := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"code":1109,"msg":"param is illegal ,please check it","success":false,"t":1653182849837,"tid":"589da661d96e11eca914e276ec45657f"}`))}}}
	changeModeErr := device.ChangeMode(clientChangueModeResponse, "newmode")
	if changeModeErr == nil {
		t.Errorf("Device mode change should fail.")
	}

}

func TestChangeModeFailedBecauseJson(t *testing.T) {

	device := TuyaDevice{Name: "Test", Host: "host.io", ClientID: "clientid", Secret: "secret", DeviceID: "deviceid", DeviceType: "alarm"}

	clientChangueModeResponse := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"code":1109,"msg":"param is illegal ,please check it","success":false,"t":1653182849837,"tid":"589da661d96e11eca914e276ec45657f"`))}}}
	changeModeErr := device.ChangeMode(clientChangueModeResponse, "newmode")
	if changeModeErr == nil {
		t.Errorf("Device mode change should fail.")
	}

}
