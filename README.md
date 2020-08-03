# go-routine-pool
golang routine协程池，类似java的Executors

# Usage

使用`pool.NewPool`创建不等待的协程池

```
	pool := pool.NewPool(2)
	var n int32 = 0
	for i := 0; i < 10; i++ {
		n := atomic.AddInt32(&n, 1)
		pool.Submit(func() {
			fmt.Println(fmt.Sprintf("n=%d", n))
		})
	}
    fmt.Println("end")
```

使用`pool.NewWaitPool`创建等待的协程池

```
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
```