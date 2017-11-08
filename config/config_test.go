package config

import "testing"
import "os"

type TestConfig struct {
	Name       string
	BaseURL    string
	Driver     string
	ConfigFile string
}

func TestGoodConfig(t *testing.T) {

	envValues := []string{"NAME", "BASE_URL", "DRIVER"}

	for i := range envValues {
		os.Unsetenv(envValues[i])
	}

	testCasesGood := []TestConfig{
		{
			Name:       "Test Config HCL",
			BaseURL:    "http://localhost:8089/base_url_hcl",
			Driver:     "testHCL",
			ConfigFile: "../test-config/good/test-config.hcl",
		},
		{
			Name:       "Test Config TOML",
			BaseURL:    "http://localhost:8089/base_url_toml",
			Driver:     "testTOML",
			ConfigFile: "../test-config/good/test-config.toml",
		},
		{
			Name:       "Test Config XML",
			BaseURL:    "http://localhost:8089/base_url_xml",
			Driver:     "testXML",
			ConfigFile: "../test-config/good/test-config.xml",
		},
		{
			Name:       "Test Config YAML",
			BaseURL:    "http://localhost:8089/base_url_yaml",
			Driver:     "testYAML",
			ConfigFile: "../test-config/good/test-config.yml",
		},
		{
			Name:       "Test Config JSON",
			BaseURL:    "http://localhost:8089/base_url_json",
			Driver:     "testJSON",
			ConfigFile: "../test-config/good/test-config.json",
		},
	}

	for i := range testCasesGood {
		cfg, err := NewConfig(testCasesGood[i].ConfigFile)
		if err != nil {
			t.Fatalf("could not read config file %s: %v", testCasesGood[i].ConfigFile, err)
		}

		if cfg.BaseURL != testCasesGood[i].BaseURL {
			t.Errorf("Expected %s, Got %s", testCasesGood[i].BaseURL, cfg.BaseURL)
		}

		if cfg.Driver != testCasesGood[i].Driver {
			t.Errorf("Expected %s, Got %s", testCasesGood[i].Driver, cfg.Driver)
		}

		if cfg.Name != testCasesGood[i].Name {
			t.Errorf("Expected %s, Got %s", testCasesGood[i].Name, cfg.Name)
		}
	}
}

func TestBadConfig(t *testing.T) {

	envValues := []string{"NAME", "BASE_URL", "DRIVER"}

	for i := range envValues {
		os.Unsetenv(envValues[i])
	}

	testCasesBad := []string{
		"../test-config/bad/test-config.hcl",
		"../test-config/bad/test-config.toml",
		"../test-config/bad/test-config.xml",
		"../test-config/bad/test-config.yml",
		"../test-config/bad/test-config.json",
	}

	for i := range testCasesBad {
		_, err := NewConfig(testCasesBad[i])
		if err == nil {
			t.Fatalf("Expected an error")
		}
	}
}

func TestEnvironmentConfig(t *testing.T) {

	envKeys := []string{"NAME", "BASE_URL", "DRIVER"}
	envValues := []string{
		"Testing Environment variables",
		"http://localhost:8088/base_url_env",
		"",
	}

	for i := range envKeys {
		os.Setenv(envKeys[i], envValues[i])
	}

	cfg := &Config{}

	err := cfg.UseCustomEnvConfig()

	if err != nil {
		t.Fatalf("could not read enviroment variable: %v", err)
	}

	if cfg.Name != envValues[0] {
		t.Errorf("Expected %s, Got %s", envValues[0], cfg.Name)
	}

	if cfg.BaseURL != envValues[1] {
		t.Errorf("Expected %s, Got %s", envValues[1], cfg.BaseURL)
	}

	if cfg.Driver != envValues[2] {
		t.Errorf("Expected %s, Got %s", envValues[2], cfg.Driver)
	}

	for i := range envKeys {
		os.Unsetenv(envKeys[i])
	}
}

func TestNoConfigFile(t *testing.T) {

	envValues := []string{"NAME", "BASE_URL", "DRIVER"}

	for i := range envValues {
		os.Unsetenv(envValues[i])
	}

	testCasesBad := []string{
		"../test-config/bad/test-config.hcl1",
	}

	for i := range testCasesBad {
		_, err := NewConfig(testCasesBad[i])
		if err == nil {
			t.Fatalf("Expected an error")
		}
	}
}

func TestInvalidConfigFileExtension(t *testing.T) {

	envValues := []string{"NAME", "BASE_URL", "DRIVER"}

	for i := range envValues {
		os.Unsetenv(envValues[i])
	}

	testCasesBad := []string{
		"../test-config/test-config.hc",
	}

	for i := range testCasesBad {
		_, err := NewConfig(testCasesBad[i])
		if err.Error() != errCfgUnsupported.Error() {
			t.Errorf("Expected \"%v\" but got \"%v\"", errCfgUnsupported, err)
		}
	}
}
