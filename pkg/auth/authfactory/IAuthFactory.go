package authfactory

import "github.com/mailslurper/mailslurper/pkg/auth/auth"

/*
IAuthFactory, when implmented, returns the correction authorization
provider based on a configuration setting
*/
type IAuthFactory interface {
	Get() auth.IAuthProvider
}
