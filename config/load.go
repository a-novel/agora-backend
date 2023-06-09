package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"
)

const (
	ENVDevelopment = "development"
	ENVProduction  = "production"
	ENVTest        = "test"
)

var (
	env  string
	main Config
)

type Config struct {
	App       string `json:"app" yaml:"app"`
	ProjectID string `json:"projectID" yaml:"projectID"`
	API       struct {
		Host string `json:"host" yaml:"host"`
		Port int    `json:"port" yaml:"port"`
	} `json:"api" yaml:"api"`
	Frontend struct {
		URLs   []string `json:"urls" yaml:"urls"`
		Routes struct {
			ValidateEmail    string `json:"validateEmail" yaml:"validateEmail"`
			ValidateNewEmail string `json:"validateNewEmail" yaml:"validateNewEmail"`
			ResetPassword    string `json:"resetPassword" yaml:"resetPassword"`
		} `json:"routes" yaml:"routes"`
	} `json:"frontend" yaml:"frontend"`
	Mailer struct {
		APIKey  string `json:"apiKey" yaml:"apiKey"`
		Sandbox bool   `json:"sandbox" yaml:"sandbox"`
		Sender  struct {
			Email string `json:"email" yaml:"email"`
			Name  string `json:"name" yaml:"name"`
		} `json:"sender" yaml:"sender"`
		Templates struct {
			EmailValidation string `json:"emailValidation" yaml:"emailValidation"`
			EmailUpdate     string `json:"emailUpdate" yaml:"emailUpdate"`
			PasswordReset   string `json:"passwordReset" yaml:"passwordReset"`
		} `json:"templates" yaml:"templates"`
	} `json:"mailer" yaml:"mailer"`
	Postgres struct {
		DSN string `json:"dsn" yaml:"dsn"`
	} `json:"postgres" yaml:"postgres"`
	Buckets struct {
		SecretKeys string `json:"secretKeys" yaml:"secretKeys"`
	} `json:"buckets" yaml:"buckets"`
	Secrets struct {
		Prefix         string        `json:"prefix" yaml:"prefix"`
		Backups        int           `json:"backups" yaml:"backups"`
		UpdateInterval time.Duration `json:"updateInterval" yaml:"updateInterval"`
	} `json:"secrets" yaml:"secrets"`
	Tokens struct {
		TTL        time.Duration `json:"ttl" yaml:"ttl"`
		RenewDelta time.Duration `json:"renewDelta" yaml:"renewDelta"`
	} `json:"tokens" yaml:"tokens"`
	IAM struct {
		ServiceAccounts struct {
			Scheduler []string `json:"scheduler" yaml:"scheduler"`
		} `json:"serviceAccounts" yaml:"serviceAccounts"`
	} `json:"iam" yaml:"iam"`
	Forum struct {
		Search struct {
			CropContent int `json:"cropContent" yaml:"cropContent"`
		} `json:"search" yaml:"search"`
	} `json:"forum" yaml:"forum"`
}

func init() {
	env = strings.ToLower(os.Getenv("ENV"))
	switch env {
	case ENVDevelopment, ENVProduction, ENVTest:
		// Do nothing.
	default:
		env = ENVDevelopment
	}

	if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(configfiles.GenericCFG))), &main); err != nil {
		panic(err.Error())
	}

	switch env {
	case ENVDevelopment:
		if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(configfiles.DevCFG))), &main); err != nil {
			panic(err.Error())
		}
	case ENVProduction:
		if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(configfiles.ProdCFG))), &main); err != nil {
			panic(err.Error())
		}
	}
}

func Main() *Config {
	return &main
}

func Env() string {
	return env
}
