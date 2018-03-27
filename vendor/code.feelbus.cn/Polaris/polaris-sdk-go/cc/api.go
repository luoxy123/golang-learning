package cc

import (
	"encoding/json"
	"net/http"

	"code.feelbus.cn/Polaris/polaris-sdk-go/polaris"
)

const (
	CSharpLang = "CSharp"
	GoLang     = "Go"
)

// FetchRequest is
type FetchRequest struct {
	Name            string            `json:"name"`
	EnvironmentName string            `json:"env"`
	Attributes      map[string]string `json:"attrs"`
	Language        string            `json:"lang"`
}

// FetchResponse is
type FetchResponse struct {
	Config          string `json:"config"`
	EnvironmentName string `json:"env"`
	Format          string `json:"format"`
	Index           string `json:"index"`
}

// Fetch is
func (c *CC) Fetch(request *FetchRequest) (*FetchResponse, error) {
	operation := &polaris.Operation{
		HTTPMethod: http.MethodPost,
		HTTPPath:   "/v1/fetch",
	}

	if request.Language == "" {
		request.Language = CSharpLang
	}

	req, err := c.NewRequest(operation, nil, request)
	if err != nil {
		return nil, err
	}

	bs, err := req.Send()
	if err != nil {
		return nil, err
	}

	var resp FetchResponse
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
