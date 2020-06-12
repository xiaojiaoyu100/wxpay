package wxpay

import "encoding/xml"

const (
	reverseURL = "https://api.mch.weixin.qq.com/secapi/pay/reverse"
)

// ReverseRequest 撤销交易请求
type ReverseRequest struct {
	XMLName       xml.Name `xml:"xml"`
	AppID         string   `xml:"appid,omitempty"`
	MchID         string   `xml:"mch_id,omitempty"`
	TransactionID string   `xml:"transaction_id,omitempty"`
	OutTradeNo    string   `xml:"out_trade_no,omitempty"`
	NonceStr      string   `xml:"nonce_str,omitempty"`
	Sign          string   `xml:"sign,omitempty"`
	SignType      string   `xml:"sign_type,omitempty"`
}

// ReverseResponse 撤销交易回复
type ReverseResponse struct {
	Meta
	AppID    string `xml:"appid"`
	MchID    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
	Recall   string `xml:"recall"`
}

// Reverse 撤销交易　
// https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_11&index=3
func (c *Client) Reverse(request *ReverseRequest) (*ReverseResponse, error) {
	request.MchID = c.mchID
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response ReverseResponse
	_, err := c.request(reverseURL, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
