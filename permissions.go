package anamericano

import (
	"context"
	"fmt"
)

// PermissionCheckResponse 권한 확인 응답을 나타냅니다
type PermissionCheckResponse struct {
	// Allowed 권한이 허용되었는지 여부
	Allowed bool `json:"allowed"`
	// Message 권한 결정에 대한 추가 설명
	Message string `json:"message"`
}

// Permission Anamericano 모델의 권한 튜플을 나타냅니다
type Permission struct {
	// ID 권한의 고유 아이디
	ID int64 `json:"id"`
	// ObjectNamespace 객체의 네임스페이스
	ObjectNamespace string `json:"objectNamespace"`
	// ObjectID 객체의 고유 아이디
	ObjectID string `json:"objectId"`
	// Relation 권한 관계
	Relation string `json:"relation"`
	// SubjectType 주체의 타입
	SubjectType string `json:"subjectType"`
	// SubjectID 주체의 고유 아이디
	SubjectID string `json:"subjectId"`
	// SubjectRelation 그룹 멤버십에 사용되는 선택적 필드
	SubjectRelation *string `json:"subjectRelation,omitempty"`
	// CreatedAt 권한이 생성된 타임스탬프
	CreatedAt string `json:"createdAt,omitempty"`
}

// String 권한의 사람이 읽을 수 있는 형태를 반환합니다
func (p *Permission) String() string {
	if p.SubjectRelation != nil {
		return fmt.Sprintf("%s:%s#%s@%s:%s#%s",
			p.ObjectNamespace, p.ObjectID, p.Relation,
			p.SubjectType, p.SubjectID, *p.SubjectRelation)
	}
	return fmt.Sprintf("%s:%s#%s@%s:%s",
		p.ObjectNamespace, p.ObjectID, p.Relation,
		p.SubjectType, p.SubjectID)
}

// CheckPermission 주체가 객체에 대해 특정 권한을 가지고 있는지 확인합니다.
//
// 예시:
//
//	req := &PermissionCheckRequest{
//	    SubjectType:     "user",
//	    SubjectID:       "hanul",
//	    Relation:        "viewer",
//	    ObjectNamespace: "document",
//	    ObjectID:        "doc1",
//	}
//	resp, err := client.CheckPermission(ctx, req)
//	if err != nil {
//	    return err
//	}
//	if resp.Allowed {
//	    fmt.Println("권한이 허용되었습니다")
//	}
func (c *Client) CheckPermission(ctx context.Context, req *PermissionCheckRequest) (*PermissionCheckResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("permission check request is nil")
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	var resp PermissionCheckResponse
	err := c.doRequest(ctx, "POST", "/api/anamericano/check", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// WritePermission 새로운 권한 관계를 생성합니다.
//
// 예시:
//
//	req := &PermissionWriteRequest{
//	    ObjectNamespace: "document",
//	    ObjectID:        "doc1",
//	    Relation:        "editor",
//	    SubjectType:     "user",
//	    SubjectID:       "koyun",
//	}
//	perm, err := client.WritePermission(ctx, req)
func (c *Client) WritePermission(ctx context.Context, req *PermissionWriteRequest) (*Permission, error) {
	if req == nil {
		return nil, fmt.Errorf("permission write request is nil")
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	var perm Permission
	err := c.doRequest(ctx, "POST", "/api/anamericano/write", req, &perm)
	if err != nil {
		return nil, err
	}

	return &perm, nil
}

// DeletePermission 권한 관계를 제거합니다.
//
// 예시:
//
//	req := &PermissionDeleteRequest{
//	    ObjectNamespace: "document",
//	    ObjectID:        "doc1",
//	    Relation:        "viewer",
//	    SubjectType:     "user",
//	    SubjectID:       "hanul",
//	}
//	err := client.DeletePermission(ctx, req)
func (c *Client) DeletePermission(ctx context.Context, req *PermissionDeleteRequest) error {
	if req == nil {
		return fmt.Errorf("permission delete request is nil")
	}

	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	return c.doRequest(ctx, "DELETE", "/api/anamericano/delete", req, nil)
}

// ReadPermissions 특정 객체에 대한 모든 권한을 가져옵니다.
//
// 예시:
//
//	perms, err := client.ReadPermissions(ctx, "document", "doc1")
//	for _, p := range perms {
//	    fmt.Printf("%s:%s가 %s 권한을 가지고 있습니다\n", p.SubjectType, p.SubjectID, p.Relation)
//	}
func (c *Client) ReadPermissions(ctx context.Context, req *PermissionReadRequest) ([]Permission, error) {
	if req == nil {
		return nil, fmt.Errorf("permission read request is nil")
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/api/anamericano/read/%s/%s", req.ObjectNamespace, req.ObjectID)
	var perms []Permission
	return perms, c.doRequest(ctx, "GET", path, nil, &perms)
}

// ExpandPermissions 객체에 대해 특정 관계를 가진 모든 주체를 가져옵니다.
// "type:id" 또는 "type:id#relation" 형식으로 주체의 아이디를 반환합니다.
//
// 예시:
//
//	subjects, err := client.ExpandPermissions(ctx, "document", "doc1", "viewer")
//	// 반환값: ["user:hanul", "user:koyun", "group:ana#member"]
func (c *Client) ExpandPermissions(ctx context.Context, req *PermissionExpendRequest) ([]string, error) {
	if req == nil {
		return nil, fmt.Errorf("permission expend request is nil")
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/api/anamericano/expand/%s/%s/%s", req.ObjectNamespace, req.ObjectID, req.Relation)
	var subjects []string
	return subjects, c.doRequest(ctx, "GET", path, nil, &subjects)
}

// ListObjects 주체가 특정 관계를 가진 모든 객체를 가져옵니다.
// 주어진 관계로 주체가 접근할 수 있는 객체 ID를 반환합니다.
//
// 예시:
//
//	objects, err := client.ListObjects(ctx, "user", "hanul", "viewer", "document")
//	// 반환값: ["doc1", "doc2", "doc5"]
func (c *Client) ListObjects(ctx context.Context, req *ListObjectsRequest) ([]string, error) {
	if req == nil {
		return nil, fmt.Errorf("permission list request is nil")
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("/api/anamericano/list/%s/%s/%s/%s", req.SubjectType, req.SubjectID, req.Relation, req.ObjectNamespace)
	var objects []string
	return objects, c.doRequest(ctx, "GET", path, nil, &objects)
}
