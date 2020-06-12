package wxpay

import (
	"encoding/xml"
	"testing"
)

var refundResponseBytes = []byte(
	`<xml>
   <return_code><![CDATA[SUCCESS]]></return_code>
   <return_msg><![CDATA[OK]]></return_msg>
   <appid><![CDATA[454555888]]></appid>
   <mch_id><![CDATA[448848]]></mch_id>
   <nonce_str><![CDATA[NfsMFbUFpdbEhPXP]]></nonce_str>
   <sign><![CDATA[B7274EB9F8925EB93100DD2085FA56C0]]></sign>
   <result_code><![CDATA[SUCCESS]]></result_code>
   <transaction_id><![CDATA[1008450740201411110005820873]]></transaction_id>
   <out_trade_no><![CDATA[1415757673]]></out_trade_no>
   <out_refund_no><![CDATA[1415701182]]></out_refund_no>
   <refund_id><![CDATA[2008450740201411110000174436]]></refund_id>
   <refund_fee>1</refund_fee>
   <settlement_refund_fee>100</settlement_refund_fee>
   <total_fee>100</total_fee>
   <settlement_total_fee>100</settlement_total_fee>
   <fee_type>CNY</fee_type>
   <cash_fee>100</cash_fee>
   <cash_fee_type>CNY</cash_fee_type>
   <cash_refund_fee>100</cash_refund_fee>
   <coupon_refund_fee>100</coupon_refund_fee>
   <coupon_refund_count>100</coupon_refund_count>
</xml>`,
)

func TestUnmarshalRefundResponse(t *testing.T) {
	response := new(RefundResponse)
	if err := xml.Unmarshal(refundResponseBytes, response); err != nil {
		t.Error(err)
	} else {
		t.Log(response)
	}
}
