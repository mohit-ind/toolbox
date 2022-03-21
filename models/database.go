package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	constants "github.com/toolboxconstants"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type DatabaseInfo struct {
	Database string `json:"dbname"   yaml:"dbname"`
	Username string `json:"user"     yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host"     yaml:"host"`
	Port     string `json:"port"     yaml:"port"`
	SSLMode  string `json:"sslmode"  yaml:"sslmode"`
}

// Validate DatabaseInfo
func (dbi *DatabaseInfo) Validate() error {
	return validation.ValidateStruct(dbi,
		validation.Field(&dbi.Database, validation.Required, validation.Length(3, 64)),
		validation.Field(&dbi.Username, validation.Required, validation.Length(3, 64)),
		validation.Field(&dbi.Password, validation.Required, validation.Length(3, 64)),
		validation.Field(&dbi.Host, validation.Required, is.Host),
		validation.Field(&dbi.Port, validation.Required, is.Port),
		validation.Field(&dbi.SSLMode, validation.Required, validation.In(constants.ValidSSLModes...)),
	)
}

func NewDatabaseInfo(dbname, user, password, host, port, sslmode string) *DatabaseInfo {
	return &DatabaseInfo{
		Database: dbname,
		Username: user,
		Password: password,
		Host:     host,
		Port:     port,
		SSLMode:  sslmode,
	}
}

func DBInfoFromEnvs(dbname_env, user_env, password_env, host_env, port_env, sslmode_env string) *DatabaseInfo {
	sslmode := os.Getenv(sslmode_env)
	if sslmode == "" {
		sslmode = constants.SSL_MODE_DISABLE
	}
	return NewDatabaseInfo(
		os.Getenv(dbname_env),
		os.Getenv(user_env),
		os.Getenv(password_env),
		os.Getenv(host_env),
		os.Getenv(port_env),
		sslmode,
	)
}

func DBInfoFromSecret(secret string) (*DatabaseInfo, error) {
	var dbSecrets struct {
		Database string      `json:"dbname"`
		Username string      `json:"username"`
		Password string      `json:"password"`
		Host     string      `json:"host"`
		Port     json.Number `json:"port"`
	}
	if err := json.Unmarshal([]byte(secret), &dbSecrets); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal database secrets JSON string")
	}

	return NewDatabaseInfo(
		dbSecrets.Database,
		dbSecrets.Username,
		dbSecrets.Password,
		dbSecrets.Host,
		string(dbSecrets.Port),
		constants.SSL_MODE_DISABLE,
	), nil
}

func DBInfoFromConnectionString(connectionString string) (*DatabaseInfo, error) {
	jsonMap := make(map[string]string)
	elements := strings.Split(connectionString, " ")
	for _, elem := range elements {
		parts := strings.Split(elem, "=")
		if len(parts) != 2 {
			return nil, errors.Errorf("Malformed connection string element: %s", elem)
		}
		jsonMap[parts[0]] = parts[1]
	}
	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to JSON marshal value map")
	}
	var dbInfo DatabaseInfo
	if err := json.Unmarshal(jsonBytes, &dbInfo); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal value map JSON into DatabaseInfo")
	}
	if err := dbInfo.Validate(); err != nil {
		return nil, errors.Wrap(err, "Invalid DatabaseInfo")
	}
	return &dbInfo, nil
}

func (dbi *DatabaseInfo) ConnectionString() string {
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=%s",
		dbi.Database,
		dbi.Username,
		dbi.Password,
		dbi.Host,
		dbi.Port,
		dbi.SSLMode,
	)
}
