package tuyadevice

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type TokenResponse struct {
	Result struct {
		AccessToken     string `json:"access_token"`
		TokenExpireTime int    `json:"expire_time"`
		RefreshToken    string `json:"refresh_token"`
		UID             string `json:"uid"`
	} `json:"result"`
	Success bool  `json:"success"`
	T       int64 `json:"t"`
}

type TokenRetriever interface {
	RetrieveToken() error
}

type TuyaDevice struct {
	Name            string
	DeviceType      string
	Host            string
	ClientID        string
	Secret          string
	DeviceID        string
	Token           string
	TokenExpireTime int64
	RefreshToken    string
}

func (device *TuyaDevice) UpdateToken(client http.Client) error {

	now := time.Now()
	currentTimestamp := now.Unix()
	if device.TokenExpireTime-currentTimestamp < 600 {
		log.Println("Device " + device.Name + " token will expire in less than 10 secons, retrive new token.")
		body := []byte(``)
		req, _ := http.NewRequest("GET", device.Host+"/v1.0/token/"+device.RefreshToken, bytes.NewReader(body))

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
		log.Println("refresh token GET response:", string(bs))
		device.Token = ret.Result.AccessToken
		now := time.Now() // current local time
		device.TokenExpireTime = now.Unix() + int64(ret.Result.TokenExpireTime)
		device.RefreshToken = ret.Result.RefreshToken
	}
	return nil
}

func (device *TuyaDevice) RetrieveToken(client http.Client) error {
	if device.TokenExpireTime == 0 { // New token
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
		now := time.Now() // current local time
		device.TokenExpireTime = now.Unix() + int64(ret.Result.TokenExpireTime)
		device.RefreshToken = ret.Result.RefreshToken

	} else {
		// refresh token
		return device.UpdateToken(client)
	}
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
