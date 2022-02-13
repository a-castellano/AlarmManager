package devices

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	config "github.com/a-castellano/AlarmManager/config_reader"
)

type AlarmMode int

const (
	FullyArmed AlarmMode = iota + 1
	Disarmed             // disarmed
	HomeArmed            // home
	Sos                  // sos
	Unknown
)

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

type DeviceManager struct {
	DevicesInfo map[string]Alarm
}

func GetDeviceManager(client http.Client, devices map[string]config.TuyaDevice) (DeviceManager, error) {
	manager := DeviceManager{}
	devicesMap := make(map[string]Alarm)
	manager.DevicesInfo = devicesMap
	for deviceName, deviceInfo := range devices {
		token, tokenError := GetToken(client, deviceInfo)
		if tokenError != nil {
			return manager, tokenError
		}
		tuyaDeviceInfo, tuyaDeviceInfoErr := GetDevice(client, deviceInfo, token)
		if tuyaDeviceInfoErr != nil {
			return manager, tuyaDeviceInfoErr
		}
		switch deviceInfo.DeviceType {
		case "99AST":
			alarmInfo := Alarm99AST{}
			if unmarshalErr := json.Unmarshal(tuyaDeviceInfo, &alarmInfo); unmarshalErr != nil {
				panic(unmarshalErr)
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
				//		fmt.Println(statusTuple.Code)
				//		fmt.Println(statusTuple.Value)
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
			manager.DevicesInfo[deviceName] = alarmInfo
		default:
			errorString := fmt.Sprintf("Alarm %s type %s not supported", deviceInfo.Name, deviceInfo.DeviceType)
			return manager, errors.New(errorString)
		}
		//fmt.Println(string(tuyaDeviceInfo))
	}
	return manager, nil
}
