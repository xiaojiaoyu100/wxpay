package wxpay

import (
	"encoding/xml"
	"testing"
)

var paidVerifyBody = []byte(
	`<xml>
  <appid><![CDATA[48484848]]></appid>
  <attach><![CDATA[支付测试]]></attach>
  <device_info><![CDATA[013467007045764]]></device_info>
  <bank_type><![CDATA[CFT]]></bank_type>
  <fee_type><![CDATA[CNY]]></fee_type>
  <sign_type><![CDATA[MD5]]></sign_type>
  <is_subscribe><![CDATA[Y]]></is_subscribe>
  <mch_id><![CDATA[5151515]]></mch_id>
  <nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
  <openid><![CDATA[441515151]]></openid>
  <out_trade_no><![CDATA[1409811653]]></out_trade_no>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <sign><![CDATA[B552ED6B279343CB493C5DD0D78AB241]]></sign>
  <sub_mch_id><![CDATA[515151]]></sub_mch_id>
  <time_end><![CDATA[20140903131540]]></time_end>
  <total_fee>1</total_fee>
  <settlement_total_fee>1</settlement_total_fee>
  <cash_fee>100</cash_fee>
  <cash_fee_type>CNY</cash_fee_type>
  <coupon_fee>100</coupon_fee>
<coupon_fee_0><![CDATA[10]]></coupon_fee_0>
<coupon_count><![CDATA[1]]></coupon_count>
<coupon_type><![CDATA[CASH]]></coupon_type>
<coupon_id><![CDATA[10000]]></coupon_id>
  <trade_type><![CDATA[JSAPI]]></trade_type>
  <transaction_id><![CDATA[1004400740201409030005092168]]></transaction_id>
</xml>`,
)

func TestUnmarshalPaidNotifyRequest(t *testing.T) {
	request := new(PaidNotifyRequest)
	if err := xml.Unmarshal(paidVerifyBody, request); err != nil {
		t.Error(err)
	} else {
		t.Log(request)
	}
}
