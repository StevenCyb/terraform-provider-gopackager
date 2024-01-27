package hasher

import "github.com/stretchr/testify/mock"

// MockHasher is an mock type for the Hasher type.
type MockHasher struct {
	mock.Mock
}

// Mocks the ReadFile method.
func (m *MockHasher) ReadFile(path string) ([]byte, error) {
	ret := m.Called(path)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}

	return ret.Get(0).([]byte), ret.Error(1) //nolint:forcetypeassert
}

// Mocks the MD5 method.
func (m *MockHasher) MD5(binaryContent []byte) string {
	ret := m.Called(binaryContent)

	return ret.String(0)
}

// Mocks the SHA1 method.
func (m *MockHasher) SHA1(binaryContent []byte) string {
	ret := m.Called(binaryContent)

	return ret.String(0)
}

// Mocks the SHA256 method.
func (m *MockHasher) SHA256(binaryContent []byte) string {
	ret := m.Called(binaryContent)

	return ret.String(0)
}

// Mocks the SHA512 method.
func (m *MockHasher) SHA512(binaryContent []byte) string {
	ret := m.Called(binaryContent)

	return ret.String(0)
}

// Mocks the SHA256Base64 method.
func (m *MockHasher) SHA256Base64(binaryContent []byte) string {
	ret := m.Called(binaryContent)

	return ret.String(0)
}

// Mocks the SHA512Base64 method.
func (m *MockHasher) SHA512Base64(binaryContent []byte) string {
	ret := m.Called(binaryContent)

	return ret.String(0)
}

func (m *MockHasher) CombinedHash(binaryContent []byte) CombinedHash {
	ret := m.Called(binaryContent)

	return ret.Get(0).(CombinedHash) //nolint:forcetypeassert
}
