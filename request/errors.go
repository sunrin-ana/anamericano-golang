package request

import "errors"

var (
	ObjectNameSpaceRequired = errors.New("objectNamespace is required")
	ObjectIdRequired        = errors.New("objectId is required")
	RelationRequired        = errors.New("relation is required")
	SubjectIdRequired       = errors.New("subjectId is required")
	SubjectTypeRequired     = errors.New("subjectType is required")
)
