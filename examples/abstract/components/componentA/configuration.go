/*
Конфигурация
*/

package componenta

import "github.com/alexflint/go-arg"

// Конфигурация компонента.
type Config struct {
	Param string `arg:"env:COMP_A_PARAM"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Param: "Default Value",
	}
	arg.MustParse(cfg)
}
