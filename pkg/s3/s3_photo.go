package s3

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client *s3.Client
}

func NewS3() *S3{
	cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion("ru-1"),
        config.WithEndpointResolver(aws.EndpointResolverFunc(
            func(service, region string) (aws.Endpoint, error) {
                return aws.Endpoint{
                    URL:           "https://s3.storage.selcloud.ru",
                    SigningRegion: "ru-1",
                }, nil
            })),
    )
    if err != nil {
        panic(err)
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.UsePathStyle = true
    })
	return &S3{
		client: client,
	}
}


func (s3c *S3) Upload(ctx context.Context, fileHeader *multipart.FileHeader, file multipart.File) (string, error) {
	uploader := manager.NewUploader(s3c.client)

        key := fileHeader.Filename

        _, err := uploader.Upload(ctx, &s3.PutObjectInput{
            Bucket:      aws.String("tinder-ru-1"),
            Key:         aws.String(key),
            Body:        file, // io.Reader, без io.ReadSeeker
            ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
        })
        if err != nil {
			return "", err
        }

        url := "https://70127a1b-bd9a-4d77-b11f-3b28f501631e.selstorage.ru/" + key

		return url, nil
}

