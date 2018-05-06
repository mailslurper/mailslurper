package auth

/*
IAuthProvider describes a provider of authentication services, such
as Basic, LDAP, etc...
*/
type IAuthProvider interface {
	Login(credentials *AuthCredentials) error
}
