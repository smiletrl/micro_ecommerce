package logger

import (
	"go.uber.org/zap"
)

type mockProvider struct {
	logger *zap.SugaredLogger
}

func NewMockProvider() Provider {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	return &mockProvider{sugar}
}

func (m mockProvider) Infow(msg string, keysAndValues ...string) {
	vals := m.convert(keysAndValues...)
	m.logger.Infow(msg, vals...)
}

func (m mockProvider) Warnw(msg string, keysAndValues ...string) {
	vals := m.convert(keysAndValues...)
	m.logger.Warnw(msg, vals...)
}

func (m mockProvider) Debugw(msg string, keysAndValues ...string) {
	vals := m.convert(keysAndValues...)
	m.logger.Debugw(msg, vals...)
}

func (m mockProvider) Errorw(msg string, keysAndValues ...string) {
	vals := m.convert(keysAndValues...)
	m.logger.Errorw(msg, vals...)
}

func (m mockProvider) Fatal(msg string, err error) {
	m.logger.Fatalw(msg, err.Error())
}

func (m mockProvider) Close() {
	m.logger.Sync()
}

func (m mockProvider) convert(keysAndValues ...string) []interface{} {
	vals := make([]interface{}, len(keysAndValues))
	for key, val := range keysAndValues {
		vals[key] = interface{}(val)
	}
	return vals
}
