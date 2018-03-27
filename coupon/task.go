package coupon

import (
	"encoding/json"
	"math"
	"strconv"
	"time"

	"code.feelbus.cn/Frameworks/feiniubus-sdk-go/feiniubus"
	"code.feelbus.cn/Frameworks/feiniubus-sdk-go/service/visitor"
	"code.feelbus.cn/Frameworks/puid-go/puid"
	"code.feelbus.cn/Polaris/polaris-sdk-go/fns"
	"code.feelbus.cn/UTP/FeiniuPay.Go/common"
	logstash "github.com/FeiniuBus/log"
)

var _globalUserIdPhones map[string]string
var _notifyClient *fns.FNS
var _visitorClient *visitor.Visitor
var _logger *logstash.Logger
var _context *common.DbContext

const (
	UnUse = iota + 1
	Used
	Expiry
	NotStarted
)

func DoTask(options common.PaymentOptions, notifyClient *fns.FNS, logger *logstash.Logger) {
	_globalUserIdPhones = make(map[string]string)
	sess, _ := feiniubus.NewSession(&feiniubus.Config{
		Endpoint: feiniubus.String(options.FeiniuBusSdk.ServiceUrl),
		Logger:   _logger,
	})
	_visitorClient = visitor.New(sess)
	_notifyClient = notifyClient
	_logger = logger
	_context = common.NewDbContext(options.MySql)
	couponChan := make(chan []CouponDetail, 1000)
	go receiveMessage(options.RabbitMQ, couponChan)
	go handleMessage(couponChan)
	select {}
}

func receiveMessage(o common.RabbitOptions, couponChan chan<- []CouponDetail) {
	queueName := "FeiniuBusCoupon-Queue-" + common.HostingEnvName
	context := common.NewAmqpContext(o, queueName, _logger)
	go context.Connect()
	for {
		msg := <-context.Data
		var task []CouponDetail
		err := json.Unmarshal(msg, &task)
		if err != nil {
			_logger.Error(err.Error())
			continue
		}
		couponChan <- task
	}
}

func handleMessage(couponChan <-chan []CouponDetail) {
	for couponTasks := range couponChan {
		l := len(couponTasks)
		if couponTasks == nil || l <= 0 {
			continue
		}
		userStack := make(map[string]struct{})
		sms := couponTasks[0].Sms
		doTasks(userStack, couponTasks)
		sendSms(userStack, sms)
	}
}

func sendSms(userStack map[string]struct{}, sms string) {
	if len(sms) > 0 {
		phones := getUserPhones(userStack)
		size := 100
		part := int(math.Ceil(float64(len(phones)) / float64(size)))
		var now = time.Now()
		for i := 0; i < part; i++ {
			input := &fns.SendSMSRequest{
				Message: sms,
				Start:   fns.JsonTime(now),
			}
			if i < part-1 {
				input.Numbers = phones[i : (i+1)*size+1]
			} else {
				input.Numbers = phones[i:]
			}
			_notifyClient.SendSMS(input)
		}
	}
}

func doTasks(userStack map[string]struct{}, tasks []CouponDetail) {
	for _, task := range tasks {
		for _, userId := range task.UserIds {
			if userId == "@all@" {
				if len(userStack) > 0 {
					for k := range userStack {
						delete(userStack, k)
					}
				}
				userStack["@all@"] = struct{}{}
				taskWithNoUserId := task
				taskWithNoUserId.UserIds = nil
				setUserInfo(&taskWithNoUserId)
				doEveryTask(taskWithNoUserId)
			} else {
				if _, ok := userStack["@all@"]; !ok {
					userStack[userId] = struct{}{}
				}
				taskWithUserId := task
				taskWithUserId.UserIds = []string{userId}
				doEveryTask(taskWithUserId)
			}
		}
	}
}

func doEveryTask(task CouponDetail) {
	startAt, err := time.ParseInLocation("2006-01-02 15:04:05", task.StartAt, time.Local)
	if err != nil {
		_logger.Error(err.Error())
		return
	}
	expireAt, err := time.ParseInLocation("2006-01-02 15:04:05", task.ExpireAt, time.Local)
	if err != nil {
		_logger.Error(err.Error())
		return
	}
	content, err := json.Marshal(task.Content)
	if err != nil {
		_logger.Error(err.Error())
		return
	}
	scopeAdcodes, err := json.Marshal(task.ScopeAdcodes)
	if err != nil {
		_logger.Error(err.Error())
		return
	}
	var bType, kind int
	switch task.BType {
	case "Carpool":
		bType = 1
		break
	case "Charter":
		bType = 2
		break
	case "Commute":
		bType = 3
		break
	}
	switch task.Kind {
	case "Cash":
		kind = 1
		break
	case "Discount":
		kind = 2
		break
	case "Fixed":
		kind = 3
		break
	}
	for _, userId := range task.UserIds {
		iUserId, _ := strconv.ParseInt(userId, 10, 64)
		for i := 0; i < task.Count; i++ {
			model := CouponRecord{
				BType:        bType,
				Cause:        task.Cause,
				Content:      content,
				StartAt:      startAt,
				CreateAt:     time.Now(),
				ExpireAt:     expireAt,
				Kind:         kind,
				Id:           puid.NewPUID(),
				Name:         task.Name,
				ScopeAdcodes: scopeAdcodes,
				State:        UnUse,
				Title:        task.Title,
				UserId:       iUserId,
				Value:        task.Value,
			}
			_context.Create(&model)
			if _context.Error != nil {
				_logger.Error(_context.Error.Error())
			}
		}
	}
}

func setUserInfo(coupon *CouponDetail) {
	pageIndex := 1
	pageSize := 1000
	count := 0
	input := &visitor.GetPassengerSimpleInput{
		Skip: 0,
		Take: pageSize,
	}
	for {
		if count > 0 && pageIndex > count/pageSize+1 {
			break
		}
		resp, err := _visitorClient.GetPassengerSimple(input)
		if err != nil {
			_logger.Error(err.Error())
			continue
		}
		count = resp.Total
		len := len(resp.Rows)
		if len > 0 {
			for _, row := range resp.Rows {
				coupon.UserIds = append(coupon.UserIds, row.ID)
				_globalUserIdPhones[row.ID] = row.Phone
			}
		}
		input.Skip = input.Skip + input.Take
		pageIndex++
	}
}

func getUserPhones(userStack map[string]struct{}) []string {
	if _, ok := userStack["@all@"]; ok {
		result := make([]string, len(_globalUserIdPhones))
		i := 0
		for _, v := range _globalUserIdPhones {
			result[i] = v
			i++
		}
		return result
	}
	userIds := make([]string, len(userStack))
	i := 0
	for k := range userStack {
		userIds[i] = k
		i++
	}
	input := &visitor.GetPassengerSimpleInput{
		Skip: 0,
		Take: 1000,
		IDs:  userIds,
	}
	resp, err := _visitorClient.GetPassengerSimple(input)
	if err != nil {
		_logger.Error(err.Error())
		return nil
	}
	result := make([]string, len(userStack))
	for i, row := range resp.Rows {
		result[i] = row.Phone
	}
	return result
}
