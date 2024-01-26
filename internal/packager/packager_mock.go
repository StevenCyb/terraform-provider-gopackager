package packager

import "github.com/stretchr/testify/mock"

// MockZIP is a mock implementation of ZIP.
type MockZIP struct {
	mock.Mock
}

// Zip is a mocked method.
func (m *MockZIP) Zip(zipPath string, files map[string]string) (string, error) {
	args := m.Called(zipPath, files)

	return args.String(0), args.Error(1)
}
