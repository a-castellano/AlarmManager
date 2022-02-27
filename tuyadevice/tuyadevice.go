package tuyadevice

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type TokenResponse struct {
	Result struct {
		AccessToken  string `json:"access_token"`
		ExpireTime   int    `json:"expire_time"`
		RefreshToken string `json:"refresh_token"`
		UID          string `json:"uid"`
	} `json:"result"`
	Success bool  `json:"success"`
	T       int64 `json:"t"`
}

type TokenRetriever interface {
	RetrieveToken() error
}

type TuyaDevice struct {
	Name         string
	DeviceType   string
	Host         string
	ClientID     string
	Secret       string
	DeviceID     string
	Token        string
	ExpireTime   int
	RefreshToken string
}

func (device *TuyaDevice) RetrieveToken(client http.Client) error {
	body := []byte(``)
	req, _ := http.NewRequest("GET", device.Host+"/v1.0/token?grant_type=1", bytes.NewReader(body))

	device.buildHeader(req, body)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	ret := TokenResponse{}
	unmarshalErr := json.Unmarshal(bs, &ret)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	log.Println("token GET response:", string(bs))
	device.Token = ret.Result.AccessToken
	device.ExpireTime = ret.Result.ExpireTime
	device.RefreshToken = ret.Result.RefreshToken

	return nil

}

func (device TuyaDevice) GetDeviceInfo(client http.Client) ([]byte, error) {
	method := "GET"
	body := []byte(``)
	req, _ := http.NewRequest(method, device.Host+"/v1.0/devices/"+device.DeviceID, bytes.NewReader(body))

	device.buildHeader(req, body)
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
