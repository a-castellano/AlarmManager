package devices

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	config "github.com/a-castellano/AlarmManager/config_reader"
)

func GetDevice(client http.Client, device config.TuyaDeviceConfig, token string) ([]byte, error) {
	method := "GET"
	body := []byte(``)
	req, _ := http.NewRequest(method, device.Host+"/v1.0/devices/"+device.DeviceID, bytes.NewReader(body))

	buildHeader(req, body, device, token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return []byte(``), err
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	log.Println("resp:", string(bs))
	return bs, nil
}
