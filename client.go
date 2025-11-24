package anamericano

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 3
	defaultRetryDelay = 1 * time.Second
	baseURL           = "https://accounts.ana.st"
)

// Client An-Americano 권한 API 클라이언트를 나타냅니다
type Client struct {
	httpClient *fasthttp.Client
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
	// MaxConnsPerHost 호스트당 최대 연결 수 (기본값: 512)
	MaxConnsPerHost int
	// MaxIdleConnDuration 유휴 연결 유지 시간 (기본값: 10초)
	MaxIdleConnDuration time.Duration
}

// Logger 로깅을 위한 인터페이스
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// Authenticator 다양한 인증 방법을 위한 인터페이스
type Authenticator interface {
	// Authenticate fasthttp 요청에 인증 헤더를 추가합니다
	AuthenticateFastHTTP(req *fasthttp.Request) error
}

// BearerTokenAuth Bearer 토큰 인증 (API 토큰용 - 레거시)
type BearerTokenAuth struct {
	Token string
}

// AuthenticateFastHTTP 요청에 Bearer 토큰을 추가합니다
func (b *BearerTokenAuth) AuthenticateFastHTTP(req *fasthttp.Request) error {
	if b.Token == "" {
		return fmt.Errorf("bearer token is empty")
	}
	req.Header.Set("Authorization", "Bearer "+b.Token)
	return nil
}

// OAuthTokenAuth OAuth 사용자 토큰 인증 (권장)
type OAuthTokenAuth struct {
	AccessToken string
}

// AuthenticateFastHTTP 요청에 OAuth 사용자 토큰을 추가합니다
func (o *OAuthTokenAuth) AuthenticateFastHTTP(req *fasthttp.Request) error {
	if o.AccessToken == "" {
		return fmt.Errorf("oauth access token is empty")
	}
	req.Header.Set("Authorization", "Bearer "+o.AccessToken)
	return nil
}

// TokenProvider 동적으로 토큰을 제공하는 인터페이스
type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
}

// DynamicTokenAuth 동적 토큰 제공자를 사용한 인증
type DynamicTokenAuth struct {
	Provider TokenProvider
	ctx      context.Context
}

// AuthenticateFastHTTP 요청 시점에 토큰을 가져와서 인증합니다
func (d *DynamicTokenAuth) AuthenticateFastHTTP(req *fasthttp.Request) error {
	ctx := d.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	token, err := d.Provider.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	if token == "" {
		return fmt.Errorf("token is empty")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

// ContextTokenAuth 컨텍스트에서 토큰을 가져오는 인증
type ContextTokenAuth struct {
	ctx context.Context
}

type contextKey string

const tokenContextKey contextKey = "anamericano_token"

// WithToken 컨텍스트에 토큰을 추가합니다
func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenContextKey, token)
}

// AuthenticateFastHTTP 컨텍스트에서 토큰을 가져와 인증합니다
func (c *ContextTokenAuth) AuthenticateFastHTTP(req *fasthttp.Request) error {
	ctx := c.ctx
	if ctx == nil {
		return fmt.Errorf("no context set")
	}
	token, ok := ctx.Value(tokenContextKey).(string)
	if !ok || token == "" {
		return fmt.Errorf("no token found in context")
	}
	req.Header.Set("Authorization", "Bearer "+token)
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
	return e.Status == fasthttp.StatusForbidden
}

// IsUnauthorized 인증되지 않은 오류인지 확인합니다
func (e *APIError) IsUnauthorized() bool {
	return e.Status == fasthttp.StatusUnauthorized
}

// IsNotFound 찾을 수 없음 오류인지 확인합니다
func (e *APIError) IsNotFound() bool {
	return e.Status == fasthttp.StatusNotFound
}

// IsBadRequest 잘못된 요청 오류인지 확인합니다
func (e *APIError) IsBadRequest() bool {
	return e.Status == fasthttp.StatusBadRequest
}

