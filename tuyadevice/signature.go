package tuyadevice

import (
	"net/http"
	"strings"
)

func (device TuyaDevice) buildSign(req *http.Request, body []byte, timeStamp string) string {
	headers := device.getHeaderStr(req)
	urlStr := getUrlStr(req)
	contentSha256 := Sha256(body)
	stringToSign := req.Method + "\n" + contentSha256 + "\n" + headers + "\n" + urlStr
	signStr := device.ClientID + device.Token + timeStamp + stringToSign
	sign := strings.ToUpper(HmacSha256(signStr, device.Secret))
	return sign
}
