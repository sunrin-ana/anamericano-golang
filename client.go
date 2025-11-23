package anamericano

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 3
	defaultRetryDelay = 1 * time.Second
	baseURL           = "https://accounts.ana.st"
)

// Client An-Americano 권한 API 클라이언트를 나타냅니다
type Client struct {
	httpClient *http.Client
	auth       Authenticator
	options    *ClientOptions
}

// ClientOptions 클라이언트 설정 옵션을 포함합니다
type ClientOptions struct {
	// Timeout HTTP 요청 타임아웃 (기본값: 30초)
	Timeout time.Duration
	// MaxRetries 실패한 요청에 대한 최대 재시도 횟수 (기본값: 3)
	MaxRetries int
	// RetryDelay 재시도 간 지연 시간 (기본값: 1초)
	RetryDelay time.Duration
	// Logger 디버그 및 에러 로깅을 위한 로거
	Logger Logger
}

// Logger 로깅을 위한 인터페이스
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// Authenticator 다양한 인증 방법을 위한 인터페이스
type Authenticator interface {
	Authenticate(req *http.Request) error
}

// BearerTokenAuth Bearer 토큰 인증
type BearerTokenAuth struct {
	Token string
}

// Authenticate 요청에 Bearer 토큰을 추가합니다
func (b *BearerTokenAuth) Authenticate(req *http.Request) error {
	if b.Token == "" {
		return fmt.Errorf("bearer token is empty")
	}
	req.Header.Set("Authorization", "Bearer "+b.Token)
	return nil
}

// APIError API 오류 응답을 나타냅니다
type APIError struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	ErrorType string `json:"error"`
	Message   string `json:"message"`
	Path      string `json:"path"`
}

// Error 오류 메시지를 반환합니다
func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s - %s (path: %s)", e.Status, e.ErrorType, e.Message, e.Path)
}

// IsPermissionDenied 권한 거부 오류인지 확인합니다
func (e *APIError) IsPermissionDenied() bool {
	return e.Status == http.StatusForbidden
}

// IsUnauthorized 인증되지 않은 오류인지 확인합니다
func (e *APIError) IsUnauthorized() bool {
	return e.Status == http.StatusUnauthorized
}

// IsNotFound 찾을 수 없음 오류인지 확인합니다
func (e *APIError) IsNotFound() bool {
	return e.Status == http.StatusNotFound
}

// IsBadRequest 잘못된 요청 오류인지 확인합니다
func (e *APIError) IsBadRequest() bool {
	return e.Status == http.StatusBadRequest
}

// NewClient 새로운 권한 API 클라이언트를 생성합니다
//
// 예시:
//
//	auth := &client.BearerTokenAuth{Token: "토큰"}
//	c := client.NewClient(auth, nil)
func NewClient(auth Authenticator, opts *ClientOptions) *Client {
	if opts == nil {
		opts = &ClientOptions{
			Timeout:    defaultTimeout,
			MaxRetries: defaultMaxRetries,
			RetryDelay: defaultRetryDelay,
		}
	}

	if opts.Timeout == 0 {
		opts.Timeout = defaultTimeout
	}
	if opts.MaxRetries == 0 {
		opts.MaxRetries = defaultMaxRetries
	}
	if opts.RetryDelay == 0 {
		opts.RetryDelay = defaultRetryDelay
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
		auth:    auth,
		options: opts,
	}
}

// doRequest 재시도 로직을 사용하여 HTTP 요청을 실행합니다
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := baseURL + path
	var lastErr error

	for attempt := 0; attempt <= c.options.MaxRetries; attempt++ {
		if attempt > 0 {
			backoffDelay := c.options.RetryDelay * time.Duration(attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoffDelay):
			}

			if c.options.Logger != nil {
				c.options.Logger.Debug("retrying request", "attempt", attempt, "url", url)
			}

			// 재시도를 위해 요청 본문 재생성
			if body != nil {
				jsonData, _ := json.Marshal(body)
				reqBody = bytes.NewBuffer(jsonData)
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		if c.auth != nil {
			if err := c.auth.Authenticate(req); err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			if c.options.Logger != nil {
				c.options.Logger.Error("request error", "error", err, "attempt", attempt)
			}
			continue
		}

		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if result != nil && len(bodyBytes) > 0 {
				if err := json.Unmarshal(bodyBytes, result); err != nil {
					return fmt.Errorf("failed to unmarshal response: %w", err)
				}
			}
			return nil
		}

		// 오류 응답 처리
		var apiErr APIError
		if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
			// 오류를 Parsing할 수 없으면 일반 오류 반환
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
		}

		// 클라이언트 오류(4xx)는 재시도하지 않음 (429 제외)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
			if c.options.Logger != nil {
				c.options.Logger.Error("client error", "status", resp.StatusCode, "message", apiErr.Message)
			}
			return &apiErr
		}

		// 서버 오류(5xx)와 429는 재시도
		lastErr = &apiErr
		if c.options.Logger != nil {
			c.options.Logger.Error("server error, will retry", "status", resp.StatusCode, "message", apiErr.Message)
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// SetAuth 클라이언트의 인증 방법을 업데이트합니다
func (c *Client) SetAuth(auth Authenticator) {
	c.auth = auth
}
