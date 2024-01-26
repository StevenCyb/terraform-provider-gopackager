package compiler

import "github.com/stretchr/testify/mock"

// MockCompiler is an mock type for the Compiler type.
type MockCompiler struct {
	mock.Mock
}

// Compile is a mock implementation of the Compiler.Compile method.
func (m *MockCompiler) Compile(conf Config) (string, string, error) {
	ret := m.Called(conf)

	return ret.Get(0).(string), ret.Get(1).(string), ret.Error(2) //nolint:forcetypeassert
}
