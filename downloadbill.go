package wxpay

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	downloadBillURL = "https://api.mch.weixin.qq.com/pay/downloadbill"
)

const (
	// BillTypeAll 返回当日所有订单信息，默认值
	BillTypeAll            = "ALL"
	// BillTypeSuccess 返回当日成功支付的订单
	BillTypeSuccess        = "SUCCESS"
	// BillTypeRefund 返回当日退款订单
	BillTypeRefund         = "REFUND"
	// BillTypeRechargeRefund 返回当日充值退款订单
	BillTypeRechargeRefund = "RECHARGE_REFUND"
)

// DownloadBillRequest 下载账单请求
// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_6
type DownloadBillRequest struct {
	XMLName  xml.Name `xml:"xml"`
	AppID    string   `xml:"appid,omitempty"`
	MchID    string   `xml:"mch_id,omitempty"`
	NonceStr string   `xml:"nonce_str,omitempty"`
	Sign     string   `xml:"sign,omitempty"`
	SignType string   `xml:"sign_type,omitempty"`
	BillDate string   `xml:"bill_date,omitempty"`
	BillType string   `xml:"bill_type,omitempty"`
	TarType  string   `xml:"tar_type,omitempty"`
}

// DownloadBillResponse 下载账单回复
type DownloadBillResponse struct {
	EntryList  []*BillEntry
	Statistics []*BillStatistics
	Bill       string
}

