package files

import (
	"testing"

	"github.com/domainry/domainry-foundation/apperror"
)

func TestErrorClassificationContract(t *testing.T) {
	tests := []struct {
		class ErrorClass
		kind  apperror.ErrorKind
		code  string
	}{
		{"bad_request", apperror.KindBadRequest, "knowledge.file_invalid"},
		{"forbidden", apperror.KindForbidden, "knowledge.file_access_denied"},
		{"conflict", apperror.KindConflict, "knowledge.file_fence_rejected"},
		{"unavailable", apperror.KindUnavailable, "knowledge.file_store_unavailable"},
	}
	for _, test := range tests {
		err := &Error{Class: test.class, Code: test.code}
		if apperror.KindOf(err) != test.kind || apperror.CodeOf(err) != test.code {
			t.Fatalf("class=%q kind=%q code=%q", test.class, apperror.KindOf(err), apperror.CodeOf(err))
		}
	}
}
