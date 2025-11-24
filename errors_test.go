package anamericano

import (
	"testing"
)

func TestErrorVariables(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{
			name: "ObjectNameSpaceRequired",
			err:  ObjectNameSpaceRequired,
			msg:  "objectNamespace is required",
		},
		{
			name: "ObjectIdRequired",
			err:  ObjectIdRequired,
			msg:  "objectId is required",
		},
		{
			name: "RelationRequired",
			err:  RelationRequired,
			msg:  "relation is required",
		},
		{
			name: "SubjectIdRequired",
			err:  SubjectIdRequired,
			msg:  "subjectId is required",
		},
		{
			name: "SubjectTypeRequired",
			err:  SubjectTypeRequired,
			msg:  "subjectType is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("error should not be nil")
			}
			if tt.err.Error() != tt.msg {
				t.Errorf("expected error message %q, got %q", tt.msg, tt.err.Error())
			}
		})
	}
}
