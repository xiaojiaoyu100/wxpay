package wxpay

import (
	"bytes"
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

func (c *Client) request(url string, in interface{}, out interface{}) ([]byte, error) {
	const (
		max = 1000 * time.Millisecond
	)
	var (
		tempDelay time.Duration
		tryNum    = 0
		err       error
		body      []byte
	)
	tryLoop:
	for {
		if tempDelay == 0 {
			tempDelay = 100 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if tempDelay > max {
			tempDelay = max
		}
		body, err = c.doRequest(url, in, out)
		tryNum++
		if tryNum > 3 {
			return body, err
		}
		switch {
		case shouldRetry(err):
			notifyAsync("doRequest err: ", err)
			time.Sleep(tempDelay)
			continue tryLoop
		default:
			for i := true; i; i = false {
				if out == nil {
					break
				}
				ot := reflect.TypeOf(out)
				if ot.Kind() != reflect.Ptr && ot.Kind() != reflect.Interface {
					break
				}
				value := reflect.ValueOf(out).Elem()
				if reflect.TypeOf(value).Kind() != reflect.Struct {
					break
				}
				metaItf := value.FieldByName("Meta").Interface()
				if m, ok := metaItf.(Meta); ok {
					if m.IsSystemErr() ||
						m.IsBizerrNeedRetry() {
						notifyAsync("doRequest err: ", m.ErrCode)
						time.Sleep(tempDelay)
						continue tryLoop
					}
				}
			}
			return body, err
		}
	}
}

func (c *Client) doRequest(url string, in interface{}, out interface{}) ([]byte, error) {
	body, err := xml.Marshal(in)
	if err != nil {
		globalLogger.printf("%s xml marshal err: %s", url, err.Error())
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		globalLogger.printf("%s new request err: %s", url, err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	globalLogger.printf("%s %s %s", req.Method, req.URL.String(), string(body))

	resp, err := selectedClient(url).Do(req)
	if err != nil {
		globalLogger.printf("%s %s do err: %s", req.Method, req.URL.String(), err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		globalLogger.printf("%s %s read resp body err: %s", req.Method, req.URL.String(), err.Error())
		return nil, err
	}
	globalLogger.printf("%s %s %s", req.Method, req.URL.String(), string(body))

	if err := xml.Unmarshal(body, &out); err != nil {
		globalLogger.printf("unmarshal body err: %s, body: %s", err.Error(), string(body))
		return nil, err
	}

	if err := checkSign(body, c.apiKey); err != nil {
		globalLogger.printf("checkSign err: %s", err.Error())
		return nil, err
	}

	return body, nil
}
