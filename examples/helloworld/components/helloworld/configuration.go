/*
 Module configuration
*/

package helloworld

import "github.com/alexflint/go-arg"

type Config struct {
	Name string `arg:"env:HELLO_NAME"`
}

// Check and save the config.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Name: "World", // default value
	}
	arg.Parse(cfg)
}
