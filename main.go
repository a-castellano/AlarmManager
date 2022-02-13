package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"time"

	config_reader "github.com/a-castellano/AlarmManager/config_reader"
	device_manager "github.com/a-castellano/AlarmManager/device_manager"
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
	deviceManager, _ := device_manager.GetDeviceManager(client, config.Devices)
	fmt.Println(deviceManager)
}
