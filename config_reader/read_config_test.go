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

func TestProcessConfigNoTuyaDeviceConfigs(t *testing.T) {
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

func TestProcessConfigNoDeviceName(t *testing.T) {
	os.Setenv("ALARM_MANAGER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_device_name/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without tuya devices should fail.")
	} else {
		if err.Error() != "Fatal error config: device home_alarm has no name." {
			t.Errorf("Error should be \"Fatal error config: device home_alarm has no name.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigRepeatedDeviceName(t *testing.T) {
	os.Setenv("ALARM_MANAGER_CONFIG_FILE_LOCATION", "./config_files_test/config_repeated_device_name/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without tuya devices should fail.")
	} else {
		if err.Error() != "Fatal error config: device name 'Home Alarm' is repeated." {
			t.Errorf("Error should be \"Fatal error config: device name 'Home Alarm' is repeated.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigRepeatedDeviceID(t *testing.T) {
	os.Setenv("ALARM_MANAGER_CONFIG_FILE_LOCATION", "./config_files_test/config_repeated_device_id/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without tuya devices should fail.")
	} else {
		if err.Error() != "Fatal error config: device ID device123 is repeated." {
			t.Errorf("Error should be \"Fatal error config: device ID device123 is repeated.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoHTTPPort(t *testing.T) {
	os.Setenv("ALARM_MANAGER_CONFIG_FILE_LOCATION", "./config_files_test/config_no_web_server/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without web server config should fail.")
	} else {
		if err.Error() != "Fatal error config: no web_server port was found." {
			t.Errorf("Error should be \"Fatal error config: no web_server port was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigOK(t *testing.T) {
	os.Setenv("ALARM_MANAGER_CONFIG_FILE_LOCATION", "./config_files_test/config_ok/")
	_, err := ReadConfig()
	if err != nil {
		t.Errorf("ReadConfig method without tuya devices should not fail.")
	}
}
