package devices

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	config "github.com/a-castellano/AlarmManager/config_reader"
	tuyadevice "github.com/a-castellano/AlarmManager/tuyadevice"
)

type RoundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *RoundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

func TestCreateTuyaDevice(t *testing.T) {
	os.Setenv("ALARM_MANAGER_CONFIG_FILE_LOCATION", "../config_reader/config_files_test/config_ok/")
	devicesConfig, readConfigErr := config.ReadConfig()
	if readConfigErr != nil {
		t.Errorf("ReadConfig method without tuya devices should not fail, error was '%s'.", readConfigErr)
	}
	deviceConfig := devicesConfig.Devices["Home Alarm"]
	device := CreateTuyaDeviceFromConfig(deviceConfig)

	if device.DeviceID != "device123" {
		t.Errorf("Processed device idshould be 'device123', not %s.", device.DeviceID)
	}

}

func TestAddOneDevice(t *testing.T) {

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}
	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"

	deviceRef := &device
	addErr := deviceManager.AddDevice(deviceRef)
	if addErr != nil {
		t.Errorf("AddDevice should not fail with only one device, error was '%s'.", addErr)
	}
}

func TestAddTwoDevicesWithSameName(t *testing.T) {

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}
	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"

	var device2 tuyadevice.TuyaDevice
	device2.Name = "Test Device"

	deviceRef := &device
	device2Ref := &device2
	deviceManager.AddDevice(deviceRef)
	addErr := deviceManager.AddDevice(device2Ref)
	if addErr == nil {
		t.Errorf("AddDevice should fail with two devices witch the same name.")
	}
}

func TestStart(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	startError := deviceManager.Start(client)

	if startError != nil {
		t.Errorf("Device Manager start should not fail. Error was %s", startError)
	}

}

func TestStartFailedToken(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	startError := deviceManager.Start(client)

	if startError == nil {
		t.Errorf("Device Manager start should fail.")
	}

}

func TestRetrieveInfo(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	clientRetrieveInfo := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"disarmed"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"}`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "99AST"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	retrieveInfoError := deviceManager.RetrieveInfo(clientRetrieveInfo)
	if retrieveInfoError != nil {
		t.Errorf("Device Manager info retrieval should not fail. Error was %s", retrieveInfoError)
	}

}

func TestRetrieveInfoFailed(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	clientRetrieveInfo := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"disarmed"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "99AST"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	retrieveInfoError := deviceManager.RetrieveInfo(clientRetrieveInfo)
	if retrieveInfoError == nil {
		t.Errorf("Device Manager info retrieval should fail.")
	}

}

func TestRetrieveInfoFailedBecauseUnknownDeviceType(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	clientRetrieveInfo := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"disarmed"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "unknown"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	retrieveInfoError := deviceManager.RetrieveInfo(clientRetrieveInfo)
	if retrieveInfoError.Error() != "Alarm Test Device type unknown not supported" {
		t.Errorf("Device Manager info retrieval should fail with error 'Alarm Test Device type unknown not supported', error was '%s'", retrieveInfoError)
	}

}

func TestRetrieveInfoAlarmDisarmed(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	clientRetrieveInfo := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"disarmed"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"}`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "99AST"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	deviceManager.RetrieveInfo(clientRetrieveInfo)

	alarmInfo := deviceManager.AlarmsInfo["Test Device"]

	if alarmInfo.ShowInfo().Firing != false {
		t.Errorf("Alarm shouldn't be firing.")
	}

	if alarmInfo.ShowInfo().Mode != Disarmed {
		t.Errorf("Alarm should be Disarmed.")
	}

}

func TestRetrieveInfoAlarmHome(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	clientRetrieveInfo := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"home"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"}`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "99AST"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	deviceManager.RetrieveInfo(clientRetrieveInfo)

	alarmInfo := deviceManager.AlarmsInfo["Test Device"]

	if alarmInfo.ShowInfo().Firing != false {
		t.Errorf("Alarm shouldn't be firing.")
	}

	if alarmInfo.ShowInfo().Mode != HomeArmed {
		t.Errorf("Alarm should be HomeArmed.")
	}

}

func TestRetrieveInfoAlarmArm(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	clientRetrieveInfo := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"arm"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"}`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "99AST"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	deviceManager.RetrieveInfo(clientRetrieveInfo)

	alarmInfo := deviceManager.AlarmsInfo["Test Device"]

	if alarmInfo.ShowInfo().Firing != false {
		t.Errorf("Alarm shouldn't be firing.")
	}

	if alarmInfo.ShowInfo().Mode != FullyArmed {
		t.Errorf("Alarm should be FullyArmed.")
	}

}

func TestRetrieveInfoAlarmArmFiring(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	clientRetrieveInfo := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"arm"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"alarm"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"}`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "99AST"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	deviceManager.RetrieveInfo(clientRetrieveInfo)

	alarmInfo := deviceManager.AlarmsInfo["Test Device"]

	if alarmInfo.ShowInfo().Firing != true {
		t.Errorf("Alarm should be firing.")
	}

	if alarmInfo.ShowInfo().Mode != FullyArmed {
		t.Errorf("Alarm should be FullyArmed.")
	}

}

