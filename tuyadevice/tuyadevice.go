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
	RetrieveToken() (string, error)
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

func (device TuyaDevice) RetrieveToken(client http.Client) (string, error) {
	body := []byte(``)
	req, _ := http.NewRequest("GET", device.Host+"/v1.0/token?grant_type=1", bytes.NewReader(body))

	var token string

	device.buildHeader(req, body)
	resp, err := client.Do(req)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	ret := TokenResponse{}
	unmarshalErr := json.Unmarshal(bs, &ret)
	if unmarshalErr != nil {
		return token, unmarshalErr
	}
	log.Println("token GET response:", string(bs))
	token = ret.Result.AccessToken

	return token, nil

}
