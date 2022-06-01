package config

import (
	"gopkg.in/dealancer/validate.v2"
	"os"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

var App *Settings

func Init(filePath string) error {
	var err error

	App, err = New(filePath)

	if err != nil {
		return eris.Wrap(err, "global config failed to initialize")
	}

	return nil
}

func New(filePath string) (*Settings, error) {
	if len(filePath) == 0 {
		return nil, eris.New("internal/config: filePath is empty")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, eris.Wrap(err, "internal/config: can not open file")
	}
	defer func() { _ = f.Close() }()

	d := yaml.NewDecoder(f)
	var out *Settings
	if err := d.Decode(&out); err != nil {
		return nil, eris.Wrap(err, "internal/config: can not decode file to yaml struct")
	}

	if err := validate.Validate(&out); err != nil {
		return nil, eris.Wrap(err, "internal/config: the provided configuration is invalid")
	}
	return out, nil
}

type Settings struct {
	Log struct {
		Type  string `yaml:"type"  validate:"one_of=human,json"`
		Level string `yaml:"level" validate:"one_of=debug,info,warning,error"`
	} `yaml:"log"`

	Databases []struct {
		Host     string `yaml:"host" validate:"empty=false"`
		Name     string `yaml:"name" validate:"empty=false"`
		Port     int    `yaml:"port" validate:"ne=0"`
		Database string `yaml:"database" validate:"empty=false"`
		Username string `yaml:"username" validate:"empty=false"`
		Password string `yaml:"password" validate:"empty=false"`
		Version  int    `yaml:"version"  validate:"one_of=12,13,14"`
		Verbose  bool   `yaml:"verbose"`
	} `yaml:"databases"`

	Servers []struct {
		Host     string `yaml:"host" validate:"empty=false"`
		Name     string `yaml:"name" validate:"empty=false"`
		Port     int    `yaml:"port" validate:"ne=0"`
		Username string `yaml:"username" validate:"empty=false"`
		Password string `yaml:"password" validate:"empty=false"`
		Version  int    `yaml:"version"  validate:"one_of=12,13,14"`
	} `yaml:"servers"`

	Outputs struct {
		Minio struct {
			Enabled         bool   `yaml:"enabled"`
			Endpoint        string `yaml:"endpoint"`
			AccessKeyID     string `yaml:"accessKeyID"`
			BucketName      string `yaml:"bucketName"`
			SecretAccessKey string `yaml:"secretAccessKey"`
		} `yaml:"minio"`
		Sftp struct {
			Enabled   bool   `yaml:"enabled"`
			Host      string `yaml:"host"`
			Port      int    `yaml:"port"`
			User      string `yaml:"user"`
			Password  string `yaml:"password"`
			Directory string `yaml:"directory"`
		} `yaml:"sftp"`
		Filesystem struct {
			Enabled bool   `yaml:"enabled"`
			Path    string `yaml:"path"`
		} `yaml:"filesystem"`
	} `yaml:"outputs"`

	Notifiers struct {
		Email struct {
			Enabled       bool   `yaml:"enabled"`
			SMTP          string `yaml:"smtp"`
			Port          int    `yaml:"port"`
			Sender        string `yaml:"sender"`
			Password      string `yaml:"password"`
			SendOnSuccess bool   `yaml:"sendOnSuccess"`
		} `yaml:"email"`
		Slack struct {
			Enabled       bool   `yaml:"enabled"`
			BotToken      string `yaml:"botToken"`
			Channel       string `yaml:"channel"`
			SendOnSuccess bool   `yaml:"sendOnSuccess"`
		} `yaml:"slack"`
	} `yaml:"notifiers"`
}
