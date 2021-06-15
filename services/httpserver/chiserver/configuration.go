package chiserver

import (
	"github.com/alexflint/go-arg"
)

// config параметры конфигурации
type Config struct {
	Port int `arg:"env:HTTP_SRV_PORT"` // порт http сервера
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Port: 8080,
	}
	arg.Parse(cfg)
}
