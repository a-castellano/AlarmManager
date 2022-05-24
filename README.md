# AlarmManager

[Actual Repo](https://git.windmaker.net/a-castellano/AlarmManager)
 [![pipeline status](https://git.windmaker.net/a-castellano/AlarmManager/badges/master/pipeline.svg)](https://git.windmaker.net/a-castellano/AlarmManager/-/commits/master) [![coverage report](https://git.windmaker.net/a-castellano/AlarmManager/badges/master/coverage.svg)](https://git.windmaker.net/a-castellano/AlarmManager/-/commits/master) [![Quality Gate Status](https://sonarqube.windmaker.net/api/project_badges/measure?project=AlarmManager&metric=alert_status)](https://sonarqube.windmaker.net/dashboard?id=AlarmManager)

Basic web API service for managing my Tuya based WiFi alarm.

## Install

Add Widmaker repo and install **windmaker-alarmmanager**:
```bash
wget -O - https://packages.windmaker.net/WINDMAKER-GPG-KEY.pub | sudo apt-key add -
sudo add-apt-repository "deb http://packages.windmaker.net/ focal main"
sudo apt-get update
sudo apt-get install windmaker-alarmmanager
```

## Configuration

This service uses a config file which folder location is defined by environment variable **ALARM_MANAGER_CONFIG_FILE_LOCATION**, inside this folder it must exists a file called **config.toml**.

```toml
[web_server]
port = 3000

[tuya_devices]
[tuya_devices.home_alarm]
name = "Home Alarm"
type = "99AST"
host = "https://openapi.tuyaeu.com"
client_id = "clientID"
secret = "secret"
device_id = "device_id"
```

Client and Device ID's are extracted from [Tuya Developer Account](https://developer.tuya.com).


## Basic usage

### Checking service aliveness
```bash
curl -s -X GET  "http://IP:PORT" | jq
{
  "success": true,
  "msg": "Service up"
}
```

### Show version
```bash
curl -s -X GET  "http://IP:PORT/version" | jq
{
  "success": true,
  "version": "0.1"
}
```


### Show devices
```bash
curl -s -X GET  "http://IP:PORT/devices" | jq
{
  "success": true,
  "data": {
    "deviceid": "Home Alarm"
  }
}
```

### Show device status
```bash
curl -s -X GET  "http://IP:PORT/devices/status/deviceid" | jq
{
  "success": true,
  "msg": "",
  "mode": "disarmed",
  "firing": false,
  "online": true
}
```

### Change device status
```bash
curl -s -X PUT  "http://IP:PORT/devices/status/deviceid" -H 'Coontent-type: application/json' -d '{"mode": "Disarmed"}' | jq
{
  "success": true,
  "msg": "",
  "mode": "disarmed",
  "firing": false,
  "online": true
}
```

