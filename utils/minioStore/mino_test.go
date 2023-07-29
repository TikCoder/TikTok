package minioStore

import (
	"github.com/minio/minio-go/v6"
	"log"
	"testing"
)

func TestInitMinio(t *testing.T) {
	endpoint := "39.106.47.128:9000"
	accessKeyID := "sODRPVeUKqCsqZe5psBH"
	secretAccessKey := "Ej9XOO4DeKYrW9AFlAcP0uZYCNuNB3lOTWYI8Zlv"
	useSSL := false
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln("创建 MinIO 客户端失败", err)
		return
	}
	log.Printf("创建 MinIO 客户端成功")
	// 创建一个叫 mybucket 的存储桶。
	bucketName := "mybucket"
	location := "beijing"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("存储桶 %s 已经存在", bucketName)
		} else {
			log.Fatalln("查询存储桶状态异常", err)
		}
	}
	log.Printf("创建存储桶 %s 成功", bucketName)

}
