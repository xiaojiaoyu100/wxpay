package wxpay

import (
	"encoding/xml"
	"testing"
)

var orderQueryRespBytes = []byte(
	`<xml>
   <return_code><![CDATA[SUCCESS]]></return_code>
   <return_msg><![CDATA[OK]]></return_msg>
   <appid><![CDATA[55151515]]></appid>
   <mch_id><![CDATA[15151515]]></mch_id>
   <device_info><![CDATA[1000]]></device_info>
   <nonce_str><![CDATA[TN55wO9Pba5yENl8]]></nonce_str>
   <sign><![CDATA[BDF0099C15FF7BC6B1585FBB110AB635]]></sign>
   <result_code><![CDATA[SUCCESS]]></result_code>
   <openid><![CDATA[oUpF8uN95-Ptaags6E_roPHg7AG0]]></openid>
   <is_subscribe><![CDATA[Y]]></is_subscribe>
   <trade_type><![CDATA[MICROPAY]]></trade_type>
   <bank_type><![CDATA[CCB_DEBIT]]></bank_type>
   <total_fee>1</total_fee>
   <cash_fee>100</cash_fee>
   <cash_fee_type>CNY</cash_fee_type>
   <settlement_total_fee>100</settlement_total_fee>
   <coupon_fee>100</coupon_fee>
   <coupon_count>1</coupon_count>
   <fee_type><![CDATA[CNY]]></fee_type>
   <transaction_id><![CDATA[1008450740201411110005820873]]></transaction_id>
   <out_trade_no><![CDATA[1415757673]]></out_trade_no>
   <attach><![CDATA[订单额外描述]]></attach>
   <time_end><![CDATA[20141111170043]]></time_end>
   <trade_state><![CDATA[SUCCESS]]></trade_state>
   <trade_state_desc><![CDATA[支付失败，请重新下单支付]]></trade_state_desc>
</xml>`,
)

func TestUnmarshalOrderQueryResponse(t *testing.T) {
	resp := new(OrderQueryResponse)
	if err := xml.Unmarshal(orderQueryRespBytes, resp); err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
