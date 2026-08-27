package processor

import (
	"context"
)

type mockWorkflow struct {
	executeErr error
	called     bool
	bucket     string
	key        string
	callCount  int
}

func (m *mockWorkflow) Execute(
	ctx context.Context,
	inputBucket string,
	inputKey string,
) error {
	m.called = true
	m.bucket = inputBucket
	m.key = inputKey
	m.callCount++

	return m.executeErr
}
