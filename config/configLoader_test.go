package config

import (
	"errors"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestReadConfigsFromYaml(t *testing.T) {
	type testCase struct {
		name string
		args string
		want ConfigsValue
	}
	tests := []testCase{
		{
			name: "Load Config 1",
			args: "./testdata/config_test.yaml",
			want: ConfigsValue{
				Password: "senha_super_secreto",
				Url:      "http://url_do_transmission/transmission/rpc/",
				Username: "usuário",
			},
		},
		{
			name: "Load Config 2",
			args: "./testdata/config_test2.yaml",
			want: ConfigsValue{
				Password: "senha_super_secreta",
				Url:      "http://url_do_transmission2/transmission/rpc/",
				Username: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoaderConfigs(tt.args); err != nil {
				t.Errorf("%s", err)
			}
			if Config != tt.want {
				t.Errorf("Expected %s do not match actual %s", tt.want, Config)
			}
		})
	}
}

func TestReadConfigsFromYamlNotFound(t *testing.T) {
	err := LoaderConfigs("./testdata/nao_existe.yaml")
	if err == nil {
		t.Errorf("File %s, shouldn't exist", "./testdata/nao_existe.yaml")
	}
}

func TestErrorDecodeYaml(t *testing.T) {
	var typeError *yaml.TypeError
	err := LoaderConfigs("./testdata/config_test3.yaml")
	if !errors.As(err, &typeError) {
		t.Errorf("Error returned was %v", err)
	}
}
