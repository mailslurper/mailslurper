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

	"github.com/mailslurper/mailslurper/pkg/auth/authscheme"
)

/*
The Configuration structure represents a JSON
configuration file with settings for how to bind
servers and connect to databases.
*/
type Configuration struct {
	WWWAddress       string `json:"wwwAddress"`
	WWWPort          int    `json:"wwwPort"`
	WWWPublicURL     string `json:"wwwPublicURL"`
	ServiceAddress   string `json:"serviceAddress"`
	ServicePort      int    `json:"servicePort"`
	ServicePublicURL string `json:"servicePublicURL"`
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

	AuthSecret           string            `json:"authSecret"`
	AuthSalt             string            `json:"authSalt"`
	AuthenticationScheme string            `json:"authenticationScheme"`
	AuthTimeoutInMinutes int               `json:"authTimeoutInMinutes"`
	Credentials          map[string]string `json:"credentials"`

	StorageType StorageType `json:"-"`
}

var ErrInvalidAdminAddress = fmt.Errorf("Invalid administrator address: wwwAddress")
var ErrInvalidServiceAddress = fmt.Errorf("Invalid service address: serviceAddress")
var ErrInvalidSMTPAddress = fmt.Errorf("Invalid SMTP address: smtpAddress")
var ErrInvalidDBEngine = fmt.Errorf("Invalid DB engine. Valid values are 'SQLite', 'MySQL', 'MSSQL': dbEngine")
var ErrInvalidDBHost = fmt.Errorf("Invalid DB host: dbHost")
var ErrInvalidDBFileName = fmt.Errorf("Invalid DB file name: dbDatabase")
var ErrKeyFileNotFound = fmt.Errorf("Key file not found: keyFile")
var ErrCertFileNotFound = fmt.Errorf("Certificate file not found: certFile")
var ErrNeedCertPair = fmt.Errorf("Please provide both a key file and a cert file: keyFile, certFile")
var ErrAdminKeyFileNotFound = fmt.Errorf("Administrator key file not found: adminKeyFile")
var ErrAdminCertFileNotFound = fmt.Errorf("Adminstartor certificate file not found: adminCertFile")
var ErrNeedAdminCertPair = fmt.Errorf("Please provide both a key file and a cert file: adminKeyFile, adminCertFile")
var ErrInvalidAuthScheme = fmt.Errorf("Invalid authentication scheme. Valid values are 'basic': authenticationScheme")
var ErrMissingAuthSecret = fmt.Errorf("Missing authentication secret. An authentication secret is requried when authentication is enabled: authSecret")
var ErrMissingAuthSalt = fmt.Errorf("Missing authentication salt. A salt value is required when authentication is enabled: authSalt")
var ErrNoUsersConfigured = fmt.Errorf("No users configured. When authentication is enabled you must have at least 1 valid user: credentials")

/*
GetDatabaseConfiguration returns a pointer to a DatabaseConnection structure with data
pulled from a Configuration structure.
*/
func (config *Configuration) GetDatabaseConfiguration() (StorageType, *ConnectionInformation) {
	connectionInformation := NewConnectionInformation(config.DBHost, config.DBPort)
	connectionInformation.SetDatabaseInformation(config.DBDatabase, config.DBUserName, config.DBPassword)

	if strings.ToLower(config.DBEngine) == "sqlite" {
		connectionInformation.SetDatabaseFile(config.DBDatabase)
	}

	result, err := GetDatabaseEngineFromName(config.DBEngine)
	if err != nil {
		panic("Unable to determine database engine")
	}

	return result, connectionInformation
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
GetPublicServiceURL returns a full protocol, address, and port for the MailSlurper service
*/
func (config *Configuration) GetPublicServiceURL() string {
	if config.ServicePublicURL != "" {
		return config.ServicePublicURL
	}

	result := "http"

	if config.CertFile != "" && config.KeyFile != "" {
		result += "s"
	}

	result += fmt.Sprintf("://%s:%d", config.ServiceAddress, config.ServicePort)
	return result
}

/*
GetPublicWWWURL returns a full protocol, address and port for the web application
*/
func (config *Configuration) GetPublicWWWURL() string {
	if config.WWWPublicURL != "" {
		return config.WWWPublicURL
	}

	result := "http"

	if config.AdminCertFile != "" && config.AdminKeyFile != "" {
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

func (config *Configuration) Validate() error {
	if config.WWWAddress == "" {
		return ErrInvalidAdminAddress
	}

	if config.ServiceAddress == "" {
		return ErrInvalidServiceAddress
	}

	if config.SMTPAddress == "" {
		return ErrInvalidSMTPAddress
	}

	if !IsValidStorageType(config.DBEngine) {
		return ErrInvalidDBEngine
	}

	if NeedDBHost(config.DBEngine) {
		if config.DBHost == "" {
			return ErrInvalidDBHost
		}
	}

	if config.DBDatabase == "" {
		return ErrInvalidDBFileName
	}

	if (config.KeyFile == "" && config.CertFile != "") || (config.KeyFile != "" && config.CertFile == "") {
		return ErrNeedCertPair
	}

	if config.KeyFile != "" && config.CertFile != "" {
		if !config.isValidFile(config.KeyFile) {
			return ErrKeyFileNotFound
		}

		if !config.isValidFile(config.CertFile) {
			return ErrCertFileNotFound
		}
	}

	if config.AdminKeyFile != "" && config.AdminCertFile != "" {
		if !config.isValidFile(config.AdminKeyFile) {
			return ErrAdminKeyFileNotFound
		}

		if !config.isValidFile(config.AdminCertFile) {
			return ErrAdminCertFileNotFound
		}
	}

	if (config.AdminKeyFile == "" && config.AdminCertFile != "") || (config.AdminKeyFile != "" && config.AdminCertFile == "") {
		return ErrNeedAdminCertPair
	}

	if config.AuthenticationScheme != "" {
		if !authscheme.IsValidAuthScheme(config.AuthenticationScheme) {
			return ErrInvalidAuthScheme
		}

		if config.AuthSecret == "" {
			return ErrMissingAuthSecret
		}

		if config.AuthSalt == "" {
			return ErrMissingAuthSalt
		}

		if len(config.Credentials) < 1 {
			return ErrNoUsersConfigured
		}
	}

	return nil
}

func (config *Configuration) isValidFile(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}
