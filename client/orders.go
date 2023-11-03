package client

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

type Order struct {
	ReferenceNo          int64   `json:"reference_no,omitempty"`
	OrderNo              string  `json:"order_no,omitempty"`
	OrderCurrncy         string  `json:"order_currncy,omitempty"`
	OrderAmt             float64 `json:"order_amt,omitempty"`
	OrderDateTime        string  `json:"order_date_time,omitempty"`
	OrderBillName        string  `json:"order_bill_name,omitempty"`
	OrderBillAddress     string  `json:"order_bill_address,omitempty"`
	OrderBillZip         string  `json:"order_bill_zip,omitempty"`
	OrderBillTel         string  `json:"order_bill_tel,omitempty"`
	OrderBillEmail       string  `json:"order_bill_email,omitempty"`
	OrderBillCountry     string  `json:"order_bill_country,omitempty"`
	OrderShipName        string  `json:"order_ship_name,omitempty"`
	OrderShipAddress     string  `json:"order_ship_address,omitempty"`
	OrderShipCountry     string  `json:"order_ship_country,omitempty"`
	OrderShipTel         string  `json:"order_ship_tel,omitempty"`
	OrderBillCity        string  `json:"order_bill_city,omitempty"`
	OrderBillState       string  `json:"order_bill_state,omitempty"`
	OrderShipCity        string  `json:"order_ship_city,omitempty"`
	OrderShipState       string  `json:"order_ship_state,omitempty"`
	OrderShipZip         string  `json:"order_ship_zip,omitempty"`
	OrderShipEmail       string  `json:"order_ship_email,omitempty"`
	OrderNotes           string  `json:"order_notes,omitempty"`
	OrderIP              string  `json:"order_ip,omitempty"`
	OrderStatus          string  `json:"order_status,omitempty"`
	OrderFraudStatus     string  `json:"order_fraud_status,omitempty"`
	OrderStatusDateTime  string  `json:"order_status_date_time,omitempty"`
	OrderCaptAmt         float64 `json:"order_capt_amt,omitempty"`
	OrderCardName        string  `json:"order_card_name,omitempty"`
	OrderDeliveryDetails string  `json:"order_delivery_details,omitempty"`
	OrderFeePercValue    float64 `json:"order_fee_perc_value,omitempty"`
	OrderFeePerc         float64 `json:"order_fee_perc,omitempty"`
	OrderFeeFlat         float64 `json:"order_fee_flat,omitempty"`
	OrderGrossAmt        float64 `json:"order_gross_amt,omitempty"`
	OrderDiscount        float64 `json:"order_discount,omitempty"`
	OrderTax             float64 `json:"order_tax,omitempty"`
	OrderBankRefNo       string  `json:"order_bank_ref_no,omitempty"`
	OrderGtwID           string  `json:"order_gtw_id,omitempty"`
	OrderBankResponse    string  `json:"order_bank_response,omitempty"`
	OrderOptionType      string  `json:"order_option_type,omitempty"`
	OrderTDS             string  `json:"order_TDS,omitempty"`
	OrderDeviceType      string  `json:"order_device_type,omitempty"`
	OrderType            string  `json:"order_type,omitempty"`
	ShipDateTime         string  `json:"ship_date_time,omitempty"`
	PaymentMode          string  `json:"payment_mode,omitempty"`
	CardType             string  `json:"card_type,omitempty"`
	SubAccID             string  `json:"sub_acc_id,omitempty"`
	OrderBinCountry      string  `json:"order_bin_country,omitempty"`
	OrderStlmtDate       string  `json:"order_stlmt_date,omitempty"`
	CardEnrolled         string  `json:"card_enrolled,omitempty"`
	MerchantParam1       string  `json:"merchant_param1,omitempty"`
	MerchantParam2       string  `json:"merchant_param2,omitempty"`
	MerchantParam3       string  `json:"merchant_param3,omitempty"`
	MerchantParam4       string  `json:"merchant_param4,omitempty"`
	MerchantParam5       string  `json:"merchant_param5,omitempty"`
	OrderBankArnNo       string  `json:"order_bank_arn_no,omitempty"`
	OrderSplitPayout     string  `json:"order_split_payout,omitempty"`
	EmiIssuingBank       string  `json:"emi_issuing_bank,omitempty"`
	TenureDuration       string  `json:"tenure_duration,omitempty"`
	EmiDiscountType      string  `json:"emi_discount_type,omitempty"`
	EmiDiscountValue     string  `json:"emi_discount_value,omitempty"`
}

type OrderResponse struct {
	OrderStatusList []Order `json:"order_Status_List,omitempty"`
	PageCount       int     `json:"page_count,omitempty"`
	TotalRecords    int     `json:"total_records,omitempty"`
	ErrorDesc       string  `json:"error_desc,omitempty"`
	ErrorCode       string  `json:"error_code,omitempty"`
}

func (c *APIClient) Orders(fromDate, toDate string) (*OrderResponse, error) {

	data := CCAvenueData{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	response, err := c.Post("orderLookup", data)
	if err != nil {
		return nil, err
	}

	query := response.Body
	defer query.Close()

	rawQuery, err := io.ReadAll(query)
	if err != nil {
		return nil, err
	}
	values, err := url.ParseQuery(string(rawQuery))
	if err != nil {
		return nil, err
	}

	fmt.Printf("Values: %+v\n", values)

	if values["status"][0] == "0" {

		payload := strings.TrimSpace(values["enc_response"][0])
		buf, err := hex.DecodeString(payload)
		if err != nil {
			return nil, err
		}

		jsonBytes, err := c.Decrypt(buf)
		if err != nil {
			return nil, err
		}

		orders := OrderResponse{}
		if err := json.Unmarshal(jsonBytes, &orders); err != nil {
			return nil, err
		}
		return &orders, nil
	}

	return nil, errors.New(string(rawQuery))
}
