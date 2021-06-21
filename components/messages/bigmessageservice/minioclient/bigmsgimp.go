/*
Implementation
*/

package minioclient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
)

// блокирующий метод сохранения данных
func (s *MinioClient) Put(key []byte, data []byte) (proxydata []byte, err error) {

	bucket := cfg.BucketName

	log.Printf(`Minio put: bucket "%s", key: "%s", len: %d bytes`, bucket, key, len(data))

	// проверка наличия коллекции
	found, err := s.client.BucketExists(shortContext(), bucket)
	if err != nil {
		return nil, errors.Wrap(err, "BucketExists fail")
	}

	// создание коллекции
	if !found {
		err := s.client.MakeBucket(shortContext(), bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "MakeBucket fail")
		}
	}

	// сохранение объекта
	r := bytes.NewReader(data)
	_, err = s.client.PutObject(shortContext(), bucket, string(key), r, r.Size(), minio.PutObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "PutObject fail")
	}

	proxydata = append(s.token, key...)
	return proxydata, nil
}

// блокирующий метод получения данных
func (s *MinioClient) Get(key []byte) ([]byte, error) {

	bucket := cfg.BucketName

	// проверка наличия коллекции
	found, err := s.client.BucketExists(shortContext(), bucket)
	if err != nil {
		return nil, errors.Wrap(err, "BucketExists fail")
	}
	if !found {
		return nil, errors.Wrap(err, fmt.Sprintf("Bucket not found: %s", bucket))
	}

	minioObj, err := s.client.GetObject(shortContext(), bucket, string(key), minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "GetObject  fail")
	}

	buf, err := ioutil.ReadAll(minioObj)
	if err != nil {
		return nil, errors.Wrap(err, "Read object fail")
	}

	log.Printf(`Minio get: bucket "%s", key: "%s", len: %d bytes`, bucket, key, len(buf))

	return buf, nil

}

// удаление данных
func (s *MinioClient) Remove(key []byte) error {
	return s.client.RemoveObject(shortContext(), cfg.BucketName, string(key), minio.RemoveObjectOptions{})
}

// проверка, что принятое сообщение является прокси сообщением
func (s *MinioClient) DataIsProxy(data []byte) bool {
	// проверка длины сообщения
	if len(data) < len(s.token) {
		return false
	}
	// проверка соответствия первой части сообщения токену
	return bytes.Equal(data[len(s.token):], s.token)
}

// Получение ключа из прокси сообщения
func (s *MinioClient) GetKey(proxyData []byte) []byte {
	if len(proxyData) < len(s.token) {
		return nil
	}
	return proxyData[len(s.token):]
}

// контекст с коротким таймутом
func shortContext() context.Context {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
	return ctx
}
