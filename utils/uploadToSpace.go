package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func SaveImageToSpaces(objectName, imageData string) error {
    // Decode the base64 data
    decodedData, err := base64.StdEncoding.DecodeString(imageData)
    if err != nil {
        return err
    }

    // Set your DigitalOcean Spaces credentials
    endpoint := os.Getenv("SP_URL")
	accessKey := os.Getenv("SP_ACCESS_KEY_ID")
	secretKey := os.Getenv("SP_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("SP_BUCKET")       

    // Initialize a new Minio client
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: true,
    })
    if err != nil {
        return err
    }

    // Upload the decoded image data to the specified object name in your Space
    contentType := "image/jpeg" 

    // objectName = "products/" + objectName

    _, err = client.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(decodedData), int64(len(decodedData)), minio.PutObjectOptions{
        ContentType: contentType,
        UserMetadata: map[string]string{
            "x-amz-acl": "public-read", 
        },
    })
    if err != nil {
        return err
    }

    return nil
}