package adaptive_breaker

import (
	"errors"
	"log"
	"sync"
	"time"
)

var (
	ErrCircuitOpen = errors.New("circuit is open")
)

// AdaptiveBreaker 熔断器
type AdaptiveBreaker struct {
	sync.RWMutex
	requests       int
	successes      int
	failures       int
	lastCheck      time.Time
	threshold      float64
	minRequests    int
	coolDownPeriod time.Duration
	state          string
}

// NewAdaptiveBreaker 初始化熔断器
func NewAdaptiveBreaker(threshold float64, minRequests int, coolDownPeriod time.Duration) *AdaptiveBreaker {
	return &AdaptiveBreaker{
		threshold:      threshold,      // 熔断器的失败率阈值
		minRequests:    minRequests,    // 开始检测失败率的最低请求数
		coolDownPeriod: coolDownPeriod, // 评估时间
		state:          "closed",       // 熔断器的状态
		lastCheck:      time.Now(),     // 上一次评估的时间
	}
}

func (b *AdaptiveBreaker) Allow() error {
	b.Lock()
	defer b.Unlock()

	now := time.Now()

	if b.state == "open" && now.Sub(b.lastCheck) >= b.coolDownPeriod {
		b.state = "half-open"
	}

	if b.state == "open" {
		// 如果熔断器是打开状态的 则返回ErrCircuitOpen错误 并且所有的请求都会被拦截
		return ErrCircuitOpen
	}

	return nil
}

func (b *AdaptiveBreaker) Report(success bool) {
	b.Lock()
	defer b.Unlock()

	b.requests++
	if success {
		b.successes++
	} else {
		b.failures++
	}

	now := time.Now()

	if b.requests >= b.minRequests && now.Sub(b.lastCheck) >= b.coolDownPeriod {
		b.evaluate()
		b.lastCheck = now
	}
}

func (b *AdaptiveBreaker) evaluate() {
	successRate := float64(b.successes) / float64(b.requests)
	if b.state == "half-open" {
		if successRate >= b.threshold {
			b.state = "closed"
			log.Println("Breaker closed due to successful recovery.")
		} else {
			b.state = "open"
			log.Println("Breaker reopened due to low success rate.")
		}
	} else if successRate < b.threshold {
		b.state = "open"
		log.Println("Breaker opened due to low success rate.")
	}

	b.requests = 0
	b.successes = 0
	b.failures = 0
}

func (b *AdaptiveBreaker) Execute(fn func() error) error {
	if err := b.Allow(); err != nil {
		return err
	}

	// 执行业务逻辑并进行错误处理
	err := fn()
	b.Report(err == nil)
	return err
}
