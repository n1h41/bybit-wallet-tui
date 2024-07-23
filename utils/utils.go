package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

func CreateSignature(apiSecret string, api_key string, recv_window string, params string) string {
	time_stamp := GetTimestamp()
	hmac256 := hmac.New(sha256.New, []byte(apiSecret))
	hmac256.Write([]byte(strconv.FormatInt(time_stamp, 10) + api_key + recv_window + params))
	signature := hex.EncodeToString(hmac256.Sum(nil))
	return signature
}

func GetTimestamp() int64 {
	now := time.Now()
	unixNano := now.UnixNano()
	time_stamp := unixNano / 1000000
	return time_stamp
}

func AddAllHeaders(req *http.Request, api_key string, signature string, time_stamp int64, recv_window string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BAPI-API-KEY", api_key)
	req.Header.Set("X-BAPI-SIGN", signature)
	req.Header.Set("X-BAPI-TIMESTAMP", strconv.FormatInt(time_stamp, 10))
	req.Header.Set("X-BAPI-SIGN-TYPE", "2")
	req.Header.Set("X-BAPI-RECV-WINDOW", recv_window)
}