// DownloadBill 下载账单
func (c *Client) DownloadBill(request *DownloadBillRequest) (*DownloadBillResponse, error) {
	const (
		max = 1000 * time.Millisecond
	)
	var (
		tempDelay time.Duration
		tryNum    = 0
		buf       bytes.Buffer
	)

tryLoop:
	for {

		if tryNum >= 3 {
			break tryLoop
		}

		if tempDelay == 0 {
			tempDelay = 100 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if tempDelay > max {
			tempDelay = max
		}

		request.MchID = c.mchID
		request.NonceStr = nonceStr()
		request.TarType = "GZIP"
		request.Sign = signStruct(request, c.apiKey)

		body, err := xml.Marshal(&request)
		if err != nil {
			globalLogger.printf("Marshal err: %v", err)
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, downloadBillURL, bytes.NewBuffer(body))
		if err != nil {
			globalLogger.printf("NewRequest err: %v", err)
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req = req.WithContext(ctx)

		resp, err := selectedClient(downloadBillURL).Do(req)
		tryNum++
		switch {
		case err != nil:
			if shouldRetry(err) {
				notifyAsync("downloadbill err: ", err)
				time.Sleep(tempDelay)
				continue tryLoop
			}
			globalLogger.printf("do err: %v", err)
			return nil, err
		default:
			defer resp.Body.Close()

			responseBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				globalLogger.printf("do err: %v", err)
				return nil, err
			}

			switch request.TarType {
			case "GZIP":
				reader, err := gzip.NewReader(bytes.NewBuffer(responseBytes))
				switch err {
				case nil:
					defer reader.Close()
					if _, err := io.Copy(&buf, reader); err != nil {
						globalLogger.printf("Copy err: %v", err)
						return nil, err
					}
				case gzip.ErrHeader:
					buf = *bytes.NewBuffer(responseBytes)
				default:
					globalLogger.printf("NewReader err: %v", err)
					return nil, err
				}
			default:
				buf = *bytes.NewBuffer(responseBytes)
			}

			var response struct {
				ReturnCode string `xml:"return_code"`
				ReturnMsg  string `xml:"return_msg"`
			}
			if err = xml.Unmarshal(buf.Bytes(), &response); err == nil {
				switch response.ReturnMsg {
				case billNoExistErr.Error():
					return nil, billNoExistErr
				case "SYSTEMERROR",
					"CompressGZip Error",
					"UnCompressGZip Error":
					notifyAsync("downloadbill err: ", err)
					time.Sleep(tempDelay)
					continue tryLoop
				default:
					return nil, errors.New(response.ReturnMsg)
				}
			} else {
				break tryLoop
			}
		}
	}

	var response DownloadBillResponse

	response.Bill = buf.String()
	response.Bill = strings.Replace(response.Bill, "`", "", -1)
	pos := strings.LastIndex(response.Bill, "%")

	firstPart := response.Bill[0 : pos+1]
	secondPart := response.Bill[pos+1:]

	orderList, err := csv.NewReader(strings.NewReader(firstPart)).ReadAll()
	if err != nil {
		globalLogger.printf("csv.NewReader err: %v", err)
		return nil, err
	}

	response.EntryList = make([]*BillEntry, 0)

	if len(orderList) >= 2 {
		for _, order := range orderList[1:] {
			if len(order) == 24 {
				entry := new(BillEntry)
				response.EntryList = append(response.EntryList, entry)

				entry.TimeEnd = order[0]
				entry.AppID = order[1]
				entry.MchID = order[2]
				entry.SubMchID = order[3]
				entry.DeviceInfo = order[4]
				entry.TransactionID = order[5]
				entry.OutTradeNo = order[6]
				entry.OpenID = order[7]
				entry.TradeType = order[8]
				entry.TradeStatus = order[9]
				entry.BankType = order[10]
				entry.FeeType = order[11]
				entry.TotalFee = order[12]
				entry.CouponFee = order[13]
				entry.RefundID = order[14]
				entry.OutRefundNo = order[15]
				entry.RefundFee = order[16]
				entry.CouponRefundFee = order[17]
				entry.RefundChannel = order[18]
				entry.RefundStatus = order[19]
				entry.Body = order[20]
				entry.Attach = order[21]
				entry.HandlingCharge = order[22]
				entry.Rate = order[23]
			}
		}
	}

	statisticsList, err := csv.NewReader(strings.NewReader(secondPart)).ReadAll()
	if err != nil {
		globalLogger.printf("csv.NewReader err: %v", err)
		return nil, err
	}

	response.Statistics = make([]*BillStatistics, 0)

	if len(statisticsList) >= 2 {
		for _, statistics := range statisticsList[1:] {
			if len(statistics) == 5 {
				s := new(BillStatistics)
				response.Statistics = append(response.Statistics, s)

				s.TradeOrderCount = statistics[0]
				s.TotalBusinessTransaction = statistics[1]
				s.TotalRefundFee = statistics[2]

				s.TotalCouponFee = statistics[3]
				s.TotalHandlingCharge = statistics[4]
			}
		}
	}

	return &response, nil

}

const (
	// BillTradeStatusSuccess 成功
	BillTradeStatusSuccess = "SUCCESS"
	// BillTradeStatusRefund 退款
	BillTradeStatusRefund  = "REFUND"
	// BillTradeStatusRevoked 撤销
	BillTradeStatusRevoked = "REVOKED"
)

// BillStatistics 账单统计
type BillStatistics struct {
	TradeOrderCount          string `csv:"总交易单数"`
	TotalBusinessTransaction string `csv:"总交易额"`
	TotalRefundFee           string `csv:"总退款金额"`
	TotalCouponFee           string `csv:"总企业红包退款金额"`
	TotalHandlingCharge      string `csv:"手续费总金额"`
}

// BillEntry 账单条目
type BillEntry struct {
	TimeEnd         string `csv:"交易时间"`
	AppID          string `csv:"公众账号ID"`
	MchID           string `csv:"商户号"`
	SubMchID        string `csv:"子商户号"`
	DeviceInfo      string `csv:"设备号"`
	TransactionID   string `csv:"微信订单号"`
	OutTradeNo      string `csv:"商户订单号"`
	OpenID         string `csv:"用户标识"`
	TradeType       string `csv:"交易类型"`
	TradeStatus     string `csv:"交易状态"`
	BankType        string `csv:"付款银行"`
	FeeType         string `csv:"货币种类"`
	TotalFee        string `csv:"总金额"`
	CouponFee       string `csv:"企业红包金额"`
	RefundID       string `csv:"微信退款单号"`
	OutRefundNo     string `csv:"商户退款单号"`
	RefundFee       string `csv:"退款金额"`
	CouponRefundFee string `csv:"企业红包退款金额"`
	RefundChannel   string `csv:"退款类型"`
	RefundStatus    string `csv:"退款状态"`
	Body            string `csv:"商品名称"`
	Attach          string `csv:"商户数据包"`
	HandlingCharge  string `csv:"手续费"`
	Rate            string `csv:"费率"`
}
