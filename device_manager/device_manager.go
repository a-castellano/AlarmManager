package devices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	config "github.com/a-castellano/AlarmManager/config_reader"
	tuyadevice "github.com/a-castellano/AlarmManager/tuyadevice"
	"github.com/go-chi/chi"
)

type AlarmMode int

const (
	FullyArmed AlarmMode = iota + 1
	Disarmed             // disarmed
	HomeArmed            // home
	Sos                  // sos
	Unknown
)

var AlarmModeMap = map[string]AlarmMode{
	"Armed":     FullyArmed,
	"Disarmed":  Disarmed,
	"HomeArmed": HomeArmed,
	"SOS":       Sos,
}

var AlarmModeAlarmValues = map[AlarmMode]string{
	FullyArmed: "arm",
	Disarmed:   "disarmed",
	HomeArmed:  "home",
	Sos:        "sos",
}

type AlarmInfo struct {
	IP        string
	LocalKey  string
	Latitude  float32
	Longitude float32
	Name      string
	Mode      AlarmMode
	Online    bool
	Firing    bool
}

type Alarm interface {
	ShowInfo() AlarmInfo
	getEquivalentMode(newMode string) (string, error)
}

type Alarm99ASTResult struct {
	ActiveTime  int     `json:"active_time"`
	BizTime     int     `json:"biz_type"`
	Category    string  `json:"category"`
	CreateTime  int     `json:"create_time"`
	Icon        string  `json:"icon"`
	ID          string  `json:"id"`
	IP          string  `json:"ip"`
	Latitude    float32 `json:"lat,string"`
	Longitude   float32 `json:"lon,string"`
	LocalKey    string  `json:"local_key"`
	Model       string  `json:"model"`
	Name        string  `json:"name"`
	Online      bool    `json:"online"`
	OwnerID     int     `json:"owner_id,string"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Sub         bool    `json:"sub"`
	TimeZone    string  `json:"time_zone"`
	UID         string  `json:"uid"`
	UpdateTime  int     `json:"update_time"`
	UUID        string  `json:"uuid"`
	Status      []struct {
		Code  string      `json:"code"`
		Value interface{} `json:"value"`
	} `json:"status"`
}

type Alarm99AST struct {
	Result    Alarm99ASTResult `json:"result"`
	Success   bool             `json:"success"`
	Time      int              `json:"t"`
	AlarmInfo AlarmInfo
}

func (a Alarm99AST) ShowInfo() AlarmInfo {
	return a.AlarmInfo
}

func (a Alarm99AST) getEquivalentMode(newMode string) (string, error) {
	var equivalentMode string
	// Check if newMode is defined
	if newModeValue, ok := AlarmModeMap[newMode]; ok {
		equivalentMode = AlarmModeAlarmValues[newModeValue]
	} else {
		errorString := fmt.Sprintf("Alarm mode '%s' is not defined.", newMode)
		return equivalentMode, errors.New(errorString)
	}
	return equivalentMode, nil
}

type DeviceManager struct {
	initiated   bool
	DevicesInfo map[string]tuyadevice.Device
	AlarmsInfo  map[string]Alarm
	mutex       sync.Mutex
}

func CreateTuyaDeviceFromConfig(deviceConfig config.TuyaDeviceConfig) tuyadevice.TuyaDevice {

	device := tuyadevice.TuyaDevice{Name: deviceConfig.Name, DeviceType: deviceConfig.DeviceType, Host: deviceConfig.Host, ClientID: deviceConfig.ClientID, Secret: deviceConfig.Secret, DeviceID: deviceConfig.DeviceID}

	return device
}

func (manager *DeviceManager) AddDevice(device tuyadevice.Device) error {
	deviceName := device.GetDeviceName()
	deviceID := device.GetDeviceID()
	if _, ok := manager.DevicesInfo[deviceID]; ok {
		return fmt.Errorf("Device called '%s' hasalready been added to device manager.", deviceName)
	} else {
		manager.DevicesInfo[deviceID] = device
	}
	return nil
}

func (manager *DeviceManager) Start(client http.Client) error {
	for deviceID, device := range manager.DevicesInfo {
		// Retrieve info foreach device
		tokenError := device.RetrieveToken(client)
		if tokenError != nil {
			return tokenError
		}
		manager.DevicesInfo[deviceID] = device

	}
	return nil
}

func (manager *DeviceManager) RetrieveInfo(client http.Client) error {

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for deviceID, device := range manager.DevicesInfo {
		deviceName := device.GetDeviceName()
		tokenError := device.RetrieveToken(client)
		if tokenError != nil {
			return tokenError
		}
		log.Println("Retrieving info from device ", deviceName)
		deviceInfo, deviceInfoErr := device.GetDeviceInfo(client)
		log.Println(string(deviceInfo))
		if deviceInfoErr != nil {
			log.Println("Fatal error retrieving info from device ", deviceName)
			return deviceInfoErr
		}
		switch device.GetDeviceType() {
		case "99AST":
			alarmInfo := Alarm99AST{}
			if unmarshalErr := json.Unmarshal(deviceInfo, &alarmInfo); unmarshalErr != nil {
				return unmarshalErr
			}
			// Retrieve Alarm Info
			alarmInfo.AlarmInfo.IP = alarmInfo.Result.IP
			alarmInfo.AlarmInfo.LocalKey = alarmInfo.Result.LocalKey
			alarmInfo.AlarmInfo.Latitude = alarmInfo.Result.Latitude
			alarmInfo.AlarmInfo.Longitude = alarmInfo.Result.Longitude
			alarmInfo.AlarmInfo.Online = alarmInfo.Result.Online
			// Check master mode value

			var masterStateSet, masterModeSet, alarmMessageSet bool
			for _, statusTuple := range alarmInfo.Result.Status {
				if masterStateSet && masterModeSet {
					break
				} else {
					switch statusTuple.Code {
					case "master_mode":
						masterModeValue := fmt.Sprintf("%v", statusTuple.Value)
						switch masterModeValue {
						case "home":
							alarmInfo.AlarmInfo.Mode = HomeArmed
						case "disarmed":
							alarmInfo.AlarmInfo.Mode = Disarmed
						case "arm":
							alarmInfo.AlarmInfo.Mode = FullyArmed
						case "sos":
							alarmInfo.AlarmInfo.Mode = Sos
						default:
							alarmInfo.AlarmInfo.Mode = Unknown
						}
						masterModeSet = true
					case "master_state":
						masterStateValue := fmt.Sprintf("%v", statusTuple.Value)
						alarmInfo.AlarmInfo.Firing = masterStateValue == "alarm"
						masterStateSet = true
					case "alarm_msg":
						alarmMessageValue := fmt.Sprintf("%v", statusTuple.Value)
						alarmMessageSet = alarmMessageValue != "AEEAUABQACAAQQByAG0AYQBkAG8"
					}

				}
			}
			if alarmInfo.AlarmInfo.Firing == false && alarmMessageSet == true {
				alarmInfo.AlarmInfo.Firing = true
			}
			manager.AlarmsInfo[deviceID] = alarmInfo
		default:
			errorString := fmt.Sprintf("Alarm %s type %s not supported", deviceName, device.GetDeviceType())
			return errors.New(errorString)
		}
	}
	manager.initiated = true
	return nil
}

func (manager *DeviceManager) ChangeMode(client http.Client, deviceID string, newMode string) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	if !manager.initiated {
		errorString := fmt.Sprintf("Device has not retrieved devices info yet.")
		return errors.New(errorString)
	} // Check if device exists
	if alarmDevice, ok := manager.AlarmsInfo[deviceID]; !ok {
		errorString := fmt.Sprintf("Device id '%s' is not a managed device.", deviceID)
		return errors.New(errorString)
	} else {
		if equivalentMode, equivalentModeError := alarmDevice.getEquivalentMode(newMode); equivalentModeError != nil {
			return equivalentModeError
		} else {
			changeModeError := manager.DevicesInfo[deviceID].ChangeMode(client, equivalentMode)
			if changeModeError != nil {
				return changeModeError
			}
		}

	}
	return nil
}

func (manager *DeviceManager) Routes() chi.Router {
	router := chi.NewRouter()
	router.Get("/devices", manager.ListDevices)
	router.Route("/devices/status/{id}", func(r chi.Router) {
		r.Use(DeviceCtx)
		r.Get("/", manager.ShowDeviceInfo)
		r.Put("/", manager.UpdateStatus)
	})
	return router
}

type DeviceListResponse struct {
	Success bool              `json:"success"`
	Data    map[string]string `json:"data"`
}

func (manager *DeviceManager) ListDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	deviceMap := make(map[string]string)
	for deviceID, device := range manager.DevicesInfo {
		// Retrieve info foreach device
		deviceMap[deviceID] = device.GetDeviceName()
	}
	jsonResponse := DeviceListResponse{Success: true, Data: deviceMap}
	jsonString, _ := json.Marshal(jsonResponse)
	w.Write([]byte(jsonString))
}

func DeviceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type DeviceStatusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
	Mode    string `json:"mode"`
	Firing  bool   `json:"firing"`
	Online  bool   `json:"online"`
}

func (manager *DeviceManager) ShowDeviceInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	deviceID := r.Context().Value("id").(string)
	var response DeviceStatusResponse
	if _, ok := manager.DevicesInfo[deviceID]; !ok {
		response.Success = false
		response.Message = fmt.Sprintf("Device id '%s' does not exist.", deviceID)
		w.WriteHeader(404)
	} else {
		response.Success = true
		response.Firing = manager.AlarmsInfo[deviceID].ShowInfo().Firing
		response.Online = manager.AlarmsInfo[deviceID].ShowInfo().Online
		response.Mode = AlarmModeAlarmValues[manager.AlarmsInfo[deviceID].ShowInfo().Mode]
	}
	jsonString, _ := json.Marshal(response)
	w.Write([]byte(jsonString))
}

type DeviceChangeStatus struct {
	Mode string `json:"mode"`
}

func (manager *DeviceManager) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var response DeviceStatusResponse
	w.Header().Set("Content-Type", "application/json")
	deviceID := r.Context().Value("id").(string)
	decoder := json.NewDecoder(r.Body)
	var deviceChangeMode DeviceChangeStatus
	err := decoder.Decode(&deviceChangeMode)
	if err != nil {
		response.Success = false
		response.Message = "Failed to decode Response"
		w.WriteHeader(400)
	} else {
		var client http.Client
		currentDeviceSratus := AlarmModeAlarmValues[AlarmModeMap[deviceChangeMode.Mode]]
		if currentDeviceSratus == AlarmModeAlarmValues[manager.AlarmsInfo[deviceID].ShowInfo().Mode] {
			response.Success = true
			response.Message = "Device status has not changed."
			w.WriteHeader(400)
		} else {
			changeModeErr := manager.ChangeMode(client, deviceID, deviceChangeMode.Mode)
			if changeModeErr != nil {
				response.Message = changeModeErr.Error()
				w.WriteHeader(400)
			} else {
				time.Sleep(1 * time.Second)
				retrieveInfoError := manager.RetrieveInfo(client)
				if retrieveInfoError != nil {
					response.Success = false
					response.Message = retrieveInfoError.Error()
					w.WriteHeader(400)
				} else {
					response.Success = true
					response.Firing = manager.AlarmsInfo[deviceID].ShowInfo().Firing
					response.Online = manager.AlarmsInfo[deviceID].ShowInfo().Online
					response.Mode = AlarmModeAlarmValues[manager.AlarmsInfo[deviceID].ShowInfo().Mode]
				}
			}
		}
	}

	jsonString, _ := json.Marshal(response)
	w.Write([]byte(jsonString))
}
