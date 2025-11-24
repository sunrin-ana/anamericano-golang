# AnAmericano (Golang Wrapper)
> 해당 프로젝트는 2025v AnA ISDT(Internal Service Development Team; 정은수, 이은교)에 의해 추진된 프로젝트입니다


> [!CAUTION]
> 해당 라이브러리는 알파 버전입니다. 사용에 각별한 주의가 필요합니다.

## 설치하기

```bash
go get github.com/sunrin-ana/anamericano-golang
```

## 빠른 시작

### 사용법 (요청마다 다른 토큰 - 권장)

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/sunrin-ana/anamericano-golang"
)

// 전역으로 하나의 클라이언트만 생성 (매우 은교리특검적)
var permissionClient = anamericano.NewClient(&anamericano.ContextTokenAuth{}, nil)

func checkPermissionHandler(w http.ResponseWriter, r *http.Request) {
    // 요청 헤더에서 사용자 토큰 추출
    userToken := r.Header.Get("Authorization")
    
    // 컨텍스트에 토큰 추가
    ctx := anamericano.WithToken(r.Context(), userToken)
    
    // 동일한 클라이언트로 모든 요청 처리
    resp, err := permissionClient.CheckPermission(ctx, &anamericano.PermissionCheckRequest{
        SubjectType:     "user",
        SubjectID:       "hanul",
        Relation:        "viewer",
        ObjectNamespace: "document",
        ObjectID:        "eungyolee-teukcom",
    })
    
    // ...
}
```

### 기존 사용법 (고정 토큰)

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/sunrin-ana/anamericano-golang"
)

func main() {
    // OAuth 사용자 토큰 사용 (권장)
    auth := &anamericano.OAuthTokenAuth{
        AccessToken: "oauth_access_token",
    }
    
    c := anamericano.NewClient(auth, nil)
    ctx := context.Background()

    canView, err := c.CheckPermission(ctx, &anamericano.PermissionCheckRequest{
        SubjectType:     "user",
        SubjectID:       "eungyolee",
        Relation:        "viewer",
        ObjectNamespace: "document",
        ObjectID:        "eungyolee-teukcom",
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Can view: %v\n", canView.Allowed)
}
```

## 메인 컨셉

### ~~절대 Zanzibar를 따라하지 않은~~ An-Americano의 권한 시스템

권한은 다음과 같은 권한 '튜플'로 정의될 수 있습니다: `object#relation@subject`

- **Object**: 타겟 객체 (예시: `document:eungyolee-teukcom`)
- **Relation**: 권한 종류 (예시: `viewer`, `editor`, `owner`)
- **Subject**: 권한 주체 (e.g., `user:eungyolee`)

## API 사용법

### Client 생성

#### 방법 1: 컨텍스트 기반 (가장 효율적 - 권장)
```go
// 하나의 클라이언트를 재사용하고 요청마다 다른 토큰 사용
client := anamericano.NewClient(&anamericano.ContextTokenAuth{}, nil)

// 각 요청마다 컨텍스트에 토큰 전달
ctx := anamericano.WithToken(context.Background(), userToken)
resp, err := client.CheckPermission(ctx, req)
```

#### 방법 2: 고정 토큰
```go
// OAuth 사용자 토큰
auth := &anamericano.OAuthTokenAuth{
    AccessToken: "oauth_access_token",
}
client := anamericano.NewClient(auth, nil)

// 또는 API 토큰 (레거시)
auth := &anamericano.BearerTokenAuth{
    Token: "api_token",
}
client := anamericano.NewClient(auth, nil)
```

#### 방법 3: 동적 토큰 제공자
```go
// 커스텀 토큰 제공자 구현
type MyTokenProvider struct {}

func (p *MyTokenProvider) GetToken(ctx context.Context) (string, error) {
    // 데이터베이스, 캐시 등에서 토큰 가져오기
    return "dynamic_token", nil
}

auth := &anamericano.DynamicTokenAuth{
    Provider: &MyTokenProvider{},
}
client := anamericano.NewClient(auth, nil)
```

#### 커스텀 옵션

```go
// 커스텀 옵션으로 생성
client := anamericano.NewClient(auth, &anamericano.ClientOptions{
    Timeout:    30 * time.Second,
    MaxRetries: 3,
    RetryDelay: 1 * time.Second,
    Logger:     &anamericano.DefaultLogger{},
})
```

