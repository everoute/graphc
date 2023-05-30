package informer_test

import (
	"testing"

	"github.com/everoute/graphc/pkg/informer"
	"github.com/everoute/graphc/pkg/schema"
)

type testType struct {
	schema.ObjectMeta

	name string
}

func TestGetTDefaultKeyFuncCrash(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("receive panic: %s, should not crash when receiving an empty object", err)
		}
	}()

	var emptyTask *testType
	_, _ = informer.DefaultKeyFunc(emptyTask)
}
