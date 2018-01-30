package mailslurper

/*
StorageType defines types of database engines MailSlurper supports
*/
type StorageType int

const (
	STORAGE_MSSQL StorageType = iota
	STORAGE_SQLITE
	STORAGE_MYSQL
)
