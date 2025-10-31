package client

import (
	"ccavenue/aescbc"
	"encoding/hex"
	"encoding/json"
	"strings"
)

type OrderLookupFilter struct {
	FromDate       string  `json:"from_date"`
	PageNumber     int     `json:"page_number"`
	OrderMaxAmount float64 `json:"order_max_amount"`
}

func (f OrderLookupFilter) Encode() (string, error) {

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

func (f OrderLookupFilter) Command() string {
	return "orderLookup"
}
