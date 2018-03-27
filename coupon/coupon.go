package coupon

import "time"

type CouponDetail struct {
	Count        int      `json:"count"`
	StartAt      string   `json:"start_at"`
	ExpireAt     string   `json:"expire_at"`
	ScopeAdcodes []string `json:"scope_adcodes"`
	Name         string   `json:"name"`
	Title        string   `json:"title"`
	Kind         string   `json:"kind"`
	Value        int      `json:"value"`
	BType        string   `json:"btype"`
	Content      []string `json:"content"`
	Cause        string   `json:"cause"`
	UserIds      []string `json:"user_ids"`
	Sms          string   `json:"sms"`
}

// CouponRecord is
type CouponRecord struct {
	Id           int64     `gorm:"primary_key;column:Id"`
	CreateAt     time.Time `gorm:"column:CreateAt"`
	StartAt      time.Time `gorm:"column:StartAt"`
	ExpireAt     time.Time `gorm:"column:ExpireAt"`
	ScopeAdcodes []byte    `gorm:"column:ScopeAdcodes"`
	Name         string    `gorm:"column:Name"`
	Title        string    `gorm:column:"Title"`
	Kind         int       `gorm:"column:Kind"`
	Value        int       `gorm:"column:Value"`
	BType        int       `gorm:"column:BType"`
	Content      []byte    `gorm:"column:Content"`
	Cause        string    `gorm:"column:Cause"`
	UserId       int64     `gorm:"column:UserId"`
	State        int       `gorm:"column:state"`
}

func (CouponRecord) TableName() string {
	return "coupon_record"
}
