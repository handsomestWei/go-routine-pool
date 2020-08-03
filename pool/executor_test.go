package pool

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// atomic add
func TestNewPool_1(t *testing.T) {
	pool := NewPool(2)
	var n int32 = 0
	for i := 0; i < 10; i++ {
		n := atomic.AddInt32(&n, 1)
		pool.Submit(func() {
			fmt.Println(fmt.Sprintf("n=%d", n))
		})
	}
	fmt.Println("end")
	fmt.Println(pool.IsShutDown())

	time.Sleep(4 * time.Second)
}

// closure add
func TestNewPool_2(t *testing.T) {
	pool := NewPool(2)
	var n int32 = 0
	for i := 0; i < 10; i++ {
		n++
		pool.Submit(func() {
			fmt.Println(fmt.Sprintf("n=%d", n))
		})
	}
	pool.ShutDown()
	fmt.Println(pool.IsShutDown())
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	pool.Submit(func() {
		fmt.Println(fmt.Sprintf("n=%d", n))
	})

	time.Sleep(4 * time.Second)
}

func TestNewWaitPool_1(t *testing.T) {
	pool := NewWaitPool(2)
	var n int32 = 0
	for i := 0; i < 10; i++ {
		n := atomic.AddInt32(&n, 1)
		pool.Submit(func() {
			fmt.Println(fmt.Sprintf("n=%d", n))
		})
	}
	pool.Wait()
	fmt.Println("end")
	fmt.Println(pool.IsShutDown())

	time.Sleep(4 * time.Second)
}
