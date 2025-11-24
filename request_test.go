package anamericano

import (
	"testing"
)

func TestPermissionCheckRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *PermissionCheckRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: &PermissionCheckRequest{
				SubjectType:     "user",
				SubjectID:       "hanul",
				Relation:        "viewer",
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: nil,
		},
		{
			name: "missing subject type",
			req: &PermissionCheckRequest{
				SubjectID:       "hanul",
				Relation:        "viewer",
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: SubjectTypeRequired,
		},
		{
			name: "missing subject id",
			req: &PermissionCheckRequest{
				SubjectType:     "user",
				Relation:        "viewer",
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: SubjectIdRequired,
		},
		{
			name: "missing relation",
			req: &PermissionCheckRequest{
				SubjectType:     "user",
				SubjectID:       "hanul",
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: RelationRequired,
		},
		{
			name: "missing object namespace",
			req: &PermissionCheckRequest{
				SubjectType: "user",
				SubjectID:   "hanul",
				Relation:    "viewer",
				ObjectID:    "doc1",
			},
			wantErr: ObjectNameSpaceRequired,
		},
		{
			name: "missing object id",
			req: &PermissionCheckRequest{
				SubjectType:     "user",
				SubjectID:       "hanul",
				Relation:        "viewer",
				ObjectNamespace: "document",
			},
			wantErr: ObjectIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionWriteRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *PermissionWriteRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: nil,
		},
		{
			name: "valid request with subject relation",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectType:     "group",
				SubjectID:       "team-alpha",
				SubjectRelation: stringPtr("member"),
			},
			wantErr: nil,
		},
		{
			name: "missing object namespace",
			req: &PermissionWriteRequest{
				ObjectID:    "doc1",
				Relation:    "viewer",
				SubjectType: "user",
				SubjectID:   "hanul",
			},
			wantErr: ObjectNameSpaceRequired,
		},
		{
			name: "missing object id",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				Relation:        "viewer",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: ObjectIdRequired,
		},
		{
			name: "missing relation",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: RelationRequired,
		},
		{
			name: "missing subject type",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectID:       "hanul",
			},
			wantErr: SubjectTypeRequired,
		},
		{
			name: "missing subject id",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectType:     "user",
			},
			wantErr: SubjectIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionDeleteRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *PermissionDeleteRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: &PermissionDeleteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: nil,
		},
		{
			name: "missing object namespace",
			req: &PermissionDeleteRequest{
				ObjectID:    "doc1",
				Relation:    "viewer",
				SubjectType: "user",
				SubjectID:   "hanul",
			},
			wantErr: ObjectNameSpaceRequired,
		},
		{
			name: "missing object id",
			req: &PermissionDeleteRequest{
				ObjectNamespace: "document",
				Relation:        "viewer",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: ObjectIdRequired,
		},
		{
			name: "missing relation",
			req: &PermissionDeleteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: RelationRequired,
		},
		{
			name: "missing subject type",
			req: &PermissionDeleteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectID:       "hanul",
			},
			wantErr: SubjectTypeRequired,
		},
		{
			name: "missing subject id",
			req: &PermissionDeleteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectType:     "user",
			},
			wantErr: SubjectIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionReadRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *PermissionReadRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: &PermissionReadRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: nil,
		},
		{
			name: "missing object namespace",
			req: &PermissionReadRequest{
				ObjectID: "doc1",
			},
			wantErr: ObjectNameSpaceRequired,
		},
		{
			name: "missing object id",
			req: &PermissionReadRequest{
				ObjectNamespace: "document",
			},
			wantErr: ObjectIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionExpendRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *PermissionExpendRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: &PermissionExpendRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
			},
			wantErr: nil,
		},
		{
			name: "missing object namespace",
			req: &PermissionExpendRequest{
				ObjectID: "doc1",
				Relation: "viewer",
			},
			wantErr: ObjectNameSpaceRequired,
		},
		{
			name: "missing object id",
			req: &PermissionExpendRequest{
				ObjectNamespace: "document",
				Relation:        "viewer",
			},
			wantErr: ObjectIdRequired,
		},
		{
			name: "missing relation",
			req: &PermissionExpendRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: RelationRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListObjectsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *ListObjectsRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: &ListObjectsRequest{
				SubjectType:     "user",
				SubjectID:       "hanul",
				Relation:        "viewer",
				ObjectNamespace: "document",
			},
			wantErr: nil,
		},
		{
			name: "missing subject type",
			req: &ListObjectsRequest{
				SubjectID:       "hanul",
				Relation:        "viewer",
				ObjectNamespace: "document",
			},
			wantErr: SubjectTypeRequired,
		},
		{
			name: "missing subject id",
			req: &ListObjectsRequest{
				SubjectType:     "user",
				Relation:        "viewer",
				ObjectNamespace: "document",
			},
			wantErr: SubjectIdRequired,
		},
		{
			name: "missing relation",
			req: &ListObjectsRequest{
				SubjectType:     "user",
				SubjectID:       "hanul",
				ObjectNamespace: "document",
			},
			wantErr: RelationRequired,
		},
		{
			name: "missing object namespace",
			req: &ListObjectsRequest{
				SubjectType: "user",
				SubjectID:   "hanul",
				Relation:    "viewer",
			},
			wantErr: ObjectNameSpaceRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
