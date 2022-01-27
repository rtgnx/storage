package s3

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	Prefix string
	Bucket string
	client *minio.Client
}

func FromCfg(cfg ConnectionProfile, bucket string) (s3 *S3, err error) {
	s3 = new(S3)

	cfg.URL = strings.TrimPrefix(cfg.URL, "https://")
	s3.Bucket = bucket
	s3.client, err = minio.New(
		cfg.URL, &minio.Options{Creds: credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, "")},
	)
	return
}

func FromConfig(alias string, bucket string) (s3 *S3, err error) {
	ReadConfig()
	s3 = new(S3)
	profile := cfg.Aliases[alias]
	profile.URL = strings.TrimPrefix(profile.URL, "https://")
	s3.Bucket = bucket
	s3.client, err = minio.New(
		profile.URL, &minio.Options{Creds: credentials.NewStaticV4(profile.AccessKey, profile.SecretKey, "")},
	)

	return s3, err
}

func FromEnv() (*S3, error) {
	var err error
	s3 := new(S3)
	s3.Bucket = os.Getenv("S3_BUCKET")
	s3.client, err = minio.New(
		os.Getenv("S3_ENDPOINT"), &minio.Options{Creds: credentials.NewStaticV4(os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), "")},
	)

	return s3, err
}

func (s3 *S3) Exists(prefix string) bool {
	_, err := s3.client.GetObject(context.Background(), s3.Bucket, prefix, minio.GetObjectOptions{})

	return err == nil
}

func (s3 *S3) WriteBytes(prefix string, b []byte) error {
	buf := bytes.NewBuffer(b)
	_, err := s3.client.PutObject(context.Background(), s3.Bucket, prefix, buf, int64(buf.Len()), minio.PutObjectOptions{})
	return err
}

func (s3 *S3) ReadBytes(prefix string) ([]byte, error) {
	obj, err := s3.client.GetObject(context.Background(), s3.Bucket, prefix, minio.GetObjectOptions{})

	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(obj)
}

func (s3 *S3) Write(prefix string, r io.Reader) error {
	_, err := s3.client.PutObject(context.Background(), s3.Bucket, prefix, r, -1, minio.PutObjectOptions{})
	return err
}

func (s3 *S3) Read(prefix string, w io.Writer) error {
	obj, err := s3.client.GetObject(context.Background(), s3.Bucket, prefix, minio.GetObjectOptions{})

	if err != nil {
		return err
	}

	_, err = io.Copy(w, obj)
	return err
}

func (s3 *S3) List(prefix string) []string {
	ls := make([]string, 0)

	list := s3.client.ListObjects(context.Background(), s3.Bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true})

	for obj := range list {
		ls = append(ls, obj.Key)
	}

	return ls
}

func (s3 *S3) MakeBucket(name string) error {
	return s3.client.MakeBucket(context.Background(), name, minio.MakeBucketOptions{})
}
