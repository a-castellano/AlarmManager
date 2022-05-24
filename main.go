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
	chi "github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
)

func updateStatus(deviceManager *device_manager.DeviceManager, client http.Client) {
	for range time.Tick(time.Second * 20) {
		log.Println("Updating deviceManager status.")
		deviceManager.RetrieveInfo(client)
	}
}

func main() {

	var version string = "0.2"

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

	log.Println("Initiating Device Manager.")
	deviceManager := device_manager.DeviceManager{DevicesInfo: make(map[string]tuyadevice.Device), AlarmsInfo: make(map[string]device_manager.Alarm)}
	for _, deviceConfig := range config.Devices {
		device := device_manager.CreateTuyaDeviceFromConfig(deviceConfig)
		deviceRef := &device
		addDeviceError := deviceManager.AddDevice(deviceRef)
		if addDeviceError != nil {
			log.Fatal(addDeviceError)
		}
	}
	log.Println("Collecting initial tokens from all devices")
	deviceManager.Start(client)
	log.Println("Obtaining info from all devices")
	deviceManager.RetrieveInfo(client)
	//	fmt.Println(deviceManager.AlarmsInfo)
	//	fmt.Println(deviceManager.AlarmsInfo)
	//changeModeErr := deviceManager.ChangeMode(client, "Home Alarm", "Disarmed")
	//	if changeModeErr != nil {
	//		log.Fatal(changeModeErr)
	//	}

	log.Println("Starting API")
	apiRouter := chi.NewRouter()
	apiRouter.Use(middleware.Logger)
	apiRouter.Use(middleware.Timeout(10 * time.Second))
	apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "msg": "Service up"}`))
	})
	apiRouter.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonResponde := fmt.Sprintf("{\"success\": true, \"version\": \"%s\"}", version)
		w.Write([]byte(jsonResponde))
	})
	apiRouter.Mount("/devices", deviceManager.Routes())

	go updateStatus(&deviceManager, client)
	listenString := fmt.Sprintf(":%d", config.WebPort)
	http.ListenAndServe(listenString, apiRouter)
}
