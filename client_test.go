package anamericano

import (
	"context"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

func TestNewClient(t *testing.T) {
	auth := &BearerTokenAuth{Token: "test-token"}
	client := NewClient(auth, nil)

	if client == nil {
		t.Fatal("expected client to be non-nil")
	}
	if client.auth != auth {
		t.Error("auth not set correctly")
	}
	if client.options.Timeout != defaultTimeout {
		t.Errorf("expected timeout %v, got %v", defaultTimeout, client.options.Timeout)
	}
	if client.options.MaxRetries != defaultMaxRetries {
		t.Errorf("expected max retries %d, got %d", defaultMaxRetries, client.options.MaxRetries)
	}
	if client.options.RetryDelay != defaultRetryDelay {
		t.Errorf("expected retry delay %v, got %v", defaultRetryDelay, client.options.RetryDelay)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	auth := &OAuthTokenAuth{AccessToken: "oauth-token"}
	opts := &ClientOptions{
		Timeout:             5 * time.Second,
		MaxRetries:          5,
		RetryDelay:          2 * time.Second,
		MaxConnsPerHost:     100,
		MaxIdleConnDuration: 5 * time.Second,
	}
	client := NewClient(auth, opts)

	if client.options.Timeout != 5*time.Second {
		t.Errorf("expected timeout %v, got %v", 5*time.Second, client.options.Timeout)
	}
	if client.options.MaxRetries != 5 {
		t.Errorf("expected max retries 5, got %d", client.options.MaxRetries)
	}
	if client.options.RetryDelay != 2*time.Second {
		t.Errorf("expected retry delay %v, got %v", 2*time.Second, client.options.RetryDelay)
	}
	if client.options.MaxConnsPerHost != 100 {
		t.Errorf("expected max conns per host 100, got %d", client.options.MaxConnsPerHost)
	}
	if client.options.MaxIdleConnDuration != 5*time.Second {
		t.Errorf("expected max idle conn duration %v, got %v", 5*time.Second, client.options.MaxIdleConnDuration)
	}
}

func TestNewClientWithPartialOptions(t *testing.T) {
	auth := &BearerTokenAuth{Token: "test"}
	opts := &ClientOptions{
		Timeout: 10 * time.Second,
	}
	client := NewClient(auth, opts)

	if client.options.Timeout != 10*time.Second {
		t.Errorf("expected timeout %v, got %v", 10*time.Second, client.options.Timeout)
	}
	if client.options.MaxRetries != defaultMaxRetries {
		t.Errorf("expected default max retries %d, got %d", defaultMaxRetries, client.options.MaxRetries)
	}
	if client.options.RetryDelay != defaultRetryDelay {
		t.Errorf("expected default retry delay %v, got %v", defaultRetryDelay, client.options.RetryDelay)
	}
	if client.options.MaxConnsPerHost != 512 {
		t.Errorf("expected default max conns per host 512, got %d", client.options.MaxConnsPerHost)
	}
}

func TestSetAuth(t *testing.T) {
	auth1 := &BearerTokenAuth{Token: "token1"}
	client := NewClient(auth1, nil)

	auth2 := &OAuthTokenAuth{AccessToken: "token2"}
	client.SetAuth(auth2)

	if client.auth != auth2 {
		t.Error("auth not updated correctly")
	}
}

func TestBearerTokenAuth_AuthenticateFastHTTP(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "test-token",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &BearerTokenAuth{Token: tt.token}
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseRequest(req)

			err := auth.AuthenticateFastHTTP(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !tt.wantErr {
				authHeader := string(req.Header.Peek("Authorization"))
				expected := "Bearer " + tt.token
				if authHeader != expected {
					t.Errorf("expected header %q, got %q", expected, authHeader)
				}
			}
		})
	}
}

func TestOAuthTokenAuth_AuthenticateFastHTTP(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid oauth token",
			token:   "oauth-access-token",
			wantErr: false,
		},
		{
			name:    "empty oauth token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &OAuthTokenAuth{AccessToken: tt.token}
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseRequest(req)

			err := auth.AuthenticateFastHTTP(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !tt.wantErr {
				authHeader := string(req.Header.Peek("Authorization"))
				expected := "Bearer " + tt.token
				if authHeader != expected {
					t.Errorf("expected header %q, got %q", expected, authHeader)
				}
			}
		})
	}
}

