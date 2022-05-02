package devices

import (
	"os"
	"testing"

	config "github.com/a-castellano/AlarmManager/config_reader"
)

func TestGetToken(t *testing.T) {
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
