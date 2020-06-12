package wxpay

import "encoding/xml"

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
const (
	closeOrderURL = "https://api.mch.weixin.qq.com/pay/closeorder"
)

// CloseOrderRequest 关闭订单
type CloseOrderRequest struct {
	XMLName    xml.Name `xml:"xml"`
	AppID      string   `xml:"appid,omitempty"`
	MchID      string   `xml:"mch_id,omitempty"`
	OutTradeNo string   `xml:"out_trade_no,omitempty"`
	NonceStr   string   `xml:"nonce_str,omitempty"`
	Sign       string   `xml:"sign,omitempty"`
	SignType   string   `xml:"sign_type,omitempty"`
}

// CloseOrderResponse 关闭订单回复
type CloseOrderResponse struct {
	Meta
	AppID    string `xml:"appid"`
	MchID    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
}

// CloseOrder 关闭订单
func (c *Client) CloseOrder(request *CloseOrderRequest) (*CloseOrderResponse, error) {
	request.MchID = c.mchID
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response CloseOrderResponse
	_, err := c.request(closeOrderURL, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
