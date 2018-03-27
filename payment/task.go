package payment

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"code.feelbus.cn/Frameworks/puid-go/puid"
	"code.feelbus.cn/UTP/FeiniuPay.Go/common"
	logstash "github.com/FeiniuBus/log"
	"github.com/pborman/uuid"
)

const (
	Charge = iota + 1
	Pay
)

const (
	UnCharged = iota
	Succeeded
	Refunded
)

var _logger *logstash.Logger
var _context *common.DbContext

func DoTask(options common.PaymentOptions, logger *logstash.Logger) {
	_logger = logger
	_context = common.NewDbContext(options.MySql)
	taskChan := make(chan *Request, 500)
	callbackChan := make(chan *callbackRequest, 1000)
	go receiveMessage(options.RabbitMQ, taskChan)
	go handleMessage(taskChan, callbackChan)
	go handleCallback(callbackChan)
	select {}
}

func receiveMessage(o common.RabbitOptions, taskChan chan<- *Request) {
	queueName := "FeiniuBusPayment-Queue-" + common.HostingEnvName

	context := common.NewAmqpContext(o, queueName, _logger)
	go context.Connect()
	for {
		msg := <-context.Data
		var task Request
		err := json.Unmarshal(msg, &task)
		if err != nil {
			_logger.Error(err.Error())
			continue
		}
		taskChan <- &task
	}
}

func handleMessage(taskChan <-chan *Request, callbackChan chan<- *callbackRequest) {
	for request := range taskChan {
		callback := doEveryTask(request)
		if request.needCallback {
			callbackChan <- callback
		} else {
			if !callback.Succeed {
				_logger.Error(fmt.Sprintf("充值操作失败：%s", callback.Message))
			}
		}
	}
}

func handleCallback(callbackChan chan *callbackRequest) {
	for callback := range callbackChan {
		if !callback.callback() {
			callback.sentCount++
			if callback.sentCount < 10 {
				callbackChan <- callback
			}
		}
	}
}
func doEveryTask(request *Request) *callbackRequest {
	log := BalanceOperateLog{
		Id:          puid.NewPUID(),
		Amount:      request.Amount,
		DeviceId:    request.DeviceId,
		IpAddress:   request.IpAddress,
		OperateAt:   time.Now(),
		OperateType: request.OperatingType,
		Operator:    request.UserId,
	}
	callback := &callbackRequest{
		notifyUrl: request.NotifyUrl,
		OrderId:   request.OutTradeId,
		salt:      uuid.New(),
	}
	if request.OperatingType == Pay {
		err := json.Unmarshal([]byte(request.Metadata), &callback.Metadata)
		if err != nil {
			callback.Message = err.Error()
			_logger.Error(err)
			return callback
		}
	}
	var balance BalanceRecord
	errCount := 0
	tran := _context.Begin()
	_context.Find(&balance, "UserId=?", request.MerchantId)
	if balance.Id <= 0 {
		errCount++
		callback.Message = fmt.Sprintf("指定的支付账户%d不存在", request.MerchantId)
	} else if !balance.Enabled || balance.IsLocked {
		errCount++
		callback.Message = "指定的支付账户不能用或已被锁定"
	} else {
		log.BalanceId = balance.Id
		balance.LastChangeTime = time.Now()
		balance.LastOperator = request.UserId
		targetId := puid.NewPUID()
		switch request.OperatingType {
		case Pay:
			if balance.Amount-balance.FrozenAmount < request.Amount {
				callback.Message = "账户可用余额不足"
				errCount++
			} else {
				log.Remark = "付款"
				request.needCallback = true
				balance.Amount = balance.Amount - request.Amount
				payment := BalancePayment{
					Id:         targetId,
					Amount:     request.Amount,
					BalanceId:  balance.Id,
					CreateAt:   time.Now(),
					OutTradeId: request.OutTradeId,
					Remark:     request.Body,
					Subject:    request.Subject,
					UserId:     request.UserId,
					Account:    strconv.FormatInt(request.MerchantId, 10),
					Channel:    request.Channel,
				}
				if err := tran.Create(&payment).Error; err != nil {
					errCount++
					_logger.Error(err.Error())
				} else {
					flow := BalanceFlow{
						BalanceId:   balance.Id,
						Channel:     request.Channel,
						CreateAt:    time.Now(),
						DeltaAmount: -request.Amount,
						Id:          puid.NewPUID(),
						Operating:   Pay,
						State:       Succeeded,
					}
					if err := tran.Create(&flow).Error; err != nil {
						errCount++
						_logger.Error(err.Error())
					}
				}
			}
			break
		case Charge:
			log.Remark = "充值"
			id, err := strconv.ParseInt(request.OutTradeId, 10, 64)
			if err != nil {
				errCount++
				_logger.Error(err.Error())
				break
			}
			charge := BalanceCharge{Id: id}
			if err := tran.Model(&charge).Update("Succeed", true).Error; err != nil {
				errCount++
				_logger.Error(err.Error())
			} else {
				balance.Amount = balance.Amount + request.Amount
				if err := _context.Select("FlowId").Find(&charge).Error; err == nil {
					flow := BalanceFlow{Id: charge.FlowId}
					if err := tran.Model(&flow).Update("State", Succeeded).Error; err != nil {
						errCount++
						_logger.Error(err.Error())
					}
				} else {
					errCount++
					_logger.Error(err.Error())
				}
			}
			break
		}
		if errCount > 0 {
			callback.Message += "|操作错误"
			tran.Rollback()

		} else {
			if err := tran.Save(&balance).Error; err != nil {
				callback.Message += "|" + err.Error()
				_logger.Error(err.Error())
			} else {
				tran.Commit()
				callback.Succeed = true
			}
			callback.Amount = request.Amount
			callback.ProviderId = targetId
		}
	}
	log.Succeed = callback.Succeed
	_context.Create(&log)
	return callback
}