type mockTokenProvider struct {
	token string
	err   error
}

func (m *mockTokenProvider) GetToken(ctx context.Context) (string, error) {
	return m.token, m.err
}

func TestDynamicTokenAuth_AuthenticateFastHTTP(t *testing.T) {
	tests := []struct {
		name     string
		provider TokenProvider
		ctx      context.Context
		wantErr  bool
	}{
		{
			name:     "valid token provider",
			provider: &mockTokenProvider{token: "dynamic-token"},
			ctx:      context.Background(),
			wantErr:  false,
		},
		{
			name:     "provider returns error",
			provider: &mockTokenProvider{err: context.DeadlineExceeded},
			ctx:      context.Background(),
			wantErr:  true,
		},
		{
			name:     "provider returns empty token",
			provider: &mockTokenProvider{token: ""},
			ctx:      context.Background(),
			wantErr:  true,
		},
		{
			name:     "nil context",
			provider: &mockTokenProvider{token: "token"},
			ctx:      nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &DynamicTokenAuth{
				Provider: tt.provider,
				ctx:      tt.ctx,
			}
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseRequest(req)

			err := auth.AuthenticateFastHTTP(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !tt.wantErr && tt.provider != nil {
				token, _ := tt.provider.GetToken(tt.ctx)
				authHeader := string(req.Header.Peek("Authorization"))
				expected := "Bearer " + token
				if authHeader != expected {
					t.Errorf("expected header %q, got %q", expected, authHeader)
				}
			}
		})
	}
}

func TestWithToken(t *testing.T) {
	ctx := context.Background()
	token := "test-token"

	newCtx := WithToken(ctx, token)

	value := newCtx.Value(tokenContextKey)
	if value == nil {
		t.Fatal("expected token in context, got nil")
	}

	tokenStr, ok := value.(string)
	if !ok {
		t.Fatal("expected string token in context")
	}

	if tokenStr != token {
		t.Errorf("expected token %q, got %q", token, tokenStr)
	}
}

func TestContextTokenAuth_AuthenticateFastHTTP(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "valid context with token",
			ctx:     WithToken(context.Background(), "context-token"),
			wantErr: false,
		},
		{
			name:    "nil context",
			ctx:     nil,
			wantErr: true,
		},
		{
			name:    "context without token",
			ctx:     context.Background(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &ContextTokenAuth{ctx: tt.ctx}
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseRequest(req)

			err := auth.AuthenticateFastHTTP(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !tt.wantErr {
				authHeader := string(req.Header.Peek("Authorization"))
				if authHeader == "" {
					t.Error("expected non-empty authorization header")
				}
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	apiErr := &APIError{
		Timestamp: "2024-01-01T00:00:00Z",
		Status:    403,
		ErrorType: "Forbidden",
		Message:   "Access denied",
		Path:      "/api/test",
	}

	errMsg := apiErr.Error()
	expected := "API error 403: Forbidden - Access denied (path: /api/test)"
	if errMsg != expected {
		t.Errorf("expected error message %q, got %q", expected, errMsg)
	}
}

func TestAPIError_IsPermissionDenied(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"forbidden", 403, true},
		{"ok", 200, false},
		{"unauthorized", 401, false},
		{"not found", 404, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := &APIError{Status: tt.status}
			if got := apiErr.IsPermissionDenied(); got != tt.want {
				t.Errorf("IsPermissionDenied() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsUnauthorized(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"unauthorized", 401, true},
		{"ok", 200, false},
		{"forbidden", 403, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := &APIError{Status: tt.status}
			if got := apiErr.IsUnauthorized(); got != tt.want {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsNotFound(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"not found", 404, true},
		{"ok", 200, false},
		{"unauthorized", 401, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := &APIError{Status: tt.status}
			if got := apiErr.IsNotFound(); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsBadRequest(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"bad request", 400, true},
		{"ok", 200, false},
		{"not found", 404, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := &APIError{Status: tt.status}
			if got := apiErr.IsBadRequest(); got != tt.want {
				t.Errorf("IsBadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
