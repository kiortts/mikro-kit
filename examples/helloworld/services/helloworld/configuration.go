package helloworld

import "github.com/alexflint/go-arg"

type Config struct {
	Name string `arg:"env:HELLO_NAME"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Name: "Word",
	}
	arg.Parse(cfg)

}
