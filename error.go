package wxpay

// Error error
type Error string

// 自定义错误
const (
	signNotMatchErr = Error("SignNotMatch")
)

// 微信的错误,请不要修改内容
const (
	billNoExistErr = Error("No Bill Exist")
)

func (err Error) Error() string {
	return string(err)
}

// IsBillNoExist 账单不存在
func IsBillNoExist(err error) bool {
	return err == billNoExistErr
}

func shouldRetry(err error) bool {
	switch err := err.(type) {
	case interface {
		Temporary() bool
	}:
		return err.Temporary()
	case interface {
		Timeout() bool
	}:
		return err.Timeout()
	default:
		return false
	}
}
