package natsclient

import "github.com/alexflint/go-arg"

// Конфигурация модуля.
type Config struct {
	Url              string `arg:"env:NATS_URL"`               // адрес сервер NATS
	StreamingCluster string `arg:"env:NATS_STREAMING_CLUSTER"` // имя кластера NATS Streaming
	// BigMessageBucket string `arg:"env:BIG_MESSAGE_BUCKET"`     // карзина для хранения больших сообщений
}

// Проверка и сохранение входяще конфигурации.
func checkConfig(config *Config) {

	if config != nil {
		cfg = config
		return
	}

	cfg = &Config{
		Url:              "",
		StreamingCluster: "nats_cluster",
	}

	arg.Parse(cfg)
}
