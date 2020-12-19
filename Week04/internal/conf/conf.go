package conf

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Base Common   `toml:"common"`
	DB   Database `json:"db"`
}

type Common struct {
	HttpAddr string `json:"addr" toml:"addr"`
}

type Database struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"passwd"`
	Name     string `toml:"name"`
	Active   int    `toml:"active"`
	Idle     int    `json:"idle"`
}

func Decode(confPath string) (c Config, err error) {
	_, err = toml.DecodeFile(confPath, &c)
	return c, err
}
