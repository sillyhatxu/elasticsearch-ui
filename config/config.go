package config

var Conf = new(Config)

type Config struct {
	ServerHost string `json:"server_host"`
	URL        string `json:"url"`
}
