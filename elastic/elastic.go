/**
 * @Author: koulei
 * @Description: TODO
 * @File:  elastic
 * @Version: 1.0.0
 * @Date: 2021/5/12 18:19
 */

package elastic

import (
	"gitee.com/flash520/utils/httpclient"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"sync"
)

var (
	client     *elastic.Client
	HttpClient = httpclient.CreateHttpClient(5, 30, 15).GetConn()
)

type Elastic struct {
}

func CreateElastic(addr, username, password string) *Elastic {
	var once sync.Once
	once.Do(func() {
		c, err := elastic.NewClient(
			elastic.SetURL("http://"+addr),
			elastic.SetSniff(false),
			elastic.SetBasicAuth(username, password),
			elastic.SetHttpClient(HttpClient),
			elastic.SetErrorLog(log.New(os.Stderr, "[ELASTIC] - ERR - ", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "[ELASTIC] - INFO - ", log.LstdFlags)),
		)
		if err != nil {
			panic(err)
		}
		client = c
	})
	return &Elastic{}
}

func (e *Elastic) GetConn() *elastic.Client {
	return client
}
