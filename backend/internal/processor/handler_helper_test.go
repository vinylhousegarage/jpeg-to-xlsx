package processor

import "context"

type mockWorkflow struct {
	executeErr  error
	called      bool
	callCount   int
	ctx         context.Context
	inputBucket string
	inputKey    string
	cognitoSub  string
}

func (m *mockWorkflow) Execute(
	ctx context.Context,
	inputBucket string,
	inputKey string,
	cognitoSub string,
) error {
	m.called = true
	m.callCount++
	m.ctx = ctx
	m.inputBucket = inputBucket
	m.inputKey = inputKey
	m.cognitoSub = cognitoSub

	return m.executeErr
}
