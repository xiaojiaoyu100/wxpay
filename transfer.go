package wxpay

import "encoding/xml"

// https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_2

const (
	transferURL = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
)

// TransferRequest 企业付款参数
type TransferRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"mch_appid,omitempty"`
	MchID          string   `xml:"mchid,omitempty"`
	DeviceInfo     string   `xml:"device_info,omitempty"`
	NonceStr       string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	PartnerTradeNo string   `xml:"partner_trade_no,omitempty"` // 商户订单号
	OpenID         string   `xml:"openid,omitempty"`
	CheckName      string   `xml:"check_name,omitempty"` // NO_CHECK：不校验真实姓名, FORCE_CHECK：强校验真实姓名
	ReUserName     string   `xml:"re_user_name,omitempty"`
	Amount         string   `xml:"amount,omitempty"` // 单位分
	Desc           string   `xml:"desc,omitempty"`   // 企业付款备注
	SpBillCreateIP string   `xml:"spbill_create_ip,omitempty"`
}

// TransferResponse 企业付款回复
type TransferResponse struct {
	Meta
	AppID          string `xml:"mch_appid"`
	MchID          string `xml:"mchid"`
	DeviceInfo     string `xml:"device_info"`
	NonceStr       string `xml:"nonce_str"`
	PartnerTradeNo string `xml:"partner_trade_no"` // 商户订单号
	PaymentNo      string `xml:"payment_no"`       // 微信付款单号
	PaymentTime    string `xml:"payment_time"`     // 付款成功时间
}

// Transfer 企业付款到零钱
// NO_AUTH
// AMOUNT_LIMIT
//
func (c *Client) Transfer(request *TransferRequest) (*TransferResponse, error) {
	request.MchID = c.mchID
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response TransferResponse
	_, err := c.request(transferURL, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
