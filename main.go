package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"time"

	config_reader "github.com/a-castellano/AlarmManager/config_reader"
	device_manager "github.com/a-castellano/AlarmManager/device_manager"
	"github.com/a-castellano/AlarmManager/tuyadevice"
)

func main() {

	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "AlarmManager")
	if e == nil {
		log.SetOutput(logwriter)
		// Remove date prefix
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}

	config, errConfig := config_reader.ReadConfig()
	if errConfig != nil {
		log.Fatal(errConfig)
		return
	}

	deviceManager := device_manager.DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]device_manager.Alarm)}
	for _, deviceConfig := range config.Devices {
		device := device_manager.CreateTuyaDeviceFromConfig(deviceConfig)
		deviceRef := &device
		addDeviceError := deviceManager.AddDevice(deviceRef)
		if addDeviceError != nil {
			log.Fatal(addDeviceError)
		}
	}
	deviceManager.Start(client)
	deviceManager.RetrieveInfo(client)
	fmt.Println(deviceManager.AlarmsInfo)
	deviceManager.RetrieveInfo(client)
	fmt.Println(deviceManager.AlarmsInfo)
}
