// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

/*
The Configuration structure represents a JSON
configuration file with settings for how to bind
servers and connect to databases.
*/
type Configuration struct {
	WWWAddress       string `json:"wwwAddress"`
	WWWPort          int    `json:"wwwPort"`
	ServiceAddress   string `json:"serviceAddress"`
	ServicePort      int    `json:"servicePort"`
	SMTPAddress      string `json:"smtpAddress"`
	SMTPPort         int    `json:"smtpPort"`
	DBEngine         string `json:"dbEngine"`
	DBHost           string `json:"dbHost"`
	DBPort           int    `json:"dbPort"`
	DBDatabase       string `json:"dbDatabase"`
	DBUserName       string `json:"dbUserName"`
	DBPassword       string `json:"dbPassword"`
	MaxWorkers       int    `json:"maxWorkers"`
	AutoStartBrowser bool   `json:"autoStartBrowser"`
	CertFile         string `json:"certFile"`
	KeyFile          string `json:"keyFile"`
	AdminCertFile    string `json:"adminCertFile"`
	AdminKeyFile     string `json:"adminKeyFile"`
	Theme            string `json:"theme"`

	AdminCookieSecret    string            `json:"-"`
	AuthenticationScheme string            `json:"authenticationScheme"`
	Credentials          map[string]string `json:"credentials"`

	StorageType StorageType
}

/*
GetDatabaseConfiguration returns a pointer to a DatabaseConnection structure with data
pulled from a Configuration structure.
*/
func (config *Configuration) GetDatabaseConfiguration() (StorageType, *ConnectionInformation) {
	result := NewConnectionInformation(config.DBHost, config.DBPort)
	result.SetDatabaseInformation(config.DBDatabase, config.DBUserName, config.DBPassword)

	if strings.ToLower(config.DBEngine) == "sqlite" {
		result.SetDatabaseFile(config.DBDatabase)
	}

	return config.getDatabaseEngineFromName(config.DBEngine), result
}

func (config *Configuration) getDatabaseEngineFromName(engineName string) StorageType {
	switch strings.ToLower(engineName) {
	case "mssql":
		return STORAGE_MSSQL

	case "mysql":
		return STORAGE_MYSQL
	}

	return STORAGE_SQLITE
}

/*
GetFullServiceAppAddress returns a full address and port for the MailSlurper service
application.
*/
func (config *Configuration) GetFullServiceAppAddress() string {
	return fmt.Sprintf("%s:%d", config.ServiceAddress, config.ServicePort)
}

/*
GetFullSMTPBindingAddress returns a full address and port for the MailSlurper SMTP
server.
*/
func (config *Configuration) GetFullSMTPBindingAddress() string {
	return fmt.Sprintf("%s:%d", config.SMTPAddress, config.SMTPPort)
}

/*
GetFullWWWBindingAddress returns a full address and port for the Web application.
*/
func (config *Configuration) GetFullWWWBindingAddress() string {
	return fmt.Sprintf("%s:%d", config.WWWAddress, config.WWWPort)
}

/*
GetCORSWWWAddress returns a full protocol, address and port for the web application
*/
func (config *Configuration) GetCORSWWWAddress() string {
	result := "http"

	if config.CertFile != "" && config.KeyFile != "" {
		result += "s"
	}

	result += fmt.Sprintf("://%s:%d", config.WWWAddress, config.WWWPort)
	return result
}

/*
GetTheme returns the configured theme. If there isn't one, the
default theme is used
*/
func (config *Configuration) GetTheme() string {
	theme := config.Theme

	if theme == "" {
		theme = "default"
	}

	return theme
}

/*
LoadConfiguration reads data from a Reader into a new Configuration structure.
*/
func LoadConfiguration(reader io.Reader) (*Configuration, error) {
	var err error
	var buffer = make([]byte, 4096)

	result := &Configuration{}
	if buffer, err = ioutil.ReadAll(reader); err != nil {
		return result, err
	}

	if err = json.Unmarshal(buffer, result); err != nil {
		return result, err
	}

	return result, nil
}

/*
LoadConfigurationFromFile reads data from a file into a Configuration object. Makes use of
LoadConfiguration().
*/
func LoadConfigurationFromFile(fileName string) (*Configuration, error) {
	var err error
	result := &Configuration{}
	var configFileHandle *os.File

	if configFileHandle, err = os.Open(fileName); err != nil {
		return result, err
	}

	if result, err = LoadConfiguration(configFileHandle); err != nil {
		return result, err
	}

	return result, nil
}

/*
SaveConfiguration saves the current state of a Configuration structure
into a JSON file.
*/
func (config *Configuration) SaveConfiguration(configFile string) error {
	var err error
	var serializedConfigFile []byte

	if serializedConfigFile, err = json.Marshal(config); err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, serializedConfigFile, 0644)
}

/*
IsAdminSSL returns true if cert files are provided for the admin
*/
func (config *Configuration) IsAdminSSL() bool {
	return config.AdminKeyFile != "" && config.AdminCertFile != ""
}

/*
IsServiceSSL returns true if cert files are provided for the SMTP server
and the services tier
*/
func (config *Configuration) IsServiceSSL() bool {
	return config.KeyFile != "" && config.CertFile != ""
}
