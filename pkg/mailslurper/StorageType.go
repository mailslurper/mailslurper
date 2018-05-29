package mailslurper

import "strings"

/*
StorageType defines types of database engines MailSlurper supports
*/
type StorageType int

const (
	STORAGE_MSSQL StorageType = iota
	STORAGE_SQLITE
	STORAGE_MYSQL
)

func GetDatabaseEngineFromName(engineName string) (StorageType, error) {
	switch strings.ToLower(engineName) {
	case "mssql":
		return STORAGE_MSSQL, nil

	case "mysql":
		return STORAGE_MYSQL, nil

	case "sqlite":
		return STORAGE_SQLITE, nil
	}

	return 0, ErrInvalidDBEngine
}

func IsValidStorageType(storageType string) bool {
	_, err := GetDatabaseEngineFromName(storageType)
	if err != nil {
		return false
	}

	return true
}

func NeedDBHost(storageType string) bool {
	if strings.ToLower(storageType) == "sqlite" {
		return false
	}

	return true
}
