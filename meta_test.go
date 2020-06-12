package wxpay

import "testing"

func TestMeta_IsNotEnough(t *testing.T) {
	meta := &Meta{
		ErrCode: "NOTENOUGH",
	}

	if !meta.IsNotEnough() {
		t.Error("NOTENOUGH")
	}
}

func TestMeta_IsTradeOverDue(t *testing.T) {
	meta := &Meta{
		ErrCode: "TRADE_OVERDUE",
	}

	if !meta.IsTradeOverDue() {
		t.Error("TRADE_OVERDUE")
	}
}

func TestMeta_IsRefundNotExist(t *testing.T) {
	meta := &Meta{
		ErrCode: "REFUNDNOTEXIST",
	}

	if !meta.IsRefundNotExist() {
		t.Error("REFUNDNOTEXIST")
	}
}
