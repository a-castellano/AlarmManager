package devices

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	config "github.com/a-castellano/AlarmManager/config_reader"
)

func buildHeader(req *http.Request, body []byte, device config.TuyaDevice, token string) {
	req.Header.Set("client_id", device.ClientID)
	req.Header.Set("sign_method", "HMAC-SHA256")

	ts := fmt.Sprint(time.Now().UnixNano() / 1e6)
	req.Header.Set("t", ts)

	if token != "" {
		req.Header.Set("access_token", token)
	}

	sign := buildSign(req, body, ts, device, token)
	req.Header.Set("sign", sign)
}

func getHeaderStr(req *http.Request) string {
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
