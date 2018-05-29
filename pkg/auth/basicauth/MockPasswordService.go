// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package basicauth

type MockPasswordService struct {
	FnIsPasswordValid func([]byte, []byte) bool
}

func (m *MockPasswordService) IsPasswordValid(password, storedPassword []byte) bool {
	return m.FnIsPasswordValid(password, storedPassword)
}
