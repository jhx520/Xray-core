package core_test

import (
	"context"
	"testing"

	. "github.com/yuzuki616/xray-core/core"
	_ "unsafe"
)

func TestFromContextPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expect panic, but nil")
		}
	}()

	MustFromContext(context.Background())
}
