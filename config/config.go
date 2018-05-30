package config

import (
	"github.com/naoina/toml"
	"path/filepath"
	"os"
	"fmt"
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
	port := map[bool]string{true: os.Getenv("PORT"), false: "8080"}[os.Getenv("PORT") != ""]
	f := filepath.ToSlash(fp)
	tf, err := os.Open(f)
	defer tf.Close()
	if err != nil {
		conf = Config{
			 SlackConfig{
				Token: os.Getenv("SLACK_TOKEN"),
				},
				ServerConfig{
					Port: fmt.Sprintf(":%s", port),
					EndPoint: os.Getenv("ENDPOINT"),
				},
				SpreadSheetConfig{
					Secret: os.Getenv("CLIENT_JSON"),
					Token: os.Getenv("TOKEN"),
					ID: os.Getenv("SHEET_ID"),
					Name: os.Getenv("SHEET_NAME"),
				},
		}
	} else if err := toml.NewDecoder(tf).Decode(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}