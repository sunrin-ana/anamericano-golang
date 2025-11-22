package anamericano

import (
	"context"
	"fmt"
)

// PermissionCheckRequest 권한 확인 요청을 나타냅니다.
// 주체(유저 또는 그룹)가 객체에 대해 특정 관계를 가지고 있는지 확인하는 데 사용됩니다.
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
type PermissionCheckRequest struct {
	// SubjectType 주체의 타입 (예: "user", "group")
	SubjectType string `json:"subjectType"`
	// SubjectID 주체의 고유 아이디
	SubjectID string `json:"subjectId"`
	// Relation 확인할 권한 관계 (예: "viewer", "editor", "owner")
	Relation string `json:"relation"`
	// ObjectNamespace 객체의 네임스페이스 (예: "document", "folder")
	ObjectNamespace string `json:"objectNamespace"`
	// ObjectID 객체의 고유 아이디
	ObjectID string `json:"objectId"`
}

// Validate 요청에 필요한 모든 필드가 있는지 확인합니다
func (r *PermissionCheckRequest) Validate() error {
	if r.SubjectType == "" {
		return fmt.Errorf("subjectType is required")
	}
	if r.SubjectID == "" {
		return fmt.Errorf("subjectId is required")
	}
	if r.Relation == "" {
		return fmt.Errorf("relation is required")
	}
	if r.ObjectNamespace == "" {
		return fmt.Errorf("objectNamespace is required")
	}
	if r.ObjectID == "" {
		return fmt.Errorf("objectId is required")
	}
	return nil
}

// PermissionCheckResponse 권한 확인 응답을 나타냅니다
type PermissionCheckResponse struct {
	// Allowed 권한이 허용되었는지 여부
	Allowed bool `json:"allowed"`
	// Message 권한 결정에 대한 추가 설명
	Message string `json:"message"`
}

// PermissionWriteRequest 권한 쓰기 요청을 나타냅니다.
// 주체와 객체 간의 새로운 권한 관계를 생성합니다.
//
// 예시 (Direct(직접적) 권한):
//
//	req := &PermissionWriteRequest{
//	    ObjectNamespace: "document",
//	    ObjectID:        "doc1",
//	    Relation:        "viewer",
//	    SubjectType:     "user",
//	    SubjectID:       "hanul",
//	}
//
// 예시 (Indirect(간접적) 권한):
//
//	req := &PermissionWriteRequest{
//	    ObjectNamespace: "document",
//	    ObjectID:        "doc1",
//	    Relation:        "viewer",
//	    SubjectType:     "group",
//	    SubjectID:       "team-alpha",
//	    SubjectRelation: stringPtr("member"),
//	}
type PermissionWriteRequest struct {
	// ObjectNamespace 객체의 네임스페이스
	ObjectNamespace string `json:"objectNamespace"`
	// ObjectID 객체의 고유 아이디
	ObjectID string `json:"objectId"`
	// Relation 권한 관계 (예: "viewer", "editor")
	Relation string `json:"relation"`
	// SubjectType 주체의 타입
	SubjectType string `json:"subjectType"`
	// SubjectID 주체의 고유 아이디
	SubjectID string `json:"subjectId"`
	// SubjectRelation 그룹에서 사용되는 선택적 필드?
	SubjectRelation *string `json:"subjectRelation,omitempty"`
}

// 필요한 필드가 모두 있는지 확인
func (r *PermissionWriteRequest) Validate() error {
	if r.ObjectNamespace == "" {
		return fmt.Errorf("objectNamespace is required")
	}
	if r.ObjectID == "" {
		return fmt.Errorf("objectId is required")
	}
	if r.Relation == "" {
		return fmt.Errorf("relation is required")
	}
	if r.SubjectType == "" {
		return fmt.Errorf("subjectType is required")
	}
	if r.SubjectID == "" {
		return fmt.Errorf("subjectId is required")
	}
	return nil
}

// PermissionDeleteRequest 권한 삭제 요청을 나타냅니다.
// 기존 권한 관계를 제거합니다.
type PermissionDeleteRequest struct {
	// ObjectNamespace 객체의 네임스페이스
	ObjectNamespace string `json:"objectNamespace"`
	// ObjectID 객체의 고유 아이디
	ObjectID string `json:"objectId"`
	// Relation 삭제할 권한 관계
	Relation string `json:"relation"`
	// SubjectType 주체의 타입
	SubjectType string `json:"subjectType"`
	// SubjectID 주체의 고유 아이디
	SubjectID string `json:"subjectId"`
}

func (r *PermissionDeleteRequest) Validate() error {
	if r.ObjectNamespace == "" {
		return fmt.Errorf("objectNamespace is required")
	}
	if r.ObjectID == "" {
		return fmt.Errorf("objectId is required")
	}
	if r.Relation == "" {
		return fmt.Errorf("relation is required")
	}
	if r.SubjectType == "" {
		return fmt.Errorf("subjectType is required")
	}
	if r.SubjectID == "" {
		return fmt.Errorf("subjectId is required")
	}
	return nil
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
func (c *Client) ReadPermissions(ctx context.Context, namespace, objectID string) ([]Permission, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if objectID == "" {
		return nil, fmt.Errorf("objectId is required")
	}

	path := fmt.Sprintf("/api/anamericano/read/%s/%s", namespace, objectID)
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
func (c *Client) ExpandPermissions(ctx context.Context, namespace, objectID, relation string) ([]string, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if objectID == "" {
		return nil, fmt.Errorf("objectId is required")
	}
	if relation == "" {
		return nil, fmt.Errorf("relation is required")
	}

	path := fmt.Sprintf("/api/anamericano/expand/%s/%s/%s", namespace, objectID, relation)
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
func (c *Client) ListObjects(ctx context.Context, subjectType, subjectID, relation, namespace string) ([]string, error) {
	if subjectType == "" {
		return nil, fmt.Errorf("subjectType is required")
	}
	if subjectID == "" {
		return nil, fmt.Errorf("subjectId is required")
	}
	if relation == "" {
		return nil, fmt.Errorf("relation is required")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	path := fmt.Sprintf("/api/anamericano/list/%s/%s/%s/%s", subjectType, subjectID, relation, namespace)
	var objects []string
	return objects, c.doRequest(ctx, "GET", path, nil, &objects)
}
