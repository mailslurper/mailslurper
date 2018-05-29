package authscheme

const (
	NONE  string = ""
	BASIC string = "basic"
)

func IsValidAuthScheme(authScheme string) bool {
	if authScheme != BASIC {
		return false
	}

	return true
}
