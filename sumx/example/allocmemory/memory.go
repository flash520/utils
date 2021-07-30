/**
 * @Author: koulei
 * @Description:
 * @File: memory
 * @Version: 1.0.0
 * @Date: 2021/7/22 19:53
 */

package allocmemory

import (
	"io"
	"os"
)

func prefixAllocMemory(buf chan []byte) {
	var n int
	var err error
	f, err := os.Open("/Users/koulei/msg.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// info, _ := f.Stat()
	var result = make([]int, 0, 1000)
	for {
		n, err = f.Read(<-buf)
		if err != nil || err == io.EOF {
			break
		}
		result = append(result, n)
	}
}

func varAllocMemory() {
	var d, i int
	for {
		d += 1 + i
		if d == 10 {
			break
		}
	}
}
