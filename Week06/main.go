package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type metrics struct {
	success int32
	fail    int32
}

type bucket struct {
	data        metrics
	windowStart int64
}

type RollingNumber struct {
	buckets []*bucket
	size    int64
	width   int64
	tail    int64
	mux     sync.RWMutex
}

func NewRollingNumber(size, width int64) *RollingNumber {
	return &RollingNumber{
		size:    size,
		width:   width,
		buckets: make([]*bucket, size),
		tail:    0,
	}
}

func (rp *RollingNumber) getCurrent() *bucket {
	rp.mux.Lock()
	defer rp.mux.Unlock()

	current := time.Now().Unix()
	if rp.tail == 0 && rp.buckets[rp.tail] == nil {
		bk := &bucket{
			data:        metrics{},
			windowStart: current,
		}
		rp.buckets[rp.tail] = bk
		return bk
	}

	last := rp.buckets[rp.tail]
	if current < last.windowStart+rp.width {
		return last
	}

	for i := 0; i < int(rp.size); i++ {
		last := rp.buckets[rp.tail]
		if current < last.windowStart+rp.width {
			return last
		} else if current-(last.windowStart+rp.width) > rp.size*rp.width {
			rp.tail = 0
			rp.buckets = make([]*bucket, rp.size)
			rp.mux.Unlock()
			return rp.getCurrent()
		} else {
			rp.tail++
			bk := &bucket{
				data:        metrics{},
				windowStart: last.windowStart + rp.width,
			}

			if rp.tail >= rp.size {
				copy(rp.buckets[:], rp.buckets[1:])
				rp.tail--
			}
			rp.buckets[rp.tail] = bk
		}
	}

	return rp.buckets[rp.tail]
}

func (rp *RollingNumber) incrSuccess() {
	bk := rp.getCurrent()
	atomic.AddInt32(&bk.data.success, 1)
}

func (rp *RollingNumber) incrFail() {
	bk := rp.getCurrent()
	atomic.AddInt32(&bk.data.fail, 1)
}

func (rp *RollingNumber) getSum() metrics {
	m := metrics{}
	rp.mux.RLock()
	defer rp.mux.RUnlock()
	for _, v := range rp.buckets {
		m.success += v.data.success
		m.fail += v.data.fail
	}
	return m
}

func main() {
	// 滑动窗口计数器
	// 用于给熔断器提供数据依据
	rw := NewRollingNumber(2, 1)
	fmt.Println(time.Now().Unix())
	if false {
		//test 1
		rw.incrSuccess()
		time.Sleep(time.Second * 1)
		rw.incrSuccess()
		time.Sleep(time.Second * 1)
		rw.incrSuccess()
		time.Sleep(time.Second * 1)
		fmt.Printf("%+v,%+v\n", rw.buckets[0], rw.buckets[1])
	}

	{
		//test 2
		rw.incrSuccess()
		time.Sleep(time.Second * 3)
		rw.incrSuccess()
		fmt.Printf("%+v,%+v\n", rw.buckets[0], rw.buckets[1])
	}
}
