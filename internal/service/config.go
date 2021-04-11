package service

import (
	"github.com/bool64/zapctxd"
	"github.com/google/uuid"
)

const (
	// CredentialsProviderKeychain indicates that the credentials provider is keychain.
	CredentialsProviderKeychain = CredentialsProviderType("keychain")
	// CredentialsProviderNone indicates that there is no credentials provider.
	CredentialsProviderNone = CredentialsProviderType("")

	// OutputFormatPrettyJSON is prettified json format.
	OutputFormatPrettyJSON = "pretty-json"
	// OutputFormatJSON is json format.
	OutputFormatJSON = "json"
	// OutputFormatCSV is csv format.
	OutputFormatCSV = "csv"
	// OutputFormatNone is no format.
	OutputFormatNone = ""
)

// Config is a global config for the application.
type Config struct {
	OutputFormat string

	Log zapctxd.Config
	N26 N26Config
}

// N26Config represents configuration for N26 Client.
type N26Config struct {
	Username            string
	Password            string
	Device              uuid.UUID
	CredentialsProvider CredentialsProviderType `toml:"credentials"`
}

// CredentialsProviderType indicates the type of a credentials provider.
type CredentialsProviderType string
