package devices

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

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
}

func CreateTuyaDeviceFromConfig(deviceConfig config.TuyaDeviceConfig) tuyadevice.TuyaDevice {

	device := tuyadevice.TuyaDevice{Name: deviceConfig.Name, DeviceType: deviceConfig.DeviceType, Host: deviceConfig.Host, ClientID: deviceConfig.ClientID, Secret: deviceConfig.Secret, DeviceID: deviceConfig.DeviceID}

	return device
}

func (manager *DeviceManager) AddDevice(device tuyadevice.Device) error {
	deviceName := device.GetDeviceName()
	if _, ok := manager.DevicesInfo[deviceName]; ok {
		return fmt.Errorf("Device called '%s' hasalready been added to device manager.", deviceName)
	} else {
		manager.DevicesInfo[deviceName] = device
	}
	return nil
}

func (manager *DeviceManager) Start(client http.Client) error {
	for deviceName, device := range manager.DevicesInfo {
		// Retrieve info foreach device
		tokenError := device.RetrieveToken(client)
		if tokenError != nil {
			return tokenError
		}
		manager.DevicesInfo[deviceName] = device

	}
	return nil
}

func (manager *DeviceManager) RetrieveInfo(client http.Client) error {

	for deviceName, device := range manager.DevicesInfo {
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
			var masterStateSet, masterModeSet bool
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
					}
				}
			}
			manager.AlarmsInfo[deviceName] = alarmInfo
		default:
			errorString := fmt.Sprintf("Alarm %s type %s not supported", deviceName, device.GetDeviceType())
			return errors.New(errorString)
		}
	}
	manager.initiated = true
	return nil
}

func (manager *DeviceManager) ChangeMode(client http.Client, deviceName string, newMode string) error {
	if !manager.initiated {
		errorString := fmt.Sprintf("Device has not retrieved devices info yet.")
		return errors.New(errorString)
	} // Check if device exists
	if alarmDevice, ok := manager.AlarmsInfo[deviceName]; !ok {
		errorString := fmt.Sprintf("Device '%s' is not a managed device.", deviceName)
		return errors.New(errorString)
	} else {
		if equivalentMode, equivalentModeError := alarmDevice.getEquivalentMode(newMode); equivalentModeError != nil {
			return equivalentModeError
		} else {
			changeModeError := manager.DevicesInfo[deviceName].ChangeMode(client, equivalentMode)
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
	//	router.Get("/device/{id}", manager.ShowDeviceInfo)
	return router
}

type DeviceListResponse struct {
	Success bool              `json:"success"`
	Data    map[string]string `json:"data"`
}

func (manager *DeviceManager) ListDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	deviceMap := make(map[string]string)
	for deviceName, device := range manager.DevicesInfo {
		// Retrieve info foreach device
		deviceMap[device.GetDeviceID()] = deviceName
	}
	jsonResponse := DeviceListResponse{Success: true, Data: deviceMap}
	jsonString, _ := json.Marshal(jsonResponse)
	w.Write([]byte(jsonString))
}
