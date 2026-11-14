package shopadmin

import (
	"errors"
	"testing"
)

func TestErrorSentinels(t *testing.T) {
	if ErrStoreRequired == nil || ErrLoggerRequired == nil {
		t.Fatal("expected non-nil sentinel errors")
	}
	if errors.Is(ErrStoreRequired, ErrLoggerRequired) {
		t.Fatal("expected distinct sentinel errors")
	}
}
