package auth

/*
IPasswordService is an interface for validating passwords
*/
type IPasswordService interface {
	IsPasswordValid(password, storedPassword []byte) bool
}
