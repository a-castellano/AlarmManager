package tuyadevice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
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

type ChangeModeResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"msg"`
	Success   bool   `json:"success"`
	Timestamp int    `json:"t"`
	TID       string `json:"tid"`
}

type Device interface {
	GetDeviceInfo(http.Client) ([]byte, error)
	GetDeviceID() string
	RetrieveToken(http.Client) error
	GetDeviceType() string
	GetDeviceName() string
	ChangeMode(http.Client, string) error
}

type TuyaDevice struct {
	Name            string `valid:"required"`
	DeviceType      string `valid:"required"`
	Host            string `valid:"required"`
	ClientID        string `valid:"required"`
	Secret          string `valid:"required"`
	DeviceID        string `valid:"required"`
	Token           string
	TokenExpireTime int64
	RefreshToken    string
}

func (device *TuyaDevice) GetDeviceType() string {
	return device.DeviceType
}

func (device *TuyaDevice) GetDeviceID() string {
	return device.DeviceID
}

func (device *TuyaDevice) GetDeviceName() string {
	return device.Name
}

func (device *TuyaDevice) Validate() error {
	_, err := govalidator.ValidateStruct(device)
	if err != nil {
		return err
	}
	return nil
}

func (device *TuyaDevice) updateToken(client http.Client) error {

	now := time.Now()
	currentTimestamp := now.Unix()
	log.Println(currentTimestamp)
	log.Println(device.TokenExpireTime)
	log.Println(device.TokenExpireTime - currentTimestamp)
	if device.TokenExpireTime-currentTimestamp < 0 {
		log.Println("Device " + device.Name + " token has expired, retrive new token.")
		body := []byte(``)
		req, _ := http.NewRequest("GET", device.Host+"/v1.0/token/"+device.RefreshToken, nil)

		device.buildHeader(req, body)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		log.Println(resp.Body)
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
		log.Println("New token expire time: ", device.TokenExpireTime)
	} else {
		log.Println("Device " + device.Name + " token has not expired yet.")
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
		return device.updateToken(client)
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

func (device TuyaDevice) ChangeMode(client http.Client, mode string) error {
	log.Println("Changing device " + device.GetDeviceName() + " mode to '" + mode + "'.")
	method := "POST"
	commandString := fmt.Sprintf("{\"commands\":[{\"code\":\"master_mode\",\"value\":\"%s\"}]}", mode)
	body := []byte(commandString)
	req, _ := http.NewRequest(method, device.Host+"/v1.0/devices/"+device.DeviceID+"/commands", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	device.buildHeader(req, body)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)

	log.Println("resp:", string(bs))

	// retrieve Response status
	response := ChangeModeResponse{}
	unmarshalErr := json.Unmarshal(bs, &response)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	if !response.Success {
		errorString := fmt.Sprintf("Device '%s' failed to change state to %s, error was '%s'.", device.GetDeviceName(), mode, response.Message)
		return errors.New(errorString)
	}
	log.Println("Changing device " + device.GetDeviceName() + " mode to '" + mode + "' succeded.")
	return nil
}
