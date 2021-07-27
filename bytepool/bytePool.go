/**
 * @Author: koulei
 * @Description:
 * @File: bytePool
 * @Version: 1.0.0
 * @Date: 2021/7/27 12:12
 */

package bytepool

// bytePool 是一个基于 channel 构造的 []byte 池
// 旨在减少大量 IO 操作中，重复定义 buffer 造成的内存分配延迟
// 提供 GET 方法获取 []byte, PUT 方法返还 []byte 放入 channel 以重复使用
type bytePool struct {
	c    chan []byte
	w    int
	wCap int
}

func NewBytePool(maxSize, w, cap int) *bytePool {
	return &bytePool{
		c:    make(chan []byte, maxSize),
		w:    w,
		wCap: cap,
	}
}

func (bp *bytePool) GET() (b []byte) {
	select {
	case b = <-bp.c:
	default:
		if bp.wCap > 0 {
			b = make([]byte, bp.w, bp.wCap)
		} else {
			b = make([]byte, bp.w)
		}
	}
	return
}

func (bp *bytePool) PUT(b []byte) {
	select {
	case bp.c <- b:
	default:
	}
}
