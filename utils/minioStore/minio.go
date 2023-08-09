package minioStore

import (
	"fmt"
	"github.com/minio/minio-go/v6"
	"strconv"
	"strings"
	"tiktok2023/config"
	"time"
)

type Minio struct {
	MinioClient    *minio.Client
	Endpoint       string
	Port           string
	VideoBuckets   string
	PictureBuckets string
}

var MinIO Minio

func InitMinio() error {
	conf := config.Conf.MinIO
	URL := conf.Url
	port := conf.Port
	endpoint := URL + ":" + port
	accessKeyID := conf.AccessKeyId
	secretAccessKey := conf.SecretAccessKey
	videoBucket := conf.VideoBuckets
	pictureBucket := conf.PictureBuckets
	minioClient, err := minio.New(
		endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		panic(err)
	}
	// 创建桶存储
	creatBucket(minioClient, videoBucket)
	creatBucket(minioClient, pictureBucket)

	MinIO = Minio{
		MinioClient:    minioClient,
		Endpoint:       endpoint,
		Port:           port,
		VideoBuckets:   videoBucket,
		PictureBuckets: pictureBucket,
	}
	return nil
}

func creatBucket(m *minio.Client, bucket string) {

	found, err := m.BucketExists(bucket)
	// todo log err
	if err != nil {

	}
	if !found {
		m.MakeBucket(bucket, "us-east-1")
	}
	// 设置桶策略
	policy := `{"Version": "2012-10-17",
				"Statement": 
					[{
						"Action":["s3:GetObject"],
						"Effect": "Allow",
						"Principal": {"AWS": ["*"]},
						"Resource": ["arn:aws:s3:::` + bucket + `/*"],
						"Sid": ""
					}]
				}`
	err = m.SetBucketPolicy(bucket, policy)
	// todo log err
	if err != nil {

	}

	// todo log info

}

// UploadFile 上传到 MinIO
func UploadFile(filetype, file, userID string) (string, error) {
	var fileName strings.Builder
	var contentType, Suffix, bucket string

	if filetype == "video" {
		contentType = "video/mp4"
		Suffix = ".mp4"
		bucket = MinIO.VideoBuckets
	} else {
		contentType = "image/jpeg"
		Suffix = ".jpg"
		bucket = MinIO.PictureBuckets
	}
	fileName.WriteString(userID)
	fileName.WriteString("_")
	fileName.WriteString(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	fileName.WriteString(Suffix)

	n, err := MinIO.MinioClient.FPutObject(bucket, fileName.String(), file, minio.PutObjectOptions{
		ContentType: contentType,
	})
	// todo log err
	if err != nil {
		// todo log err save
		return "", err
	}
	// todo log info
	fmt.Printf("upload file %d byte success,fileName:%s", n, fileName.String())
	url := "http://" + MinIO.Endpoint + "/" + bucket + "/" + fileName.String()
	return url, nil
}