// NewClient 새로운 권한 API 클라이언트를 생성합니다
//
// 컨텍스트 기반 인증 (가장 효율적 - 클라이언트 1개 재사용):
//
//	client := anamericano.NewClient(&anamericano.ContextTokenAuth{}, nil)
//	ctx := anamericano.WithToken(context.Background(), userToken)
//	resp, err := client.CheckPermission(ctx, req)
//
// OAuth 토큰 사용 (고정 토큰):
//
//	auth := &anamericano.OAuthTokenAuth{AccessToken: "oauth_access_token"}
//	c := anamericano.NewClient(auth, nil)
//
// 동적 토큰 제공자 사용 (효율적):
//
//	auth := &anamericano.DynamicTokenAuth{Provider: yourTokenProvider}
//	c := anamericano.NewClient(auth, nil)
func NewClient(auth Authenticator, opts *ClientOptions) *Client {
	if opts == nil {
		opts = &ClientOptions{
			Timeout:             defaultTimeout,
			MaxRetries:          defaultMaxRetries,
			RetryDelay:          defaultRetryDelay,
			MaxConnsPerHost:     512,
			MaxIdleConnDuration: 10 * time.Second,
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
	if opts.MaxConnsPerHost == 0 {
		opts.MaxConnsPerHost = 512
	}
	if opts.MaxIdleConnDuration == 0 {
		opts.MaxIdleConnDuration = 10 * time.Second
	}

	return &Client{
		httpClient: &fasthttp.Client{
			ReadTimeout:                   opts.Timeout,
			WriteTimeout:                  opts.Timeout,
			MaxConnsPerHost:               opts.MaxConnsPerHost,
			MaxIdleConnDuration:           opts.MaxIdleConnDuration,
			DisableHeaderNamesNormalizing: false,
			DisablePathNormalizing:        false,
			NoDefaultUserAgentHeader:      false,
		},
		auth:    auth,
		options: opts,
	}
}

// doRequest 재시도 로직을 사용하여 HTTP 요청을 실행합니다
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := baseURL + path
	var lastErr error

	// 컨텍스트를 authenticator에 전달 (ContextTokenAuth용)
	if ctxAuth, ok := c.auth.(*ContextTokenAuth); ok {
		ctxAuth.ctx = ctx
	}
	if dynAuth, ok := c.auth.(*DynamicTokenAuth); ok {
		dynAuth.ctx = ctx
	}

	// 요청 본문을 한 번만 마샬링하여 재시도 시 재사용 (메모리 할당 최적화)
	var jsonData []byte
	var marshalErr error
	if body != nil {
		jsonData, marshalErr = json.Marshal(body)
		if marshalErr != nil {
			return fmt.Errorf("failed to marshal request body: %w", marshalErr)
		}
	}

	for attempt := 0; attempt <= c.options.MaxRetries; attempt++ {
		if attempt > 0 {
			backoffDelay := c.options.RetryDelay * time.Duration(attempt)

			timer := time.NewTimer(backoffDelay)
			select {
			case <-ctx.Done():
				// 타이머 중지 및 채널 드레인
				if !timer.Stop() {
					<-timer.C
				}
				return ctx.Err()
			case <-timer.C:
				// 타이머 정상 만료 - 이미 채널에서 값을 읽었으므로 정리 불필요
			}

			if c.options.Logger != nil {
				c.options.Logger.Debug("retrying request", "attempt", attempt, "url", url)
			}
		}

		// 익명 함수로 스코프 생성하여 즉시 릴리즈
		err := func() error {
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseRequest(req)
			defer fasthttp.ReleaseResponse(resp)

			req.SetRequestURI(url)
			req.Header.SetMethod(method)

			// 요청 본문 설정 (이미 마샬링된 데이터 사용)
			if len(jsonData) > 0 {
				req.SetBody(jsonData)
				req.Header.SetContentType("application/json")
			}

			// 인증 헤더 추가
			if c.auth != nil {
				if err := c.auth.AuthenticateFastHTTP(req); err != nil {
					return fmt.Errorf("authentication failed: %w", err)
				}
			}

			// 타임아웃이 있는 요청 실행
			if err := c.httpClient.DoTimeout(req, resp, c.options.Timeout); err != nil {
				return fmt.Errorf("request failed: %w", err)
			}

			statusCode := resp.StatusCode()
			// Body()는 내부 버퍼를 반환하므로 한 번만 호출하고 재사용
			bodyBytes := resp.Body()

			// 성공 응답 처리
			if statusCode >= 200 && statusCode < 300 {
				if result != nil && len(bodyBytes) > 0 {
					// bodyBytes를 직접 사용 (복사 방지)
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
				// string() 변환은 복사를 일으키지만 에러 케이스이므로 허용
				return fmt.Errorf("HTTP %d: %s", statusCode, string(bodyBytes))
			}

			// 클라이언트 오류(4xx)는 재시도하지 않음 (429 제외)
			if statusCode >= 400 && statusCode < 500 && statusCode != 429 {
				if c.options.Logger != nil {
					c.options.Logger.Error("client error", "status", statusCode, "message", apiErr.Message)
				}
				return &apiErr
			}

			// 서버 오류(5xx)와 429는 재시도
			if c.options.Logger != nil {
				c.options.Logger.Error("server error, will retry", "status", statusCode, "message", apiErr.Message)
			}
			return &apiErr
		}()

		if err == nil {
			return nil
		}

		// APIError인 경우 재시도 여부 판단
		if apiErr, ok := err.(*APIError); ok {
			// 4xx (429 제외)는 즉시 반환
			if apiErr.Status >= 400 && apiErr.Status < 500 && apiErr.Status != 429 {
				return apiErr
			}
			// 5xx와 429는 재시도
			lastErr = apiErr
		} else {
			// 네트워크 오류 등은 재시도
			lastErr = err
			if c.options.Logger != nil {
				c.options.Logger.Error("request error", "error", err, "attempt", attempt)
			}
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// SetAuth 클라이언트의 인증 방법을 업데이트합니다
func (c *Client) SetAuth(auth Authenticator) {
	c.auth = auth
}
