package devices

import (
	"net/http"
	"strings"

	config "github.com/a-castellano/AlarmManager/config_reader"
)

func buildSign(req *http.Request, body []byte, t string, device config.TuyaDevice, token string) string {
	headers := getHeaderStr(req)
	urlStr := getUrlStr(req)
	contentSha256 := Sha256(body)
	stringToSign := req.Method + "\n" + contentSha256 + "\n" + headers + "\n" + urlStr
	signStr := device.ClientID + token + t + stringToSign
	sign := strings.ToUpper(HmacSha256(signStr, device.Secret))
	return sign
}
