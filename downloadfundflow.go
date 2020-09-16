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
	downloadFundFlowURL = "https://api.mch.weixin.qq.com/pay/downloadfundflow"
)

// 账单的资金来源账户
const (
	// AccountTypeBasic 基本账户
	AccountTypeBasic = "Basic"
	// AccountTypeOperation 运营账户
	AccountTypeOperation = "Operation"
	// AccountTypeFees 手续费账户
	AccountTypeFees = "Fees"
)

// DownloadFundFlowRequest 下载资金账单请求
// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_18
type DownloadFundFlowRequest struct {
	XMLName     xml.Name `xml:"xml"`
	AppID       string   `xml:"appid,omitempty"`
	MchID       string   `xml:"mch_id,omitempty"`
	NonceStr    string   `xml:"nonce_str,omitempty"`
	Sign        string   `xml:"sign,omitempty"`
	SignType    string   `xml:"sign_type,omitempty"`
	BillDate    string   `xml:"bill_date,omitempty"`
	AccountType string   `xml:"account_type,omitempty"`
	TarType     string   `xml:"tar_type,omitempty"`
}

// DownloadFundFlowResponse 下载账单回复
type DownloadFundFlowResponse struct {
	EntryList  []*FundFlowEntry
	Statistics []*FundFlowStatistics
	FundFlow   string
}

// DownloadFundFlow 下载资金账单
func (c *Client) DownloadFundFlow(request *DownloadFundFlowRequest) (*DownloadFundFlowResponse, error) {
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
		request.SignType = SignTypeSHA256
		request.Sign = signStruct(request, c.apiKey)

		body, err := xml.Marshal(&request)
		if err != nil {
			globalLogger.printf("Marshal err: %v", err)
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, downloadFundFlowURL, bytes.NewBuffer(body))
		if err != nil {
			globalLogger.printf("NewRequest err: %v", err)
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req = req.WithContext(ctx)

		resp, err := selectedClient(downloadFundFlowURL).Do(req)
		tryNum++
		switch {
		case err != nil:
			if shouldRetry(err) {
				notifyAsync("downloadfundflow err: ", err)
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
				ErrCode    string `xml:"err_code"`
				ErrCodeDes string `xml:"err_code_des"`
			}
			if err = xml.Unmarshal(buf.Bytes(), &response); err == nil {
				switch response.ErrCode {
				case noBillExistErrForFundFlow.Error():
					return nil, billNoExistErr
				case systemerror.Error():
					notifyAsync("downloadfundflow err: ", err)
					time.Sleep(tempDelay)
					continue tryLoop
				default:
					return nil, errors.New(response.ErrCodeDes)
				}
			} else {
				break tryLoop
			}
		}
	}

	var response DownloadFundFlowResponse

	response.FundFlow = buf.String()
	response.FundFlow = strings.Replace(response.FundFlow, "`", "", -1)
	pos := strings.LastIndex(response.FundFlow, "资")

	firstPart := response.FundFlow[0:pos]
	secondPart := response.FundFlow[pos:]

	reader := csv.NewReader(strings.NewReader(firstPart))
	reader.LazyQuotes = true
	orderList, err := reader.ReadAll()
	if err != nil {
		globalLogger.printf("csv.NewReader err: %v", err)
		return nil, err
	}

	response.EntryList = make([]*FundFlowEntry, 0)

	if len(orderList) >= 2 {
		for _, order := range orderList[1:] {
			if len(order) == 11 {
				entry := new(FundFlowEntry)
				response.EntryList = append(response.EntryList, entry)

				entry.TimeEnd = order[0]
				entry.TransactionID = order[1]
				entry.FundFlowID = order[2]
				entry.BusinessName = order[3]
				entry.BusinessType = order[4]
				entry.TradeType = order[5]
				entry.TradeFee = order[6]
				entry.Balance = order[7]
				entry.Operator = order[8]
				entry.Remark = order[9]
				entry.BusinessID = order[10]
			}
		}
	}

	reader = csv.NewReader(strings.NewReader(secondPart))
	reader.LazyQuotes = true
	statisticsList, err := reader.ReadAll()
	if err != nil {
		globalLogger.printf("csv.NewReader err: %v", err)
		return nil, err
	}

	response.Statistics = make([]*FundFlowStatistics, 0)

	if len(statisticsList) >= 2 {
		for _, statistics := range statisticsList[1:] {
			if len(statistics) == 5 {
				s := new(FundFlowStatistics)
				response.Statistics = append(response.Statistics, s)

				s.TradeOrderCount = statistics[0]
				s.IncomeCount = statistics[1]
				s.TotalIncomeFee = statistics[2]

				s.SpendCount = statistics[3]
				s.TotalSpendFee = statistics[4]
			}
		}
	}

	return &response, nil

}

// FundFlowStatistics 账单统计
type FundFlowStatistics struct {
	TradeOrderCount string `csv:"资金流水总笔数"`
	IncomeCount     string `csv:"收入笔数"`
	TotalIncomeFee  string `csv:"收入金额"`
	SpendCount      string `csv:"支出笔数"`
	TotalSpendFee   string `csv:"支出金额"`
}

// FundFlowEntry 账单条目
type FundFlowEntry struct {
	TimeEnd       string `csv:"记账时间"`
	TransactionID string `csv:"微信支付业务单号"`
	FundFlowID    string `csv:"资金流水单号"`
	BusinessName  string `csv:"业务名称"`
	BusinessType  string `csv:"业务类型"`
	TradeType     string `csv:"收支类型"`
	TradeFee      string `csv:"收支金额（元）"`
	Balance       string `csv:"账户结余（元）"`
	Operator      string `csv:"资金变更提交申请人"`
	Remark        string `csv:"备注"`
	BusinessID    string `csv:"业务凭证号"`
}
