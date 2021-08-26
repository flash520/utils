/**
 * @Author: koulei
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2021/8/16 12:54
 */

package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
)

var (
	x  int64
	l  sync.Mutex
	wg sync.WaitGroup
)

func add() {
	x++
	wg.Done()
}
func ladd() {
	l.Lock()
	x++
	l.Unlock()
	wg.Done()
}
func atomicadd() {
	atomic.AddInt64(&x, 1)
	wg.Done()
}

func main() {
	v4, _ := uuid.NewV4()
	s := v4.String()
	fmt.Println("X-Request-ID: ", strings.ToUpper(s))
	start := time.Now()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int64) {
			atomicadd()
		}(int64(i))
	}
	wg.Wait()
	fmt.Println(x)
	fmt.Println(time.Since(start))
}
