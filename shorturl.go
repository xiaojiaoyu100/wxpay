package wxpay

import (
	"encoding/xml"
)

// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_9&index=10

const (
	shortURL = "https://api.mch.weixin.qq.com/tools/shorturl"
)

// ShortURLRequest 转换短链接请求
type ShortURLRequest struct {
	XMLName  xml.Name `xml:"xml"`
	AppID    string   `xml:"appid,omitempty"`
	MchID    string   `xml:"mch_id,omitempty"`
	NonceStr string   `xml:"nonce_str,omitempty"`
	LongURL  string   `xml:"long_url,omitempty"`
	Sign     string   `xml:"sign,omitempty"`
	SignType string   `xml:"sign_type,omitempty"`
}

// ShortURLResponse 转换短链接回复
type ShortURLResponse struct {
	Meta
	AppID    string `xml:"appid"`
	MchID    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
	ShortURL string `xml:"short_url"`
}

// ShortURL 二维码转短链接
func (c *Client) ShortURL(request *ShortURLRequest) (*ShortURLResponse, error) {
	request.MchID = c.mchID
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response ShortURLResponse
	_, err := c.request(shortURL, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
