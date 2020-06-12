package wxpay

import (
	"crypto/rand"
	"fmt"
)

func nonceStr() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		notifyAsync("nonceStr err: ", err.Error())
		return ""
	}
	return fmt.Sprintf("%x", b)
}
