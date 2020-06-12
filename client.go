package wxpay

// Client 微信支付
type Client struct {
	apiKey string
	mchID  string
}

// New 产生新的微信支付client
func New(apiKey, mchID string) *Client {
	return &Client{
		apiKey: apiKey,
		mchID:  mchID,
	}
}
