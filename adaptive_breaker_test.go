package adaptive_breaker

import (
	"errors"
	"log"
	"testing"
	"time"
)

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
