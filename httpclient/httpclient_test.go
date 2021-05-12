/**
 * @Author: koulei
 * @Description: TODO
 * @File:  client_test
 * @Version: 1.0.0
 * @Date: 2021/5/12 00:30
 */

package httpclient

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestHttpClient_GetConn(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	c := CreateHttpClient(2, 30, 15).GetConn()
	res, err := c.Get("https://sssxx.com/")
	if err != nil {
		log.Error(err.Error())
	}

	defer res.Body.Close()

	log.Info(res.StatusCode)

}
