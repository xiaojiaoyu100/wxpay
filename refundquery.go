package wxpay

import (
	"encoding/xml"
	"strconv"
)

// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_5
const (
	refundQueryURL = "https://api.mch.weixin.qq.com/pay/refundquery"
)

const (
	// RefundStatusSuccess 退款成功
	RefundStatusSuccess = "SUCCESS"
	// RefundStatusRefundClose 退款关闭
	RefundStatusRefundClose = "REFUNDCLOSE"
	// RefundStatusProcessing 退款处理中
	RefundStatusProcessing = "PROCESSING"
	// RefundStatusChange 退款异常，退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，可前往商户平台（pay.weixin.qq.com）-交易中心，手动处理
	RefundStatusChange = "CHANGE"
)

// RefundQueryRequest 退款查询请求
type RefundQueryRequest struct {
	XMLName       xml.Name `xml:"xml"`
	AppID         string   `xml:"appid,omitempty"`
	MchID         string   `xml:"mch_id,omitempty"`
	NonceStr      string   `xml:"nonce_str,omitempty"`
	Sign          string   `xml:"sign,omitempty"`
	SignType      string   `xml:"sign_type,omitempty"`
	TransactionID string   `xml:"transaction_id,omitempty"`
	OutTradeNo    string   `xml:"out_trade_no,omitempty"`
	OutRefundNo   string   `xml:"out_refund_no,omitempty"`
	RefundID      string   `xml:"refund_id,omitempty"`
	Offset        int      `xml:"offset,omitempty"`
}

// RefundDetail 退款详情
type RefundDetail struct {
	OutRefundNo         string // 商户退款单号
	RefundID            string // 微信退款单号
	RefundChannel       string // 退款渠道
	RefundFee           int64  // 申请退款金额
	SettlementRefundFee int64  // 退款金额
	CouponRefundFee     int64  // 总代金券退款金额
	CouponRefundCount   int    // 退款代金券使用数量
	RefundStatus        string // 退款状态
	RefundAccount       string // 退款资金来源
	RefundRecvAccount   string // 退款入账账户
	RefundSuccessTime   string // 退款成功时间
}

// RefundQueryResponse 退款查询回复
type RefundQueryResponse struct {
	Meta
	AppID              string          `xml:"appid"`
	MchID              string          `xml:"mch_id"`
	NonceStr           string          `xml:"nonce_str"`
	Sign               string          `xml:"sign"`
	TotalRefundCount   int             `xml:"total_refund_count"` // 订单总退款次数, 订单总共已发生的部分退款次数，当请求参数传入offset后有返回
	TransactionID      string          `xml:"transaction_id"`
	OutTradeNo         string          `xml:"out_trade_no"`
	TotalFee           int64           `xml:"total_fee"`
	SettlementTotalFee int64           `xml:"settlement_total_fee"`
	FeeType            string          `xml:"fee_type"`
	CashFee            int64           `xml:"cash_fee"`
	RefundCount        int             `xml:"refund_count"`
	RefundDetails      []*RefundDetail `xml:"-"`
}

// RefundQuery 退款查询
func (c *Client) RefundQuery(request *RefundQueryRequest) (*RefundQueryResponse, error) {
	request.MchID = c.mchID
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response RefundQueryResponse
	body, err := c.request(refundQueryURL, request, &response)
	if err != nil {
		return nil, err
	}
	tempMap := make(Map)
	if err := xml.Unmarshal(body, &tempMap); err != nil {
		return nil, err
	}

	response.RefundDetails = make([]*RefundDetail, 0, response.RefundCount)
	for i := 0; i < response.RefundCount; i++ {

		rd := new(RefundDetail)
		response.RefundDetails = append(response.RefundDetails, rd)

		index := strconv.Itoa(i)
		key := "out_refund_no_" + index
		if val, ok := tempMap[key]; ok {
			rd.OutRefundNo = val
		}

		key = "refund_id_" + index
		if val, ok := tempMap[key]; ok {
			rd.RefundID = val
		}

		key = "refund_channel_" + index
		if val, ok := tempMap[key]; ok {
			rd.RefundChannel = val
		}

		key = "refund_fee_" + index
		if val, ok := tempMap[key]; ok {
			intVal, _ := strconv.ParseInt(val, 10, 0)
			rd.RefundFee = intVal
		}

		key = "settlement_refund_fee_" + index
		if val, ok := tempMap[key]; ok {
			intVal, _ := strconv.ParseInt(val, 10, 0)
			rd.SettlementRefundFee = intVal
		}

		key = "coupon_refund_fee_" + index
		if val, ok := tempMap[key]; ok {
			intVal, _ := strconv.ParseInt(val, 10, 0)
			rd.CouponRefundFee = intVal
		}

		key = "coupon_refund_count_" + index
		if val, ok := tempMap[key]; ok {
			intVal, _ := strconv.ParseInt(val, 10, 0)
			rd.CouponRefundCount = int(intVal)
		}

		key = "refund_status_" + index
		if val, ok := tempMap[key]; ok {
			rd.RefundStatus = val
		}

		key = "refund_account_" + index
		if val, ok := tempMap[key]; ok {
			rd.RefundAccount = val
		}

		// 注意这里是微信的拼写错误
		key = "refund_recv_accout_" + index
		if val, ok := tempMap[key]; ok {
			rd.RefundRecvAccount = val
		}

		key = "refund_success_time_" + index
		if val, ok := tempMap[key]; ok {
			rd.RefundSuccessTime = val
		}
	}

	return &response, nil
}
