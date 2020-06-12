package wxpay

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
	"time"
)

const (
	unifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

// UnifiedOrderRequest 下单请求
type UnifiedOrderRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid,omitempty"`
	MchID          string   `xml:"mch_id,omitempty"`
	DeviceInfo     string   `xml:"device_info,omitempty"`
	NonceStr       string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	SignType       string   `xml:"sign_type,omitempty"`
	Body           string   `xml:"body,omitempty"`
	Detail         string   `xml:"detail,omitempty"`
	Attach         string   `xml:"attach,omitempty"`
	OutTradeNo     string   `xml:"out_trade_no,omitempty"`
	FeeType        string   `xml:"fee_type,omitempty"`
	TotalFee       int64    `xml:"total_fee,omitempty"`
	SpBillCreateIP string   `xml:"spbill_create_ip,omitempty"`
	TimeStart      string   `xml:"time_start,omitempty"`
	TimeExpire     string   `xml:"time_expire,omitempty"`
	GoodsTag       string   `xml:"goods_tag,omitempty"`
	NotifyURL      string   `xml:"notify_url,omitempty"`
	TradeType      string   `xml:"trade_type,omitempty"`
	ProductID      string   `xml:"product_id,omitempty"`
	LimitPay       string   `xml:"limit_pay,omitempty"`
	OpenID         string   `xml:"openid,omitempty"`
	SceneInfo      string   `xml:"scene_info,omitempty"`
}

// TimeExpire 默认时间为30分钟
func TimeExpire() string {
	now := time.Now().UTC()
	duration := time.Duration(30*time.Minute) + time.Duration(8*time.Hour)
	return now.Add(duration).Format("20060102150405")
}

// UnifiedOrderResponse 下单回复
type UnifiedOrderResponse struct {
	Meta
	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	PrepayID   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
	CodeURL    string `xml:"code_url"`
	MWebURL    string `xml:"mweb_url"`
}

// UnifiedOrder 下单
func (c *Client) UnifiedOrder(request *UnifiedOrderRequest) (*UnifiedOrderResponse, error) {
	request.MchID = c.mchID
	request.NonceStr = nonceStr()
	request.TimeExpire = TimeExpire()
	request.Sign = signStruct(request, c.apiKey)

	if len(request.Body) == 0 {
		return nil, errors.New("body is zero")
	}

	if len(request.OutTradeNo) == 0 {
		return nil, errors.New("out_trade_no is zero")
	}

	if request.TotalFee <= 0 {
		return nil, errors.New("wrong total_fee")
	}

	if len(request.SpBillCreateIP) == 0 {
		return nil, errors.New("spbill_create_ip is zero")
	}

	if len(request.NotifyURL) == 0 {
		return nil, errors.New("notify_url is zero")
	}

	switch request.TradeType {
	case TradeTypeNative:
	case TradeTypeJs:
		if len(request.OpenID) == 0 {
			return nil, errors.New("openid is zero")
		}
	case TradeTypeMWeb:
		if len(request.SceneInfo) == 0 {
			return nil, errors.New("scene_info is zero")
		}
	default:
		return nil, errors.New("wrong trade_type")
	}

	var response UnifiedOrderResponse
	_, err := c.request(unifiedOrderURL, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// BrandWCPayRequest 微信内H5调起支付
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7&index=6
type BrandWCPayRequest struct {
	AppID     string `xml:"appId" json:"appId"`
	Timestamp string `xml:"timeStamp" json:"timeStamp"`
	NonceStr  string `xml:"nonceStr" json:"nonceStr"`
	Package   string `xml:"package" json:"package"`
	SignType  string `xml:"signType" json:"signType"`
	PaySign   string `xml:"paySign" json:"paySign"`
}

// GetBrandWCPayRequest 生成微信内h5支付调用支付字符串
func (c *Client) GetBrandWCPayRequest(resp *UnifiedOrderResponse) string {
	brandWCPayRequest := &BrandWCPayRequest{
		AppID:     resp.AppID,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  nonceStr(),
		Package:   "prepay_id=" + resp.PrepayID,
		SignType:  "MD5",
	}
	brandWCPayRequest.PaySign = signStruct(brandWCPayRequest, c.apiKey)
	bytes, err := json.Marshal(brandWCPayRequest)
	if err != nil {
		globalLogger.printf("%s marshal err: %s", "GetBrandWCPayRequest: ", err.Error())
	}
	globalLogger.printf("GetBrandWCPayRequest: %s", string(bytes))
	return string(bytes)
}
