package devices

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	config "github.com/a-castellano/AlarmManager/config_reader"
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

func GetToken(client http.Client, device config.TuyaDevice) (string, error) {
	method := "GET"
	body := []byte(``)
	req, _ := http.NewRequest(method, device.Host+"/v1.0/token?grant_type=1", bytes.NewReader(body))

	var token string

	buildHeader(req, body, device, token)
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
