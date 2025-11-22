package anamericano

import "fmt"

// DefaultLogger 기본 로거
type DefaultLogger struct{}

func (d *DefaultLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
}

func (d *DefaultLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] "+msg+"\n", args...)
}

func (d *DefaultLogger) Debug(msg string, args ...interface{}) {
	// TODO: 디버그 구현하기 차기 부장이 해주겠지 ^^
}

// NoOpLogger 아무것도 하지 않는 로거
type NoOpLogger struct{}

func (n *NoOpLogger) Info(msg string, args ...interface{})  {}
func (n *NoOpLogger) Error(msg string, args ...interface{}) {}
func (n *NoOpLogger) Debug(msg string, args ...interface{}) {}
