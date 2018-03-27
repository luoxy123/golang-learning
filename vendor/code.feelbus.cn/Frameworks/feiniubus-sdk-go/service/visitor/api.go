package visitor

import (
	"code.feelbus.cn/Frameworks/feiniubus-sdk-go/feiniubus"
)

// GetPassengerSimpleInput is
type GetPassengerSimpleInput struct {
	Take int      `json:"take,omitempty"`
	Skip int      `json:"skip,omitempty"`
	IDs  []string `json:"ids,omitempty"`
}

// GetPassengerSimpleOutput is
type GetPassengerSimpleOutput struct {
	Total int               `json:"total,omitempty"`
	Rows  []SimplePassenger `json:"rows,omitempty"`
}

// SimplePassenger is
type SimplePassenger struct {
	ID    string `json:"id,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// GetPassengerSimpleRequest is
func (v *Visitor) GetPassengerSimpleRequest(input *GetPassengerSimpleInput) (req *feiniubus.Request, output *GetPassengerSimpleOutput) {
	op := &feiniubus.Operation{
		Name:           "GetPassengerSimple",
		HTTPMethod:     "POST",
		HTTPPath:       "/account/passenger/simplelist",
		UseQueryString: false,
	}

	if input == nil {
		input = &GetPassengerSimpleInput{}
	}

	output = &GetPassengerSimpleOutput{}
	op.Data = output
	op.Content = input
	op.Unmarshaler = &feiniubus.JSONUnmarshaler{}

	req = v.NewRequest(op)
	return
}

// GetPassengerSimple is
func (v *Visitor) GetPassengerSimple(input *GetPassengerSimpleInput) (*GetPassengerSimpleOutput, error) {
	req, out := v.GetPassengerSimpleRequest(input)
	return out, req.Send()
}

type GetMerchantInfoInput struct {
	ID     string
	Adcode string
}

type GetMerchantInfoOutput struct {
	Items []MerchantInfo
}

type MerchantInfo struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Adcode     string       `json:"adcode"`
	Accounting string       `json:"accounting_method"`
	PayMethod  string       `json:"pay_method"`
	Contact    string       `json:"contact"`
	Phone      string       `json:"phone"`
	PayTimed   int          `json:"pay_timed"`
	Deps       []Department `json:"departments"`
}

type Department struct {
	ID    string `json:"department_id"`
	Name  string `json:"department_name"`
	Users []User `json:"users"`
}

type User struct {
	ID    string `json:"user_id"`
	Name  string `json:"user_name"`
	Phone string `json:"user_phone"`
}

func (v *Visitor) GetMerchantInfoRequest(input *GetMerchantInfoInput) (req *feiniubus.Request, out *GetMerchantInfoOutput) {
	op := &feiniubus.Operation{
		Name:           "GetMerchantInfo",
		HTTPMethod:     "GET",
		HTTPPath:       "/MerchantInformation/sing",
		UseQueryString: true,
		Params:         make(map[string]string),
	}

	if input == nil {
		input = &GetMerchantInfoInput{}
	}

	op.Params["adcode"] = input.Adcode
	if input.ID != "" {
		op.Params["id"] = input.ID
	}
	out = &GetMerchantInfoOutput{}
	op.Data = out
	op.Unmarshaler = &GetMerchantInfoUnmarshal{}
	req = v.NewRequest(op)
	return
}

func (v *Visitor) GetMerchantInfo(input *GetMerchantInfoInput) (*GetMerchantInfoOutput, error) {
	req, out := v.GetMerchantInfoRequest(input)
	return out, req.Send()
}

type GetDriverListInput struct {
	DynamicQuery interface{}
}

type GetDriverListOutput struct {
	Total int       `json:"total"`
	Rows  []*Driver `json:"rows"`
}

type Driver struct {
	ID            string `json:"id"`
	Phone         string `json:"phone"`
	Number        string `json:"number"`
	Name          string `json:"name"`
	IDCard        string `json:"id_card_number"`
	Address       string `json:"address"`
	Avatar        string `json:"avatar"`
	BPhone        string `json:"business_phone"`
	Adcode        string `json:"adcode"`
	Own           bool   `json:"is_own"`
	Gender        string `json:"gender"`
	CreateTime    string `json:"create_time"`
	Disabled      bool   `json:"disabled"`
	LastLoginTime string `json:"last_login_time"`
	CompanyID     string `json:"company_id"`
	Company       string `json:"company"`
}

func (v *Visitor) GetDriverListRequest(input *GetDriverListInput) (req *feiniubus.Request, out *GetDriverListOutput) {
	op := &feiniubus.Operation{
		Name:       "GetDriverList",
		HTTPMethod: "POST",
		HTTPPath:   "/account/driver/list",
	}

	op.Content = input.DynamicQuery
	out = &GetDriverListOutput{}
	op.Data = out
	op.Unmarshaler = &feiniubus.JSONUnmarshaler{}

	req = v.NewRequest(op)
	return
}

func (v *Visitor) GetDriverList(input *GetDriverListInput) (*GetDriverListOutput, error) {
	req, out := v.GetDriverListRequest(input)
	return out, req.Send()
}

type GetPassengerListInput struct {
	DynamicQuery interface{}
}

type GetPassengerListOutput struct {
	Total int          `json:"total"`
	Rows  []*Passenger `json:"rows"`
}

type Passenger struct {
	ID            string `json:"id"`
	Phone         string `json:"phone"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	CreateTime    string `json:"create_time"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	Avatar        string `json:"avatar"`
	RTerminal     string `json:"register_terminal"`
	LastLoginTime string `json:"last_login_time"`
	Disabled      bool   `json:"disabled"`
	Invitor       string `json:"invite_id"`
}

func (v *Visitor) GetPassengerListRequest(input *GetPassengerListInput) (req *feiniubus.Request, out *GetPassengerListOutput) {
	op := &feiniubus.Operation{
		Name:       "GetPassengerList",
		HTTPMethod: "POST",
		HTTPPath:   "/account/passenger/list",
	}

	op.Content = input.DynamicQuery
	out = &GetPassengerListOutput{}
	op.Data = out
	op.Unmarshaler = &feiniubus.JSONUnmarshaler{}

	req = v.NewRequest(op)
	return
}

func (v *Visitor) GetPassengerList(input *GetPassengerListInput) (*GetPassengerListOutput, error) {
	req, out := v.GetPassengerListRequest(input)
	return out, req.Send()
}

type GetMerchantListInput struct {
	DynamicQuery interface{}
}

type GetMerchantListOutput struct {
	Total int         `json:"total"`
	Rows  []*Merchant `json:"rows"`
}

type Merchant struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Contact           string `json:"contact"`
	Phone             string `json:"phone"`
	License           string `json:"business_license"`
	SmsName           string `json:"sms_name"`
	Adcode            string `json:"adcode"`
	Accounting        string `json:"accounting_method"`
	Source            string `json:"order_source"`
	BillingCycle      string `json:"billing_cycle"`
	Salesman          string `json:"salesman"`
	SendSms           bool   `json:"send_sms"`
	SmsTemplate       string `json:"sms_template"`
	Address           string `json:"address"`
	TypeID            string `json:"type_id"`
	GlobalPrice       bool   `json:"global_price"`
	Audit             string `json:"audit"`
	CreateTime        string `json:"create_time"`
	Commissioner      string `json:"commissioner"`
	CommissionerPhone string `json:"commissioner_phone"`
	Disabled          bool   `json:"disabled"`
	PayMethod         string `json:"pay_method"`
	PayTime           int    `json:"pay_timed"`
}

func (v *Visitor) GetMerchantListRequest(input *GetMerchantListInput) (req *feiniubus.Request, out *GetMerchantListOutput) {
	op := &feiniubus.Operation{
		Name:       "GetMerchantList",
		HTTPMethod: "POST",
		HTTPPath:   "/MerchantInformation/list",
	}

	op.Content = input.DynamicQuery
	out = &GetMerchantListOutput{}
	op.Data = out
	op.Unmarshaler = &feiniubus.JSONUnmarshaler{}

	req = v.NewRequest(op)
	return
}

func (v *Visitor) GetMerchantList(input *GetMerchantListInput) (*GetMerchantListOutput, error) {
	req, out := v.GetMerchantListRequest(input)
	return out, req.Send()
}

type GetMerchantUserListInput struct {
	DynamicQuery interface{}
}

type GetMerchantUserListOutput struct {
	Total int             `json:"total"`
	Rows  []*MerchantUser `json:"rows"`
}

type MerchantUser struct {
	ID            string `json:"id"`
	Phone         string `json:"phone"`
	Name          string `json:"name"`
	Gender        string `json:"gender"`
	Avatar        string `json:"avatar"`
	Merchant      string `json:"merchant_id"`
	Department    string `json:"department_id"`
	Type          string `json:"account_type"`
	Disabled      bool   `json:"disabled"`
	LastLoginTime string `json:"last_login_time"`
	CreateTime    string `json:"create_time"`
	Role          string `json:"role"`
	Comment       string `json:"comment"`
	Account       string `json:"account"`
	DName         string `json:"department_name"`
}

func (v *Visitor) GetMerchantUserListRequest(input *GetMerchantUserListInput) (req *feiniubus.Request, out *GetMerchantUserListOutput) {
	op := &feiniubus.Operation{
		Name:       "GetMerchantUserList",
		HTTPMethod: "POST",
		HTTPPath:   "/account/merchant/list",
	}

	op.Content = input.DynamicQuery
	out = &GetMerchantUserListOutput{}
	op.Data = out
	op.Unmarshaler = &feiniubus.JSONUnmarshaler{}

	req = v.NewRequest(op)
	return
}

func (v *Visitor) GetMerchantUserList(input *GetMerchantUserListInput) (*GetMerchantUserListOutput, error) {
	req, out := v.GetMerchantUserListRequest(input)
	return out, req.Send()
}

type GetMerchantDepartmentListInput struct {
	DynamicQuery interface{}
}

type GetMerchantDepartmentListOutput struct {
	Total int                   `json:"total"`
	Rows  []*MerchantDepartment `json:"rows"`
}

type MerchantDepartment struct {
	ID       string `json:"id"`
	Merchant string `json:"merchant_id"`
	Name     string `json:"name"`
}

func (v *Visitor) GetMerchantDepartmentListRequest(input *GetMerchantDepartmentListInput) (req *feiniubus.Request, out *GetMerchantDepartmentListOutput) {
	op := &feiniubus.Operation{
		Name:       "GetMerchantDepartmentList",
		HTTPMethod: "POST",
		HTTPPath:   "/Merchant/Department/List",
	}

	op.Content = input.DynamicQuery
	out = &GetMerchantDepartmentListOutput{}
	op.Data = out
	op.Unmarshaler = &feiniubus.JSONUnmarshaler{}

	req = v.NewRequest(op)
	return
}

func (v *Visitor) GetMerchantDepartmentList(input *GetMerchantDepartmentListInput) (*GetMerchantDepartmentListOutput, error) {
	req, out := v.GetMerchantDepartmentListRequest(input)
	return out, req.Send()
}

type GetAdminInput struct {
	ID string
}

type GetAdminOutput struct {
	Admin Admin
}

type Admin struct {
	ID         string       `json:"id"`
	Phone      string       `json:"phone"`
	NickName   string       `json:"nick_name"`
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	Gender     string       `json:"gender"`
	CreateTime string       `json:"create_time"`
	Avatar     string       `json:"avatatr"`
	Number     string       `json:"number"`
	Adcode     string       `json:"adcode"`
	CityName   string       `json:"city_name"`
	Super      bool         `json:"is_super"`
	Disabled   bool         `json:"disabled"`
	Department string       `json:"department_id"`
	Salesman   bool         `json:"is_salesman"`
	Roles      []*AdminRole `json:"roles"`
}

type AdminRole struct {
	ID   string `json:"id"`
	Name string `json:"string"`
}

func (v *Visitor) GetAdminRequest(input *GetAdminInput) (req *feiniubus.Request, out *GetAdminOutput) {
	op := &feiniubus.Operation{
		Name:           "GetAdmin",
		HTTPMethod:     "GET",
		HTTPPath:       "/account/admin",
		UseQueryString: true,
		Params:         make(map[string]string),
	}

	op.Params["id"] = input.ID
	out = &GetAdminOutput{}
	op.Data = out
	op.Unmarshaler = &GetAdminUnmarshal{}

	req = v.NewRequest(op)
	return
}

func (v *Visitor) GetAdmin(input *GetAdminInput) (*GetAdminOutput, error) {
	req, out := v.GetAdminRequest(input)
	return out, req.Send()
}

type GetMerchantDepartmentInput struct {
	ID string
}

type GetMerchantDepartmentOutput struct {
	Dep MerchantDepartment
}

func (v *Visitor) GetMerchantDepartmentRequest(input *GetMerchantDepartmentInput) (req *feiniubus.Request, out *GetMerchantDepartmentOutput) {
	op := &feiniubus.Operation{
		Name:           "GetMerchantDepartment",
		HTTPMethod:     "GET",
		HTTPPath:       "/Merchant/Department",
		UseQueryString: true,
		Params:         make(map[string]string),
	}

	op.Params["id"] = input.ID
	out = &GetMerchantDepartmentOutput{}
	op.Data = out
	op.Unmarshaler = &GetMerchantDepartmentUnmarshal{}

	req = v.NewRequest(op)
	return
}

func (v *Visitor) GetMerchantDepartment(input *GetMerchantDepartmentInput) (*GetMerchantDepartmentOutput, error) {
	req, out := v.GetMerchantDepartmentRequest(input)
	return out, req.Send()
}

type GetMerchantUserWithNameInput struct {
	DynamicQuery interface{}
}

type GetMerchantUserWithNameOutput struct {
	Total int      `json:"total"`
	Rows  []*MUser `json:"rows"`
}

type MUser struct {
	ID                string `json:"id"`
	Phone             string `json:"phone"`
	Name              string `json:"name"`
	Disabled          bool   `json:"disabled"`
	CreateTime        string `json:"create_time"`
	Role              string `json:"role"`
	Account           string `json:"account"`
	AccountType       string `json:"account_type"`
	MName             string `json:"merchant_name"`
	MID               string `json:"merchant_id"`
	MAccountingMethod string `json:"merchant_accounting_method"`
	MPayMethod        string `json:"merchant_pay_method"`
	MSalesman         string `json:"merchant_salesman"`
	MAcode            string `json:"merchant_adcode"`
	DID               string `json:"department_id"`
	DName             string `json:"department_name"`
}

func (v *Visitor) GetMerchantUserWithNameRequest(input *GetMerchantUserWithNameInput) (req *feiniubus.Request, out *GetMerchantUserWithNameOutput) {
	op := &feiniubus.Operation{
		Name:           "GetMerchantUserWithName",
		HTTPMethod:     "POST",
		HTTPPath:       "/account/merchant_users",
		UseQueryString: true,
	}

	op.Content = input.DynamicQuery
	out = &GetMerchantUserWithNameOutput{}
	op.Data = out
	op.Unmarshaler = &feiniubus.JSONUnmarshaler{}

	req = v.NewRequest(op)
	return
}

func (v *Visitor) GetMerchantUserWithName(input *GetMerchantUserWithNameInput) (*GetMerchantUserWithNameOutput, error) {
	req, out := v.GetMerchantUserWithNameRequest(input)
	return out, req.Send()
}
