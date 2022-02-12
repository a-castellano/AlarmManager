package config

import (
	"os"
	"testing"
)

func TestProcessNoConfigFilePresent(t *testing.T) {

	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without any valid config file should fail.")
	} else {
		if err.Error() != "Environment variable ALARM_MANAGER_CONFIG_FILE_LOCATION is not defined." {
			t.Errorf("Error should be 'Environment variable ALARM_MANAGER_CONFIG_FILE_LOCATION is not defined.', but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoTuyaDevices(t *testing.T) {
	os.Setenv("ALARM_MANAGER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_tuya_devices/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without tuya devices should fail.")
	} else {
		if err.Error() != "Fatal error config: no tuya_devices field was found." {
			t.Errorf("Error should be \"Fatal error config: no tuya_devices field was found.\" but error was '%s'.", err.Error())
		}
	}
}