func TestChangeAlarmModeNonRetrievedInfo(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "99AST"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	clientChangueModeResponse := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}
	changeModeError := deviceManager.ChangeMode(clientChangueModeResponse, "NonExistentDevice", "NonExistentMode")
	if changeModeError == nil {
		t.Errorf("Device Manager change mode should fail.")
	} else {
		if changeModeError.Error() != "Device has not retrieved devices info yet." {
			t.Errorf("Device Manager change mode error should be 'Device has not retrieved devices info yet.' instead of '%s'", changeModeError)
		}
	}

}

func TestChangeAlarmModeNonExistentDevice(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	clientRetrieveInfo := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"disarmed"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"}`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "99AST"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	deviceManager.RetrieveInfo(clientRetrieveInfo)

	clientChangueModeResponse := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}
	changeModeError := deviceManager.ChangeMode(clientChangueModeResponse, "NonExistentDevice", "NonExistentMode")
	if changeModeError == nil {
		t.Errorf("Device Manager change mode should fail.")
	} else {
		if changeModeError.Error() != "Device 'NonExistentDevice' is not a managed device." {
			t.Errorf("Device Manager change mode error should be 'Device 'NonExistentDevice' is not a managed device.' instead of '%s'", changeModeError)
		}
	}

}

func TestChangeAlarmModeNonExistentMode(t *testing.T) {

	clientGetToken := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}

	clientRetrieveInfo := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"active_time":1634987857,"biz_type":18,"category":"mal","create_time":1620050314,"icon":"smart/icon/ay15427647462366edzT/153535979f068afab73c91841c844c82.png","id":"1234456789cca88fafe1","ip":"199.46.115.128","lat":"37.9988","local_key":"bc10cf0dca9aa13f","lon":"-5.0338","model":"99AST-西语","name":"Multifunction alarm","online":true,"owner_id":"11154007","product_id":"2aelhoqe23e7vxjr","product_name":"Multifunction alarm ","status":[{"code":"master_mode","value":"disarmed"},{"code":"delay_set","value":0},{"code":"alarm_time","value":1},{"code":"switch_alarm_sound","value":true},{"code":"switch_alarm_light","value":false},{"code":"switch_mode_sound","value":true},{"code":"switch_mode_light","value":true},{"code":"switch_kb_sound","value":true},{"code":"switch_kb_light","value":true},{"code":"password_set","value":""},{"code":"charge_state","value":true},{"code":"switch_low_battery","value":false},{"code":"alarm_call_number","value":"AQkAAQ=="},{"code":"alarm_sms_number","value":""},{"code":"switch_alarm_call","value":true},{"code":"switch_alarm_sms","value":true},{"code":"telnet_state","value":"sim_card_no"},{"code":"zone_attribute","value":"disarmed"},{"code":"muffling","value":false},{"code":"alarm_msg","value":"AEEAUABQACAARABlAHMAZQByAG0AYQBkAG8="},{"code":"alarm_delay_time","value":0},{"code":"switch_mode_dl_sound","value":false},{"code":"master_state","value":"normal"},{"code":"master_information","value":""},{"code":"factory_reset","value":false},{"code":"night_light_bright","value":1},{"code":"sub_class","value":"detector"},{"code":"sub_type","value":"motion_sensor"},{"code":"sub_admin","value":"CEAFEQH///8OAHAAYQBzAGkAbABsAG8="},{"code":"sub_state","value":"normal"}],"sub":false,"time_zone":"+01:00","uid":"eujJ01152904a15dpPln","update_time":1639405182,"uuid":"1531440084cca88fafe1"},"success":true,"t":1645128085588,"tid":"62fa5cb3902c11eceec15ef357c3f603"}`))}}}

	deviceManager := DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]Alarm)}

	var device tuyadevice.TuyaDevice
	device.Name = "Test Device"
	device.DeviceType = "99AST"

	deviceRef := &device
	deviceManager.AddDevice(deviceRef)

	deviceManager.Start(clientGetToken)

	deviceManager.RetrieveInfo(clientRetrieveInfo)

	clientChangueModeResponse := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`{"result":{"access_token":"testtoken","expire_time":7200,"refresh_token":"refesh","uid":"bay1635003708553hilW"},"success":true,"t":1644740470593}`))}}}
	changeModeError := deviceManager.ChangeMode(clientChangueModeResponse, "Test Device", "NonExistentMode")

	if changeModeError == nil {
		t.Errorf("Device Manager change mode should fail bacause mode is 'NonExistentMode'.")
	} else {
		if changeModeError.Error() != "Alarm mode 'NonExistentMode' is not defined." {
			t.Errorf("Device Manager change mode error should be 'Alarm mode 'NonExistentMode' is not defined.' but error was '%s'.", changeModeError)
		}
	}
}
