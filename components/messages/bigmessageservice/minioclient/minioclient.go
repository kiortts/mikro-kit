package minioclient

import (
	"github.com/kiortts/mikro-kit/components"
	"github.com/kiortts/mikro-kit/components/messages"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type MinioClient struct {
	components.AbstractComponent
	client *minio.Client
	token  []byte
}

// статическая проверка реализаии интерфесов
var _ components.Runnable = (*MinioClient)(nil)
var _ messages.BigMsgService = (*MinioClient)(nil)

// статические переменные
var cfg *Config

func New(config *Config) *MinioClient {
	checkConfig(config)
	s := &MinioClient{
		token: []byte(uuid.NewV3(uuid.NamespaceX500, "MinioClient").String()),
	}
	return s
}

// Запуск клиента в работу.
// Реализация интерфейса components.Runnable.
func (s *MinioClient) Run(mainParams *components.MainParams) error {

	client, err := makeDefaultMinioClient()
	if err != nil {
		return errors.Wrap(err, "getDefaultMinioClient")
	}
	s.client = client

	return nil
}

// создание объекта доступа к данны minio
func makeDefaultMinioClient() (*minio.Client, error) {

	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "minio.New err")
	}

	return minioClient, nil
}
