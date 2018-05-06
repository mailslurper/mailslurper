package mailslurper

/*
A Page is the basis for all displayed HTML pages
*/
type Page struct {
	Error   bool
	Message string
	Theme   string
	Title   string
}
