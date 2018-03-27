package fns

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"code.feelbus.cn/Polaris/polaris-sdk-go/polaris"
)

type PublishRequest struct {
	Topic   string `json:"topic"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type PublishResponse struct {
	MessageId string `json:"message_id"`
}

func (fns *FNS) Publish(request *PublishRequest) (*PublishResponse, error) {
	op := &polaris.Operation{
		HTTPMethod: http.MethodPost,
		HTTPPath:   "/v1/publish",
	}

	req, err := fns.NewRequest(op, nil, request)
	if err != nil {
		return nil, err
	}

	bs, err := req.Send()
	if err != nil {
		return nil, err
	}

	var resp PublishResponse
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type ConfirmSubscriptionRequest struct {
	MessageID       string
	SubscriptionARN string
	Token           string
}

func (fns *FNS) ConfirmSubscription(request *ConfirmSubscriptionRequest) error {
	op := &polaris.Operation{
		HTTPMethod: http.MethodGet,
		HTTPPath:   "/v1/subscription/validate",
	}

	values := url.Values{}
	values.Set("message_id", request.MessageID)
	values.Set("ss_arn", request.SubscriptionARN)
	values.Set("token", request.Token)

	req, err := fns.NewRequest(op, values, nil)
	if err != nil {
		return err
	}

	_, err = req.Send()
	return err
}

type CancelSMSRequest struct {
	IDs []string
}

func (fns *FNS) CancelSMS(request *CancelSMSRequest) error {
	if len(request.IDs) == 0 {
		return errors.New("必须至少指定一个ID值")
	}

	op := &polaris.Operation{
		HTTPMethod: http.MethodDelete,
		HTTPPath:   "/v1/sms/cancel",
	}

	values := url.Values{}
	values.Set("ids", strings.Join(request.IDs, ","))

	req, err := fns.NewRequest(op, values, nil)
	if err != nil {
		return err
	}

	_, err = req.Send()
	return err
}

type DeletePushRequest struct {
	TaskID string
}

func (fns *FNS) DeletePush(request *DeletePushRequest) error {
	op := &polaris.Operation{
		HTTPMethod: http.MethodDelete,
		HTTPPath:   "/v1/push/delete",
	}

	values := url.Values{}
	values.Set("id", request.TaskID)

	req, err := fns.NewRequest(op, values, nil)
	if err != nil {
		return err
	}

	_, err = req.Send()
	return err
}

type SchedulePushRequest struct {
	Type        string            `json:"type"`
	Range       string            `json:"range"`
	Content     string            `json:"content"`
	Tags        []string          `json:"tags"`
	Identifiers []string          `json:"identifiers"`
	Subject     string            `json:"subject"`
	Terminal    string            `json:"terminal"`
	Title       string            `json:"title"`
	Extra       map[string]string `json:"extra"`
	TimeToLive  int               `json:"time_to_live"`
	Start       JsonTime          `json:"start"`
}

type SchedulePushResponse struct {
	TaskID string `json:"task_id"`
}

func (fns *FNS) SchedulePush(request *SchedulePushRequest) (*SchedulePushResponse, error) {
	op := &polaris.Operation{
		HTTPMethod: http.MethodPost,
		HTTPPath:   "/v1/push/schedule",
	}

	req, err := fns.NewRequest(op, nil, request)
	if err != nil {
		return nil, err
	}

	bs, err := req.Send()
	if err != nil {
		return nil, err
	}

	var resp SchedulePushResponse
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type SendPushRequest struct {
	Type        string            `json:"type"`
	Range       string            `json:"range"`
	Content     string            `json:"content"`
	Tags        []string          `json:"tags"`
	Identifiers []string          `json:"identifiers"`
	Subject     string            `json:"subject"`
	Terminal    string            `json:"terminal"`
	Title       string            `json:"title"`
	Extra       map[string]string `json:"extra"`
	TimeToLive  int               `json:"time_to_live"`
}

type SendPushResponse struct {
	TaskID string `json:"task_id"`
}

func (fns *FNS) SendPush(request *SendPushRequest) (*SendPushResponse, error) {
	op := &polaris.Operation{
		HTTPMethod: http.MethodPost,
		HTTPPath:   "/v1/push/send",
	}

	req, err := fns.NewRequest(op, nil, request)
	if err != nil {
		return nil, err
	}

	bs, err := req.Send()
	if err != nil {
		return nil, err
	}

	var resp SendPushResponse
	err = json.Unmarshal(bs, &resp)
	return &resp, err
}

type SendSMSRequest struct {
	Start   JsonTime `json:"start"`
	Numbers []string `json:"numbers"`
	Message string   `json:"message"`
}

type SendSMSResponse struct {
	TaskID string `json:"task_id"`
}

func (fns *FNS) SendSMS(request *SendSMSRequest) (*SendSMSResponse, error) {
	op := &polaris.Operation{
		HTTPMethod: http.MethodPost,
		HTTPPath:   "/v1/sms/send",
	}

	req, err := fns.NewRequest(op, nil, request)
	if err != nil {
		return nil, err
	}

	bs, err := req.Send()
	if err != nil {
		return nil, err
	}

	var resp SendSMSResponse
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type UpdatePushRequest struct {
	TaskID      string            `json:"-"`
	Type        string            `json:"type"`
	Range       string            `json:"range"`
	Content     string            `json:"content"`
	Tags        []string          `json:"tags"`
	Identifiers []string          `json:"identifiers"`
	Subject     string            `json:"subject"`
	Terminal    string            `json:"terminal"`
	Title       string            `json:"title"`
	Extra       map[string]string `json:"extra"`
	TimeToLive  int               `json:"time_to_live"`
	Start       JsonTime          `json:"start"`
}

func (fns *FNS) UpdatePush(request *UpdatePushRequest) error {
	op := &polaris.Operation{
		HTTPMethod: http.MethodPut,
		HTTPPath:   "/v1/push/update",
	}

	values := url.Values{}
	values.Set("id", request.TaskID)

	req, err := fns.NewRequest(op, values, request)
	if err != nil {
		return err
	}

	_, err = req.Send()
	return err
}

type UpdateSMSRequest struct {
	TaskID  string   `json:"-"`
	Start   JsonTime `json:"start"`
	Numbers []string `json:"numbers"`
	Message string   `json:"message"`
}

func (fns *FNS) UpdateSMS(request *UpdateSMSRequest) error {
	op := &polaris.Operation{
		HTTPMethod: http.MethodPut,
		HTTPPath:   "/v1/sms/update",
	}

	values := url.Values{}
	values.Set("id", request.TaskID)

	req, err := fns.NewRequest(op, values, request)
	if err != nil {
		return err
	}

	_, err = req.Send()
	return err
}

type SmsStatus string

const (
	SmsStatusNotYet  SmsStatus = "NotYet"
	SmsStatusSucceed SmsStatus = "Succeed"
	SmsStatusFailed  SmsStatus = "Failed"
)

type SmsStatusRequest struct {
	Id string `json:"id"`
}

type SmsStatusResponse struct {
	Id        string                 `json:"id"`
	SmsStatus SmsStatus              `json:"status"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func (fns *FNS) SmsStatus(request *SmsStatusRequest) (*SmsStatusResponse, error) {
	if request == nil {
		return nil, errors.New("request不能为空")
	}

	if request.Id == "" {
		return nil, errors.New("Id不能为空")
	}

	op := &polaris.Operation{
		HTTPMethod: http.MethodGet,
		HTTPPath:   "/v1/sms/sms_status",
	}

	values := url.Values{}
	values.Set("id", request.Id)

	req, e := fns.NewRequest(op, values, request)
	if e != nil {
		return nil, e
	}

	resp, e := req.Send()
	if e != nil {
		return nil, e
	}

	var r SmsStatusResponse
	e = json.Unmarshal(resp, &r)
	if e != nil {
		return nil, e
	}

	return &r, nil
}

type SmssStatusRequest struct {
	Ids []string `json:"ids"`
}

type SmssStatusResponseItem SmsStatusResponse

type SmssStatusResponse struct {
	SmssStatus []SmssStatusResponseItem `json:"smss_status"`
}

func (fns *FNS) SmssStatus(request *SmssStatusRequest) (*SmssStatusResponse, error) {
	if request == nil {
		return nil, errors.New("request不能为空")
	}

	if request.Ids == nil {
		return nil, errors.New("Ids不能为空")
	}

	if len(request.Ids) == 0 {
		return nil, errors.New("Ids需要至少包含一个参数")
	}

	op := &polaris.Operation{
		HTTPMethod: http.MethodGet,
		HTTPPath:   "/v1/sms/smss_status",
	}

	values := url.Values{}
	values.Set("ids", strings.Join(request.Ids, ","))

	req, e := fns.NewRequest(op, values, request)
	if e != nil {
		return nil, e
	}

	resp, e := req.Send()
	if e != nil {
		return nil, e
	}

	var r SmssStatusResponse
	e = json.Unmarshal(resp, &r)
	if e != nil {
		return nil, e
	}

	return &r, nil
}

type RegisterAliasRequest struct {
	ID      string `json:"identifier"`
	Subject string `json:"subject"`
	Alias   string `json:"alias"`
}

func (fns *FNS) RegisterAlias(request *RegisterAliasRequest) error {
	op := &polaris.Operation{
		HTTPMethod: http.MethodPost,
		HTTPPath:   "/v1/push/register/alias",
	}

	req, err := fns.NewRequest(op, nil, request)
	if err != nil {
		return err
	}

	_, err = req.Send()
	return err
}
