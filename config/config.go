package config

import (
	"github.com/naoina/toml"
	"path/filepath"
	"os"
)

type Config struct {
	Slack SlackConfig
	Server ServerConfig
	SpreadSheet SpreadSheetConfig
}

type SlackConfig struct {
	Token string `toml:"token"`

}

type ServerConfig struct {
	Port string `toml:"port"`
	EndPoint string `toml:"endpoint"`
}

type SpreadSheetConfig struct {
	Secret string `toml:"secret"`
	Token string `toml:"token"`
	ID string `toml:"id"`
	Name string `toml:"name"`
}

func LoadToml(fp string) (*Config, error) {
	var conf Config
	f := filepath.ToSlash(fp)
	tf, err := os.Open(f)
	defer tf.Close()
	if err != nil {
		return nil, err
	}
	if err := toml.NewDecoder(tf).Decode(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}