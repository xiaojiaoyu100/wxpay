package wxpay

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
)

// PaidNotifyRequest 付款通知请求
type PaidNotifyRequest struct {
	XMLName xml.Name `xml:"xml"`
	Meta
	AppID             string `xml:"appid"`
	MchID              string `xml:"mch_id"`
	DeviceInfo         string `xml:"device_info"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	SignType           string `xml:"sign_type"`
	OpenID             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	BankType           string `xml:"bank_type"`
	TotalFee           int64  `xml:"total_fee"`
	SettlementTotalFee int64  `xml:"settlement_total_fee"` // 应结订单金额
	FeeType            string `xml:"fee_type"`
	CashFee            int64  `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          int64  `xml:"coupon_fee"`     // 总代金券金额
	CouponCount        int    `xml:"coupon_count"`   // 代金券使用数量
	TransactionID      string `xml:"transaction_id"` // 微信支付订单号
	OutTradeNo         string `xml:"out_trade_no"`   // 商户订单号
	Attach             string `xml:"attach"`         // 商家数据包
	TimeEnd            string `xml:"time_end"`       // 支付完成时间
}

// PaidNotifyResponse 付款通知回复
type PaidNotifyResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

// PaidVerify 付款通知
func (c *Client) PaidVerify(body []byte) (*PaidNotifyRequest, *PaidNotifyResponse, error) {
	var (
		notifyRequest  PaidNotifyRequest
		notifyResponse PaidNotifyResponse
	)
	err := xml.Unmarshal(body, &notifyRequest)
	if err != nil {
		return nil, nil, err
	}

	if !notifyRequest.ResultCodeSuccess() {
		return nil, nil, errors.New("业务结果不成功")
	}

	err = checkSign(body, c.apiKey)
	if err != nil {
		return nil, nil, err
	}

	notifyResponse.ReturnCode = "SUCCESS"
	notifyResponse.ReturnMsg = "OK"
	return &notifyRequest, &notifyResponse, nil
}

// PaidNotifyVerify 付款通知验证
func (c *Client) PaidNotifyVerify(request *http.Request) (*PaidNotifyRequest, *PaidNotifyResponse, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, nil, err
	}
	globalLogger.printf("%s: %s", request.URL.String(), string(body))
	return c.PaidVerify(body)
}
