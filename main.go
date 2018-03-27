package main

import (
	"encoding/json"

	"code.feelbus.cn/Polaris/polaris-sdk-go/cc"
	"code.feelbus.cn/Polaris/polaris-sdk-go/credentials"
	"code.feelbus.cn/Polaris/polaris-sdk-go/fns"
	"code.feelbus.cn/Polaris/polaris-sdk-go/polaris"
	"code.feelbus.cn/UTP/FeiniuPay.Go/common"
	"code.feelbus.cn/UTP/FeiniuPay.Go/coupon"
	"code.feelbus.cn/UTP/FeiniuPay.Go/payment"
	logstash "github.com/FeiniuBus/log"
)

var _options common.PaymentOptions
var _logger *logstash.Logger
var _notifyClient *fns.FNS

// consts
const (
	DevPolarisServiceURL = "http://172.16.2.117:5100"
	DevPolarisAccessKey  = "11E7BD5E566B06AE878AFA163EE05ADE"
	DevPolarisSecret     = "1f3eaebb8b33f1b97012caff2b7ee1e8209780d3d4411448814f6b9887647840"
	DevDatabaseName      = "feiniu_ms_payment"

	ProPolarisServiceURL = "http://ip-172-31-23-28.cn-north-1.compute.internal"
	ProPolarisAccessKey  = "11E7E0B738E88DE69D58022EB3F88FF8"
	ProPolarisSecret     = "50b54c45d1540f8b265d705ef9adb8973213d852470c072499264f9466678daa"
	ProDatabaseName      = "feiniubus_ms_payment"
)

func init() {
	cfg := polaris.DefaultConfig()
	attr := make(map[string]string)
	if !common.IsProduction() {
		cfg.Address = polaris.String(DevPolarisServiceURL)
		cfg.Credentials = credentials.NewStaticCredentials(DevPolarisAccessKey, DevPolarisSecret)
		if common.HostingEnvName == "Preview" {
			attr["mysql_db"] = ProDatabaseName
		} else {
			attr["mysql_db"] = DevDatabaseName
		}
	} else {
		cfg.Address = polaris.String(ProPolarisServiceURL)
		cfg.Credentials = credentials.NewStaticCredentials(ProPolarisAccessKey, ProPolarisSecret)
		attr["mysql_db"] = ProDatabaseName
	}
	env := common.HostingEnvName
	attr["mysql_profile"] = "mysql"
	attr["mysql_env"] = env
	attr["rabbitmq_profile"] = "rabbitmq"
	attr["rabbitmq_env"] = env
	request := &cc.FetchRequest{
		EnvironmentName: env,
		Name:            "pay_appsettings",
		Attributes:      attr,
		Language:        cc.GoLang,
	}
	sess := polaris.NewSession()
	response, err := cc.New(sess, cfg).Fetch(request)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(response.Config), &_options)
	if err != nil {
		panic(err)
	}
	_notifyClient = fns.New(sess, cfg)
	l, err := logstash.NewLogstash(false, _options.Logger.Host, _options.Logger.Port)
	if err != nil {
		panic(err)
	}
	_logger = l.With("Application", "Feiniubus-Payment-Background")
}
func main() {
	go coupon.DoTask(_options, _notifyClient, _logger)
	go payment.DoTask(_options, _logger)
	select {}
}
