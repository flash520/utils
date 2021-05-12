/**
 * @Author: koulei
 * @Description: TODO
 * @File:  minio
 * @Version: 1.0.0
 * @Date: 2021/5/12 18:27
 */

package minio

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

var (
	transport = func() *http.Transport {
		tr := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 90 * time.Second,
			}).DialContext,
			MaxIdleConns:          256,
			MaxIdleConnsPerHost:   256,
			MaxConnsPerHost:       8,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			DisableKeepAlives:     false,
			DisableCompression:    true,
		}
		return tr
	}

	client *minio.Client
	err    error
)

type Minio struct {
}

func CreateMinio(addr, accessKey, secretKey string) *Minio {
	var once sync.Once
	once.Do(func() {
		client, err = minio.New(addr, &minio.Options{
			Creds:     credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure:    false,
			Transport: transport(),
		})
		if err != nil {
			log.Println(err.Error())
		}
	})

	return &Minio{}
}

func (m *Minio) GetConn() *minio.Client {
	return client
}
