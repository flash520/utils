/**
 * @Author: koulei
 * @Description:
 * @File: main_test
 * @Version: 1.0.0
 * @Date: 2021/8/8 21:45
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRateLimiter(t *testing.T) {

	failure := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(i int) {
			get, err := http.Get("http://localhost/")
			if err != nil {
				panic(err)
			}
			body, err := ioutil.ReadAll(get.Body)
			if err != nil {
				panic(err)
			}
			fmt.Printf("request: %d, response: %s\n", i, string(body))
		}(i)
	}
	<-failure
	fmt.Println("Testing Success")
}
