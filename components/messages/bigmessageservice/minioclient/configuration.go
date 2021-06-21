/*
Конфигурация
*/

package minioclient

import "github.com/alexflint/go-arg"

// Конфигурация модуля.
type Config struct {
	Endpoint        string `arg:"env:MINIO_ENDPOINT"`
	AccessKeyID     string `arg:"env:MINIO_ACCESS_KEY"`
	SecretAccessKey string `arg:"env:MINIO_SECRET_KEY"`
	BucketName      string `arg:"env:MINIO_BIG_MSG_BUCKET"`
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Endpoint: "localhost:9000",
		// Endpoint:        "192.168.31.25:9000",
		AccessKeyID:     "minio",
		SecretAccessKey: "12345678",
	}
	arg.MustParse(cfg)
}
