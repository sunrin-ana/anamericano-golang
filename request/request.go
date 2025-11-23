package request

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
		return SubjectTypeRequired
	}
	if r.SubjectID == "" {
		return SubjectIdRequired
	}
	if r.Relation == "" {
		return RelationRequired
	}
	if r.ObjectNamespace == "" {
		return ObjectNameSpaceRequired
	}
	if r.ObjectID == "" {
		return ObjectIdRequired
	}
	return nil
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

// Validate 필요한 필드가 모두 있는지 확인
func (r *PermissionWriteRequest) Validate() error {
	if r.ObjectNamespace == "" {
		return ObjectNameSpaceRequired
	}
	if r.ObjectID == "" {
		return ObjectIdRequired
	}
	if r.Relation == "" {
		return RelationRequired
	}
	if r.SubjectType == "" {
		return SubjectTypeRequired
	}
	if r.SubjectID == "" {
		return SubjectIdRequired
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
		return ObjectNameSpaceRequired
	}
	if r.ObjectID == "" {
		return ObjectIdRequired
	}
	if r.Relation == "" {
		return RelationRequired
	}
	if r.SubjectType == "" {
		return SubjectTypeRequired
	}
	if r.SubjectID == "" {
		return SubjectIdRequired
	}
	return nil
}

// PermissionReadRequest 권한 읽기 요청을 나타냅니다.
// 특정 객체에 대한 모든 권한을 가져옵니다.
type PermissionReadRequest struct {
	// ObjectNamespace 객체의 네임스페이스
	ObjectNamespace string `json:"objectNamespace"`
	// ObjectID 객체의 고유 아이디
	ObjectID string `json:"objectId"`
}

// Validate 필요한 필드가 모두 있는지 확인
func (r *PermissionReadRequest) Validate() error {
	if r.ObjectNamespace == "" {
		return ObjectNameSpaceRequired
	}
	if r.ObjectID == "" {
		return ObjectIdRequired
	}
	return nil
}

type PermissionExpendRequest struct {
	// ObjectNamespace 객체의 네임스페이스
	ObjectNamespace string `json:"objectNamespace"`
	// ObjectID 객체의 고유 아이디
	ObjectID string `json:"objectId"`
	// Relation 삭제할 권한 관계
	Relation string `json:"relation"`
}

func (r *PermissionExpendRequest) Validate() error {
	if r.ObjectNamespace == "" {
		return ObjectNameSpaceRequired
	}
	if r.ObjectID == "" {
		return ObjectIdRequired
	}
	if r.Relation == "" {
		return RelationRequired
	}
	return nil
}

type ListObjectsRequest struct {
	// ObjectNamespace 객체의 네임스페이스
	ObjectNamespace string `json:"objectNamespace"`
	// Relation 삭제할 권한 관계
	Relation string `json:"relation"`
	// SubjectType 주체의 타입
	SubjectType string `json:"subjectType"`
	// SubjectID 주체의 고유 아이디
	SubjectID string `json:"subjectId"`
}

func (r *ListObjectsRequest) Validate() error {
	if r.SubjectType == "" {
		return SubjectTypeRequired
	}
	if r.SubjectID == "" {
		return SubjectIdRequired
	}
	if r.Relation == "" {
		return RelationRequired
	}
	if r.ObjectNamespace == "" {
		return ObjectNameSpaceRequired
	}
	return nil
}
