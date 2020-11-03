package wxpay

import (
	"crypto/tls"
	"encoding/pem"
	"net/http"
	"time"

	"golang.org/x/crypto/pkcs12"
)

var client = &http.Client{
	Timeout: 60 * time.Second,
}

var tlsClient = &http.Client{
	Timeout: 60 * time.Second,
}

func selectedClient(url string) *http.Client {
	switch url {
	case refundURL, reverseURL, transferURL, transferInfoURL, downloadFundFlowURL:
		return tlsClient
	default:
		return client
	}
}

// SetTLSClient tls client
func SetTLSClient(pfxData []byte, password string) error {
	blocks, err := pkcs12.ToPEM(pfxData, password)
	if err != nil {
		globalLogger.printf("ToPEM err: %v", err)
		return err
	}
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		globalLogger.printf("X509KeyPair err: %v", err)
		return err
	}
	tlsClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	return nil
}