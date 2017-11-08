package config

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/camelcase"
	"github.com/hashicorp/hcl"
	yaml "gopkg.in/yaml.v2"
)

var errCfgUnsupported = errors.New("config file format not supported. Supported formats are json, xml, yaml, toml, hcl")

// Config for testing
type Config struct {
	Name    string `json:"name" xml:"name" yaml:"name" toml:"name" hcl:"name"`
	BaseURL string `json:"base_url" xml:"base_url" yaml:"base_url" toml:"base_url" hcl:"base_url"`
	Driver  string `json:"driver" xml:"driver" yaml:"driver" toml:"driver" hcl:"driver"`
}

// DefaultConfig returns a new default config
func defaultConfig() *Config {
	return &Config{
		Name:    "No name",
		BaseURL: "No base url",
		Driver:  "No Driver",
	}
}

// NewConfig returns a new config
func NewConfig(path string) (*Config, error) {
	return readConfig(path)
}

// readConfig reads config from the provided file
func readConfig(path string) (*Config, error) {

	cfg := defaultConfig()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(path) {
	case ".json":
		err := json.Unmarshal(data, cfg)
		if err != nil {
			return nil, err
		}
	case ".xml":
		fmt.Println("XML")
		err := xml.Unmarshal(data, cfg)
		if err != nil {
			return nil, err
		}
	case ".yml":
		err := yaml.Unmarshal(data, cfg)
		if err != nil {
			return nil, err
		}
	case ".toml":
		err := toml.Unmarshal(data, cfg)
		if err != nil {
			return nil, err
		}
	case ".hcl":
		obj, err := hcl.Parse(string(data))
		if err != nil {
			return nil, err
		}

		if err = hcl.DecodeObject(&cfg, obj); err != nil {
			return nil, err
		}
	default:
		return nil, errCfgUnsupported
	}

	return cfg, nil
}

func getEnvName(field string) string {
	camSplit := camelcase.Split(field)
	rst := strings.Join(camSplit, "_")
	return strings.ToUpper(rst)
}

// UseCustomEnvConfig updates configs with user environment configs
func (c *Config) UseCustomEnvConfig() error {
	cfg := reflect.ValueOf(c).Elem()
	cTyp := cfg.Type()

	for i := range make([]struct{}, cTyp.NumField()) {
		field := cTyp.Field(i)

		cm := getEnvName(field.Name)
		env := os.Getenv(strings.ToUpper(cm))

		if env == "" {
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			cfg.FieldByName(field.Name).SetString(env)
		}
	}

	return nil
}
