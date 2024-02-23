package client

import (
	"ccavenue/aescbc"
	"encoding/hex"
	"encoding/json"
	"strings"
)

type PayoutFilter struct {
	SettlementDate string `json:"settlement_date,omitempty"`
}

func (f PayoutFilter) Encode() (string, error) {

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

func (f PayoutFilter) Command() string {
	return "payoutSummary"
}
