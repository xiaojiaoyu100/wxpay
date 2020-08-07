package wxpay

import "encoding/xml"

// https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_3

const (
	transferInfoURL = "https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo"
)

const (
	transferStatusSuccess    = "SUCCESS"
	transferStatusFail       = "FAILED"
	transferStatusProcessing = "PROCESSING"
)

// TransferInfoRequest 查询企业付款参数
type TransferInfoRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid,omitempty"`
	MchID          string   `xml:"mch_id,omitempty"`
	NonceStr       string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	PartnerTradeNo string   `xml:"partner_trade_no,omitempty"` // 商户订单号
}

// TransferInfoResponse 企业付款响应
type TransferInfoResponse struct {
	Meta
	AppID          string `xml:"appid"`
	MchID          string `xml:"mch_id"`
	PartnerTradeNo string `xml:"partner_trade_no"` // 商户订单号
	DetailId       string `xml:"detail_id"`        // 付款单号
	Status         string `xml:"status"`           // 转账状态
	Reason         string `xml:"reason"`           // 失败原因
	OpenID         string `xml:"openid"`           // 收款用户openId
	TransferName   string `xml:"transfer_name"`    // 收款用户姓名
	PaymentAmount  string `xml:"payment_amount"`   // 收款金额
	TransferTime   string `xml:"transfer_time"`    // 转账时间
	PaymentTime    string `xml:"payment_time"`     // 付款成功时间
	Remark         string `xml:"desc"`             // 企业付款备注
}

// GetTransferInfo 企业付款查询
func (c *Client) GetTransferInfo(request *TransferInfoRequest) (*TransferInfoResponse, error) {
	request.MchID = c.mchID
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response TransferInfoResponse
	_, err := c.request(transferInfoURL, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (resp *TransferInfoResponse) PaidSuccess() bool {
	return resp.Status == transferStatusSuccess
}

func (resp *TransferInfoResponse) PaidFail() bool {
	return resp.Status == transferStatusFail
}

func (resp *TransferInfoResponse) PayProcessing() bool {
	return resp.Status == transferStatusProcessing
}
