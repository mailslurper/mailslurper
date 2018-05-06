package basicauth

type MockPasswordService struct {
	FnIsPasswordValid func([]byte, []byte) bool
}

func (m *MockPasswordService) IsPasswordValid(password, storedPassword []byte) bool {
	return m.FnIsPasswordValid(password, storedPassword)
}
