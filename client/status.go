package client

import (
	"ccavenue/aescbc"
	"encoding/hex"
	"encoding/json"
	"strings"
)

type StatusFilter struct {
	OrderNo string `json:"order_no,omitempty"`
}

func (f StatusFilter) Encode() (string, error) {

	jsonBytes, err := json.Marshal(f)
	if err != nil {
		return "", err
	}

	encReqBytes, err := aescbc.NewCrypter().Encrypt(jsonBytes)
	if err != nil {
		return "", err
	}

	encStr := strings.ToUpper(hex.EncodeToString(encReqBytes))

	return encStr, nil
}

func (f StatusFilter) Command() string {
	return "orderStatusTracker"
}
