# go-routine-pool
golang routine协程池，包含一个简单的有限状态机。类似java的Executors

# Usage

使用`pool.NewPool`创建不等待的协程池

```
// 初始化协程池，指定协程数
pool := pool.NewPool(2)
var n int32 = 0
for i := 0; i < 10; i++ {
	n := atomic.AddInt32(&n, 1)
	pool.Submit(func() {
		fmt.Println(fmt.Sprintf("n=%d", n))
	})
}
// 不用等待所有协程处理结束。直接往下处理
fmt.Println("end")
```

使用`pool.NewWaitPool`创建等待的协程池

```
// 初始化协程池，指定协程数
pool := pool.NewWaitPool(2)
var n int32 = 0
for i := 0; i < 10; i++ {
	n := atomic.AddInt32(&n, 1)
	pool.Submit(func() {
		fmt.Println(fmt.Sprintf("n=%d", n))
	})
}
// 等待
pool.Wait()
// 等所有协程处理结束，再继续往下处理
fmt.Println("end")
```
