package devices

import (
	"os"
	"testing"

	config "github.com/a-castellano/AlarmManager/config_reader"
	tuyadevice "github.com/a-castellano/AlarmManager/tuyadevice"
)

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
