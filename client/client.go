package client

import (
	"ccavenue/aescbc"
	"ccavenue/config"
	"encoding/hex"
	"encoding/json"
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

type Crypter interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

func (*APIClient) Encrypt(buf []byte) ([]byte, error) {

	return aescbc.NewCrypter().Encrypt(buf)
}

func (*APIClient) Decrypt(buf []byte) ([]byte, error) {

	return aescbc.NewCrypter().Decrypt(buf)
}

type CCAvenueData struct {
	OrderNo    string `json:"order_no,omitempty"`
	OrderEmail string `json:"order_email,omitempty"`
	FromDate   string `json:"from_date,omitempty"`
	ToDate     string `json:"to_date,omitempty"`
	PageNumber int    `json:"page_number,omitempty"`
}

type CCAvenueParams struct {
	EncRequest   string `json:"enc_request,omitempty"`
	AccessCode   string `json:"access_code,omitempty"`
	Command      string `json:"command,omitempty"`
	RequestType  string `json:"request_type,omitempty"`
	ResponseType string `json:"response_type,omitempty"`
	Version      string `json:"version,omitempty"`
}

func NewClient(cfg config.Config, timeout time.Duration) (*APIClient, error) {

	return &APIClient{
		Host:         cfg.Host,
		Client:       &http.Client{Timeout: timeout},
		AccessCode:   cfg.AccessCode,
		RequestType:  "JSON",
		ResponseType: "JSON",
		Version:      "1.1",
	}, nil
}

func (c *APIClient) Post(command string, data CCAvenueData, destPtr any) error {

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encReqBytes, err := c.Encrypt(jsonBytes)
	if err != nil {
		return err
	}

	q := url.Values{}
	q.Add("enc_request", strings.ToUpper(hex.EncodeToString(encReqBytes)))
	q.Add("access_code", c.AccessCode)
	q.Add("command", command)
	q.Add("request_type", c.RequestType)
	q.Add("response_type", c.ResponseType)
	q.Add("version", c.Version)

	req, err := http.NewRequest("POST", c.Host, strings.NewReader(q.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	query := response.Body
	defer query.Close()

	rawQuery, err := io.ReadAll(query)
	if err != nil {
		return err
	}
	values, err := url.ParseQuery(string(rawQuery))
	if err != nil {
		return err
	}

	if values["status"][0] == "0" {

		payload := strings.TrimSpace(values["enc_response"][0])
		buf, err := hex.DecodeString(payload)
		if err != nil {
			return err
		}

		jsonBytes, err := c.Decrypt(buf)
		if err != nil {
			return err
		}

		fmt.Println(string(jsonBytes))

		if err := json.Unmarshal(jsonBytes, destPtr); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("%s", fmt.Sprintf("%+v", values))
}
