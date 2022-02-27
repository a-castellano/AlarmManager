package config

import (
	"errors"
	"reflect"

	viperLib "github.com/spf13/viper"
)

type TuyaDeviceConfig struct {
	Name       string
	DeviceType string
	Host       string
	ClientID   string
	Secret     string
	DeviceID   string
}

type Config struct {
	Devices map[string]TuyaDeviceConfig
}

func ReadConfig() (Config, error) {
	var configFileLocation string
	var config Config

	var envVariable string = "ALARM_MANAGER_CONFIG_FILE_LOCATION"

	requiredVariables := []string{"tuya_devices"}

	tuyaDevicesRequiredVariables := []string{"name", "type", "host", "client_id", "secret", "device_id"}

	viper := viperLib.New()

	//Look for config file location defined as env var
	viper.BindEnv(envVariable)
	configFileLocation = viper.GetString(envVariable)
	if configFileLocation == "" {
		// Get config file from default location
		return config, errors.New(errors.New("Environment variable ALARM_MANAGER_CONFIG_FILE_LOCATION is not defined.").Error())
	}
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configFileLocation)

	if err := viper.ReadInConfig(); err != nil {
		return config, errors.New(errors.New("Fatal error reading config file: ").Error() + err.Error())
	}

	for _, requiredVariable := range requiredVariables {
		if !viper.IsSet(requiredVariable) {
			return config, errors.New("Fatal error config: no " + requiredVariable + " field was found.")
		}
	}

	devices := make(map[string]TuyaDeviceConfig)
	deviceIDs := make(map[string]bool)
	deviceNames := make(map[string]bool)
	readedDevices := viper.GetStringMap("tuya_devices")

	for deviceKey, deviceInfo := range readedDevices {
		deviceInfoValue := reflect.ValueOf(deviceInfo)
		if deviceInfoValue.Kind() != reflect.Map {
			return config, errors.New("Fatal error config: device " + deviceKey + " not a map.")
		} else {

			deviceInfoValueMap := deviceInfoValue.Interface().(map[string]interface{})
			var device TuyaDeviceConfig

			keys := make(map[string]bool)
			for key_name := range deviceInfoValueMap {
				keys[key_name] = true
			}

			for _, requiredDeviceKey := range tuyaDevicesRequiredVariables {
				if _, ok := keys[requiredDeviceKey]; !ok {
					return config, errors.New("Fatal error config: device " + deviceKey + " has no " + requiredDeviceKey + ".")
				} else {
					value := reflect.ValueOf(deviceInfoValueMap[requiredDeviceKey]).Interface().(string)
					switch requiredDeviceKey {
					case "name":
						if _, ok := deviceNames[value]; ok {
							return config, errors.New("Fatal error config: device name '" + value + "' is repeated.")
						}
						device.Name = value
					case "type":
						device.DeviceType = value
					case "host":
						device.Host = value
					case "client_id":
						device.ClientID = value
					case "secret":
						device.Secret = value
					case "device_id":
						if _, ok := deviceIDs[value]; ok {
							return config, errors.New("Fatal error config: device ID " + value + " is repeated.")
						}
						device.DeviceID = value
					}

				}
			}

			deviceNames[device.Name] = true
			deviceIDs[device.DeviceID] = true
			devices[device.Name] = device
		}
	}
	config.Devices = devices
	return config, nil
}
