package anamericano

import (
	"context"
	"testing"
)

func TestPermission_String(t *testing.T) {
	tests := []struct {
		name       string
		permission Permission
		want       string
	}{
		{
			name: "direct permission",
			permission: Permission{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			want: "document:doc1#viewer@user:hanul",
		},
		{
			name: "indirect permission with subject relation",
			permission: Permission{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectType:     "group",
				SubjectID:       "team-alpha",
				SubjectRelation: stringPtr("member"),
			},
			want: "document:doc1#viewer@group:team-alpha#member",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.permission.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestClient_CheckPermission(t *testing.T) {
	auth := &BearerTokenAuth{Token: "test-token"}
	client := NewClient(auth, nil)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *PermissionCheckRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid request - missing subject type",
			req: &PermissionCheckRequest{
				SubjectID:       "hanul",
				Relation:        "viewer",
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing subject id",
			req: &PermissionCheckRequest{
				SubjectType:     "user",
				Relation:        "viewer",
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing relation",
			req: &PermissionCheckRequest{
				SubjectType:     "user",
				SubjectID:       "hanul",
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing object namespace",
			req: &PermissionCheckRequest{
				SubjectType: "user",
				SubjectID:   "hanul",
				Relation:    "viewer",
				ObjectID:    "doc1",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing object id",
			req: &PermissionCheckRequest{
				SubjectType:     "user",
				SubjectID:       "hanul",
				Relation:        "viewer",
				ObjectNamespace: "document",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.CheckPermission(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_WritePermission(t *testing.T) {
	auth := &BearerTokenAuth{Token: "test-token"}
	client := NewClient(auth, nil)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *PermissionWriteRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid request - missing object namespace",
			req: &PermissionWriteRequest{
				ObjectID:    "doc1",
				Relation:    "viewer",
				SubjectType: "user",
				SubjectID:   "hanul",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing object id",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				Relation:        "viewer",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing relation",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing subject type",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectID:       "hanul",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing subject id",
			req: &PermissionWriteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectType:     "user",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.WritePermission(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WritePermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeletePermission(t *testing.T) {
	auth := &BearerTokenAuth{Token: "test-token"}
	client := NewClient(auth, nil)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *PermissionDeleteRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid request - missing object namespace",
			req: &PermissionDeleteRequest{
				ObjectID:    "doc1",
				Relation:    "viewer",
				SubjectType: "user",
				SubjectID:   "hanul",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing object id",
			req: &PermissionDeleteRequest{
				ObjectNamespace: "document",
				Relation:        "viewer",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing relation",
			req: &PermissionDeleteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				SubjectType:     "user",
				SubjectID:       "hanul",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing subject type",
			req: &PermissionDeleteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectID:       "hanul",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing subject id",
			req: &PermissionDeleteRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
				Relation:        "viewer",
				SubjectType:     "user",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeletePermission(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ReadPermissions(t *testing.T) {
	auth := &BearerTokenAuth{Token: "test-token"}
	client := NewClient(auth, nil)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *PermissionReadRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid request - missing object namespace",
			req: &PermissionReadRequest{
				ObjectID: "doc1",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing object id",
			req: &PermissionReadRequest{
				ObjectNamespace: "document",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.ReadPermissions(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPermissions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ExpandPermissions(t *testing.T) {
	auth := &BearerTokenAuth{Token: "test-token"}
	client := NewClient(auth, nil)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *PermissionExpendRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid request - missing object namespace",
			req: &PermissionExpendRequest{
				ObjectID: "doc1",
				Relation: "viewer",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing object id",
			req: &PermissionExpendRequest{
				ObjectNamespace: "document",
				Relation:        "viewer",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing relation",
			req: &PermissionExpendRequest{
				ObjectNamespace: "document",
				ObjectID:        "doc1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.ExpandPermissions(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandPermissions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ListObjects(t *testing.T) {
	auth := &BearerTokenAuth{Token: "test-token"}
	client := NewClient(auth, nil)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *ListObjectsRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid request - missing subject type",
			req: &ListObjectsRequest{
				SubjectID:       "hanul",
				Relation:        "viewer",
				ObjectNamespace: "document",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing subject id",
			req: &ListObjectsRequest{
				SubjectType:     "user",
				Relation:        "viewer",
				ObjectNamespace: "document",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing relation",
			req: &ListObjectsRequest{
				SubjectType:     "user",
				SubjectID:       "hanul",
				ObjectNamespace: "document",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing object namespace",
			req: &ListObjectsRequest{
				SubjectType: "user",
				SubjectID:   "hanul",
				Relation:    "viewer",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.ListObjects(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListObjects() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
