/**
 * @Author: koulei
 * @Description: TODO
 * @File:  httpClient
 * @Version: 1.0.0
 * @Date: 2021/5/12 00:27
 */

package httpclient

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var httpclient *http.Client

type httpClient struct {
}

// CreateHttpClient 初始化自定义http 客户端
func CreateHttpClient(dialTimeout, keepAliveTimeout, responseTimeout time.Duration) *httpClient {
	var once sync.Once
	once.Do(func() {
		httpclient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   dialTimeout * time.Second,
					KeepAlive: keepAliveTimeout * time.Second,
				}).DialContext,
				MaxIdleConns:          10,
				MaxIdleConnsPerHost:   10,
				MaxConnsPerHost:       8,
				IdleConnTimeout:       90 * time.Second,
				ResponseHeaderTimeout: responseTimeout * time.Second,
				DisableKeepAlives:     false,
				DisableCompression:    true,
			},
		}
	})

	return &httpClient{}
}

// GetConn 获取http客户端连接
func (h *httpClient) GetConn() *http.Client {
	return httpclient
}
