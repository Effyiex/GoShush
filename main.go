package main

import (
	"effyiex/shush/server"
	"flag"
	"os"

	"github.com/BurntSushi/toml"
)

const DEFAULT_CONFIG_FILE = "config.toml"

func LoadConfig(path string) server.HostConfig {
	var _, err = os.Open(path)
	if err != nil {
		var cfg server.HostConfig = server.HostConfig{
			Port:        "80",
			Label:       "GoShush",
			ShortLabel:  "Shush",
			FrontColor:  [3]byte{155, 255, 155},
			AccentColor: [3]byte{255, 155, 55},
			BackColor:   [3]byte{33, 33, 39},
			DataFolder:  "shush-data",
		}
		file, _ := os.Create(path)
		toml.NewEncoder(file).Encode(cfg)
		return cfg
	} else {
		data, _ := os.ReadFile(path)
		var cfg server.HostConfig
		toml.Unmarshal(data, &cfg)
		return cfg
	}
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "cfg", DEFAULT_CONFIG_FILE, "path to config file")
	flag.Parse()
	server.NewServer(LoadConfig(configPath)).Run()
}
