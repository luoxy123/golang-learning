package payment

import (
	"time"
)

type Request struct {
	OperatingType int    `json:"operating_type"`
	OutTradeId    string `json:"out_trade_id"`
	MerchantId    int64  `json:"merchant_id"`
	UserId        int64  `json:"user_id"`
	Amount        int    `json:"amount"`
	Subject       string `json:"subject"`
	Body          string `json:"body"`
	DeviceId      string `json:"device_id"`
	IpAddress     string `json:"ip_address"`
	NotifyUrl     string `json:"notify_url"`
	Metadata      string `json:"metadata"`
	Channel       string `json:"channel"`
	needCallback  bool
}

type BalanceRecord struct {
	Id             int64     `gorm:"primary_key;column:Id"`
	UserId         int64     `gorm:"column:UserId"`
	Amount         int       `gorm:"column:Amount"`
	Enabled        bool      `gorm:"column:Enabled"`
	FrozenAmount   int       `gorm:"column:FrozenAmount"`
	LastChangeTime time.Time `gorm:"column:LastChangeTime"`
	LastOperator   int64     `gorm:"column:LastOperator"`
	IsLocked       bool      `gorm:"column:IsLocked"`
}

func (BalanceRecord) TableName() string {
	return "balance_record"
}

type BalanceOperateLog struct {
	Id          int64     `gorm:"primary_key;column:Id"`
	BalanceId   int64     `gorm:"column:BalanceId"`
	Operator    int64     `gorm:"column:Operator"`
	OperateAt   time.Time `gorm:"column:OperateAt"`
	OperateType int       `gorm:"column:OperateType"`
	Amount      int       `gorm:"column:Amount"`
	Succeed     bool      `gorm:"column:Succeed"`
	IpAddress   string    `gorm:"column:IpAddress"`
	DeviceId    string    `gorm:"column:DeviceId"`
	Remark      string    `gorm:"column:Remark"`
}

func (BalanceOperateLog) TableName() string {
	return "balance_operate_log"
}

type BalancePayment struct {
	Id         int64     `gorm:"primary_key;column:Id"`
	Account    string    `gorm:"column:Account"`
	BalanceId  int64     `gorm:"column:BalanceId"`
	OutTradeId string    `gorm:"column:OutTradeId"`
	CreateAt   time.Time `gorm:"column:CreateAt"`
	UserId     int64     `gorm:"column:UserId"`
	Amount     int       `gorm:"column:Amount"`
	Subject    string    `gorm:"column:Subject"`
	Remark     string    `gorm:"column:Remark"`
	Channel    string    `gorm:"column:Channel"`
}

func (BalancePayment) TableName() string {
	return "balance_payment"
}

type BalanceCharge struct {
	Id         int64     `gorm:"primary_key;column:Id"`
	FlowId     int64     `gorm:"column:FlowId"`
	BalanceId  int64     `gorm:"column:BalanceId"`
	MerchantId int64     `gorm:"column:MerchantId"`
	UserId     int64     `gorm:"column:UserId"`
	Channel    string    `gorm:"column:Channel"`
	Amount     int       `gorm:"column:Amount"`
	CreateAt   time.Time `gorm:"column:CreateAt"`
	Remark     string    `gorm:"column:Remark"`
	Succeed    bool      `gorm:"column:Succeed"`
}

func (BalanceCharge) TableName() string {
	return "balance_charge"
}

type BalanceFlow struct {
	Id          int64     `gorm:"primary_key;column:Id"`
	BalanceId   int64     `gorm:"column:BalanceId"`
	CreateAt    time.Time `gorm:"column:CreateAt"`
	DeltaAmount int       `gorm:"column:DeltaAmount"`
	Operating   int       `gorm:"column:Operating"`
	Channel     string    `gorm:"column:Channel"`
	State       int       `gorm:"column:State"`
}

func (BalanceFlow) TableName() string {
	return "balance_flow"
}
