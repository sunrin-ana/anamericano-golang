package anamericano

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestDefaultLogger_Info(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := &DefaultLogger{}
	logger.Info("test message", "key", "value")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "[INFO]") {
		t.Error("expected [INFO] prefix")
	}
	if !strings.Contains(output, "test message") {
		t.Error("expected test message in output")
	}
}

func TestDefaultLogger_Error(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := &DefaultLogger{}
	logger.Error("error message", "error", "some error")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "[ERROR]") {
		t.Error("expected [ERROR] prefix")
	}
	if !strings.Contains(output, "error message") {
		t.Error("expected error message in output")
	}
}

func TestDefaultLogger_Debug(t *testing.T) {
	logger := &DefaultLogger{}
	logger.Debug("debug message")
}

func TestNoOpLogger_Info(t *testing.T) {
	logger := &NoOpLogger{}
	logger.Info("test")
}

func TestNoOpLogger_Error(t *testing.T) {
	logger := &NoOpLogger{}
	logger.Error("test")
}

func TestNoOpLogger_Debug(t *testing.T) {
	logger := &NoOpLogger{}
	logger.Debug("test")
}
