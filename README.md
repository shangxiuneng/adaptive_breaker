# 基于golang实现的自适应熔断器

参考 [Handling Overload](https://sre.google/sre-book/handling-overload/)

## 业务背景
在我的业务开发中，遇到了以下场景：  
大部分情况下需要同步的调用下游业务，保证实时性，当下游业务不可用时，将消息发送到kafka中，
由下游自行订阅kafka。于是什么时候同步调用转异步调用，异步调用何时恢复为同步调用变成了问题，在此背景下实现了本仓库。

## 实现说明

TODO


## 如何使用

```go
func TestNewAdaptiveBreaker(t *testing.T) {
	breaker := NewAdaptiveBreaker(0.8, 10, 5*time.Second)

	for i := 0; i < 100; i++ {
		err := breaker.Execute(service)
		if err != nil {
			if errors.Is(err, ErrCircuitOpen) {
				log.Println("Circuit open: request blocked.")
			} else {
				log.Println("Request failed:", err)
			}
		} else {
			log.Println("Request succeeded.")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// 业务逻辑
func service() error {
	return nil
}

```

1. 初始化函数:  
```go
NewAdaptiveBreaker(threshold float64, minRequests int, coolDownPeriod time.Duration)
```
`threshold` 触发熔断的失败阈值  
`minRequests` 最小请求数量  
`coolDownPeriod` 冷却时间(或者可以理解为评估时间)    


2. 执行业务逻辑:  
```go
breaker.Execute(service)
```
`service`是我们的业务逻辑函数。


## 如果需要二次开发
1. 修改`Execute`函数的实现，当业务发生错误时，可以自定义根据那些指标进行上报。
2. 增加数据的持久化，比如把数据存储到redis中，因为目前的实现是单实例的，多个节点时间没法共享下游集群的状态。