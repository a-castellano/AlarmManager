package tuyadevice

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (device TuyaDevice) buildHeader(req *http.Request, body []byte) {
	req.Header.Set("client_id", device.ClientID)
	req.Header.Set("sign_method", "HMAC-SHA256")

	timeStamp := fmt.Sprint(time.Now().UnixNano() / 1e6)
	req.Header.Set("t", timeStamp)

	if device.Token != "" {
		req.Header.Set("access_token", device.Token)
	}

	sign := device.buildSign(req, body, timeStamp)
	req.Header.Set("sign", sign)
}

func (device TuyaDevice) getHeaderStr(req *http.Request) string {
	signHeaderKeys := req.Header.Get("Signature-Headers")
	if signHeaderKeys == "" {
		return ""
	}
	keys := strings.Split(signHeaderKeys, ":")
	headers := ""
	for _, key := range keys {
		headers += key + ":" + req.Header.Get(key) + "\n"
	}
	return headers
}