### Permission Operations

#### 1. 권한 확인

주체가 객체에 대해 권한이 있는지 확인하기

```go
resp, err := client.CheckPermission(ctx, &anamericano.PermissionCheckRequest{
    SubjectType:     "user",
    SubjectID:       "eungyolee",
    Relation:        "viewer",
    ObjectNamespace: "document",
    ObjectID:        "eungyolee-teukcom",
})

if resp.Allowed {
    // 권한이 있는거임
}
```

#### 2. 권한 업데이트

주체에게 객체에 대한 권한을 업데이트

```go
perm, err := client.WritePermission(ctx, &anamericano.PermissionWriteRequest{
    ObjectNamespace: "document",
    ObjectID:        "eungyolee-teukcom",
    Relation:        "editor",
    SubjectType:     "user",
    SubjectID:       "hanul",
})

fmt.Printf("Created permission: %s\n", perm.String())
```

**그룹 권한 업데이트:**
```go
perm, err := client.WritePermission(ctx, &anamericano.PermissionWriteRequest{
    ObjectNamespace: "document",
    ObjectID:        "eungyolee-teukcom",
    Relation:        "viewer",
    SubjectType:     "group",
    SubjectID:       "ana",
    SubjectRelation: stringPtr("member"),
})
```

#### 3. 권한 삭제

권한을 삭제합니다

```go
err := client.DeletePermission(ctx, &anamericano.PermissionDeleteRequest{
    ObjectNamespace: "document",
    ObjectID:        "eungyolee-teukcom",
    Relation:        "viewer",
    SubjectType:     "user",
    SubjectID:       "hanul",
})
```

#### 4. 권한 읽기

모든 권한을 읽습니다

```go
perms, err := client.ReadPermissions(ctx, "document", "eungyolee-teukcom")
for _, p := range perms {
    fmt.Printf("%s can %s\n", p.SubjectID, p.Relation)
}
```

#### 5. 특정 권한 읽기

특정 권한을 가진 모든 주체를 가져옵니다

```go
subjects, err := client.ExpandPermissions(ctx, "document", "eungyolee-teukcom", "viewer")
// Returns: ["user:eungyolee", "user:hanul", "group:ana#member"]

for _, subject := range subjects {
    fmt.Println(subject)
}
```

#### 6. 모든 객체 가져오기

객체가 접근할 수 있는 모든 객체를 불러옵니다

```go
docs, err := client.ListObjects(ctx, "user", "eungyolee", "viewer", "document")
// Returns: ["eungyolee-teukcom", "eungyolee-babo", "eungyolee-kimanjja"]

for _, docID := range docs {
    fmt.Printf("Eungyolee can view: %s\n", docID)
}
```

## 예외처리

다음과 같이 할 수 있음:

```go
resp, err := client.CheckPermission(ctx, req)
if err != nil {
    if apiErr, ok := err.(*anamericano.APIError); ok {
        switch {
        case apiErr.IsUnauthorized():
            // API Token이 이상함
            fmt.Println("다시 발급하쇼")
        case apiErr.IsPermissionDenied():
            // 발급된 API Token으로 권한이 없음
            fmt.Println("임원진에게 문의하쇼")
        case apiErr.IsNotFound():
            // 객체가 없음
            fmt.Println("ㅁ?ㄹ")
        case apiErr.IsBadRequest():
            // 암튼 클라이언트 오류
            fmt.Printf("Bad request: %s\n", apiErr.Message)
        default:
            // 다른 오류
            fmt.Printf("API error: %s\n", apiErr.Error())
        }
    } else {
        // 네트워크 오류
        fmt.Printf("Error: %v\n", err)
    }
}
```

## 필요 사항

- Go 1.24 또는 그 이상
- AnA ISDT(Internal Service Development Team) 가입

## License

MIT License

## 지원

선린인터넷고등학교 3호관 341실로 찾아오세요

## Changelog

### v0.0.1 (2025-11-22)
- Initial release

### v0.0.2 (2025-11-23)
- Dynamic Token

### v0.0.3 (2025-11-23)
- Refactored library

### v0.0.4 (2025-11-24)
- Reverted Structure
- migration to fasthttp from net/http
- Memeory leak fix
