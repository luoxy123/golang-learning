package payment

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type callbackRequest struct {
	ProviderId int64                  `json:"provider_id,string"`
	OrderId    string                 `json:"order_id"`
	Metadata   map[string]interface{} `json:"metadata"`
	Amount     int                    `json:"amount,string"`
	Succeed    bool                   `json:"result"`
	Message    string                 `json:"message"`
	Sign       string                 `json:"signature"`
	notifyUrl  string
	salt       string
	sentCount  int
}

func (request *callbackRequest) callback() bool {
	if len(request.notifyUrl) <= 0 {
		_logger.Warn("飞牛支付没有回调url，忽略本次回调")
		return true
	}
	signUrl, err := request.toUrl()
	if err != nil {
		_logger.Error(err.Error())
		return false
	}
	sign := md5.Sum([]byte(signUrl))
	request.Sign = base64.StdEncoding.EncodeToString(sign[:])
	b, err := json.Marshal(request)
	if err != nil {
		_logger.Error(err.Error())
		return false
	}
	_logger.Info("飞牛支付回调参数:", string(b))
	body := bytes.NewReader(b)
	resp, err := http.Post(request.notifyUrl, "application/json", body)
	if err != nil {
		_logger.Error(err.Error())
		return false
	}
	if resp.StatusCode != 200 {
		e, _ := ioutil.ReadAll(resp.Body)
		_logger.Error(string(e), resp.Status, "请求地址：", request.notifyUrl)
		return false
	}
	return true
}

func (request callbackRequest) toUrl() ([]byte, error) {
	buff, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	var mp map[string]interface{}
	err = json.Unmarshal(buff, &mp)
	if err != nil {
		return nil, err
	}
	mp["salt"] = request.salt
	var keys []string
	for k, _ := range mp {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b bytes.Buffer
	for _, k := range keys {
		if k != "sign" {
			b.WriteString(fmt.Sprintf("%s=%v&", k, mp[k]))
		}
	}
	return b.Bytes()[:b.Len()-1], nil
}
