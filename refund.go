package wxpay

import "encoding/xml"

// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_4

const (
	refundURL = "https://api.mch.weixin.qq.com/secapi/pay/refund"
)

// RefundRequest 退款请求
type RefundRequest struct {
	XMLName       xml.Name `xml:"xml"`
	AppID         string   `xml:"appid,omitempty"`
	MchID         string   `xml:"mch_id,omitempty"`
	NonceStr      string   `xml:"nonce_str,omitempty"`
	Sign          string   `xml:"sign,omitempty"`
	SignType      string   `xml:"sign_type,omitempty"`
	TransactionID string   `xml:"transaction_id,omitempty"`
	OutTradeNo    string   `xml:"out_trade_no,omitempty"`
	OutRefundNo   string   `xml:"out_refund_no,omitempty"`
	TotalFee      int64    `xml:"total_fee,omitempty"`
	RefundFee     int64    `xml:"refund_fee,omitempty"`
	RefundFeeType string   `xml:"refund_fee_type,omitempty"`
	RefundDesc    string   `xml:"refund_desc,omitempty"`
	RefundAccount string   `xml:"refund_account,omitempty"`
	NotifyURL     string   `xml:"notify_url,omitempty"` // 可以不传
}

// RefundResponse 退款回复
type RefundResponse struct {
	Meta
	AppID               string `xml:"appid"`
	MchID               string `xml:"mch_id"`
	NonceStr            string `xml:"nonce_str"`
	Sign                string `xml:"sign"`
	TransactionID       string `xml:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no"`
	OutRefundNo         string `xml:"out_refund_no"`
	RefundID            string `xml:"refund_id"`
	RefundFee           int64  `xml:"refund_fee"`
	SettlementRefundFee int64  `xml:"settlement_refund_fee"`
	TotalFee            int64  `xml:"total_fee"`
	SettlementTotalFee  int64  `xml:"settlement_total_fee"`
	FeeType             string `xml:"fee_type"`
	CashFee             int64  `xml:"cash_fee"`
	CashFeeType         string `xml:"cash_fee_type"`
	CashRefundFee       int64  `xml:"cash_refund_fee"`
	CouponRefundFee     int64  `xml:"coupon_refund_fee"`
	CouponRefundCount   int    `xml:"coupon_refund_count"`
}

// Refund 退款
func (c *Client) Refund(request *RefundRequest) (*RefundResponse, error) {
	request.MchID = c.mchID
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response RefundResponse
	_, err := c.request(refundURL, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
