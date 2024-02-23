package client

import (
	"ccavenue/aescbc"
	"ccavenue/config"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type APIClient struct {
	Host         string
	Client       *http.Client
	AccessCode   string
	RequestType  string
	ResponseType string
	Version      string
}

type Filter interface {
	Encode() (string, error)
	Command() string
}

type CCAvenueParams struct {
	EncRequest   string `json:"enc_request,omitempty"`
	AccessCode   string `json:"access_code,omitempty"`
	Command      string `json:"command,omitempty"`
	RequestType  string `json:"request_type,omitempty"`
	ResponseType string `json:"response_type,omitempty"`
	Version      string `json:"version,omitempty"`
}

const TIMEOUT = 15 * time.Second

func NewClient(cfg config.Config, version string) (*APIClient, error) {

	return &APIClient{
		Host:         cfg.Host,
		Client:       &http.Client{Timeout: TIMEOUT},
		AccessCode:   cfg.AccessCode,
		RequestType:  "JSON",
		ResponseType: "JSON",
		Version:      version,
	}, nil
}

func (c *APIClient) Post(f Filter) (*[]byte, error) {

	encStr, err := f.Encode()
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("enc_request", encStr)
	q.Add("access_code", c.AccessCode)
	q.Add("command", f.Command())
	q.Add("request_type", c.RequestType)
	q.Add("response_type", c.ResponseType)
	q.Add("version", c.Version)

	req, err := http.NewRequest("POST", c.Host, strings.NewReader(q.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	encResponse := response.Body
	defer encResponse.Close()

	rawQuery, err := io.ReadAll(encResponse)
	if err != nil {
		return nil, err
	}
	values, err := url.ParseQuery(string(rawQuery))
	if err != nil {
		return nil, err
	}

	if values["status"][0] == "0" {

		encArr, ok := values["enc_response"]
		if !ok || len(encArr) == 0 {
			return nil, fmt.Errorf("invalid encrypted response: %v, %d\n %+v", ok, len(encArr), encArr)
		}

		respStr, err := decode(encArr[0])
		if err != nil {
			return nil, err
		}

		return respStr, nil
	}

	return nil, fmt.Errorf("error parsing response%s", fmt.Sprintf("%+v", values))
}

func decode(encStr string) (*[]byte, error) {

	payload := strings.TrimSpace(encStr)
	buf, err := hex.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	resp, err := aescbc.NewCrypter().Decrypt(buf)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
