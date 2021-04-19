package bootstrap

import (
	"fmt"
)

type testingT struct {
	err error
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	t.err = fmt.Errorf(format, args...) // nolint: goerr113
}

func (t *testingT) LastError() error {
	return t.err
}

func t() *testingT {
	return &testingT{}
}
