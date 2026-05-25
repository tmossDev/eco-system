package repository

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"tmossDev.github.com/eco-system/shared-components/backend/package/env"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

const (
	defaultMediaBucket       = "product-media"
	defaultMediaPublicPrefix = "/api/product-media"
	thumbnailMaxDimension    = 160
	detailMaxDimension       = 1200
	maxUploadBytes           = 10 << 20
)

type ProductMediaObject struct {
	Body          io.ReadCloser
	ContentType   string
	ContentLength int64
}

type ProductMediaStore interface {
	SaveProductImage(productID uint64, fileName string, body io.Reader) (thumbnailURL string, detailURL string, err error)
	GetObject(objectKey string) (*ProductMediaObject, error)
}

type s3ProductMediaStore struct {
	client     *s3.S3
	bucket     string
	publicBase string
	once       sync.Once
	onceErr    error
}

func NewS3ProductMediaStoreFromEnv() ProductMediaStore {
	endpoint := env.Getenv("S3_ENDPOINT", "http://minio.eco-foundation.svc.cluster.local:9000")
	region := env.Getenv("S3_REGION", "us-east-1")
	accessKey := env.Getenv("S3_ACCESS_KEY_ID", "minioadmin")
	secretKey := env.Getenv("S3_SECRET_ACCESS_KEY", "minioadmin")
	forcePathStyle := strings.EqualFold(env.Getenv("S3_FORCE_PATH_STYLE", "true"), "true")

	awsConfig := aws.NewConfig().
		WithEndpoint(endpoint).
		WithRegion(region).
		WithS3ForcePathStyle(forcePathStyle).
		WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, ""))

	return &s3ProductMediaStore{
		client:     s3.New(session.Must(session.NewSession(awsConfig))),
		bucket:     env.Getenv("PRODUCT_MEDIA_BUCKET", defaultMediaBucket),
		publicBase: strings.TrimRight(env.Getenv("PRODUCT_MEDIA_PUBLIC_BASE_URL", defaultMediaPublicPrefix), "/"),
	}
}

func (store *s3ProductMediaStore) SaveProductImage(productID uint64, fileName string, body io.Reader) (string, string, error) {
	if err := store.ensureBucket(); err != nil {
		return "", "", err
	}

	limitedBody := io.LimitReader(body, maxUploadBytes+1)
	uploadBytes, err := io.ReadAll(limitedBody)
	if err != nil || len(uploadBytes) == 0 || len(uploadBytes) > maxUploadBytes {
		return "", "", types.NewInvalidInputError()
	}

	sourceImage, _, err := image.Decode(bytes.NewReader(uploadBytes))
	if err != nil {
		return "", "", types.NewInvalidInputError()
	}

	baseName := sanitizedBaseName(fileName)
	imageID := uuid.NewString()
	baseKey := fmt.Sprintf("products/%d/%s-%s", productID, imageID, baseName)
	thumbnailKey := baseKey + "-thumb.jpg"
	detailKey := baseKey + "-detail.jpg"

	thumbnail, err := encodeJPEG(resizeToFit(sourceImage, thumbnailMaxDimension), 78)
	if err != nil {
		return "", "", types.NewInternalServerError()
	}

	detail, err := encodeJPEG(resizeToFit(sourceImage, detailMaxDimension), 84)
	if err != nil {
		return "", "", types.NewInternalServerError()
	}

	if err := store.putObject(thumbnailKey, thumbnail); err != nil {
		return "", "", err
	}
	if err := store.putObject(detailKey, detail); err != nil {
		return "", "", err
	}

	return store.publicURL(thumbnailKey), store.publicURL(detailKey), nil
}

func (store *s3ProductMediaStore) GetObject(objectKey string) (*ProductMediaObject, error) {
	cleanKey := strings.TrimLeft(objectKey, "/")
	if cleanKey == "" || strings.Contains(cleanKey, "..") {
		return nil, types.NewInvalidInputError()
	}

	output, err := store.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(store.bucket),
		Key:    aws.String(cleanKey),
	})
	if err != nil {
		logger.Errorf("Unable to get product media object %q: %s", cleanKey, err.Error())
		return nil, types.NewNoTFoundOrNoRecordError()
	}

	contentType := aws.StringValue(output.ContentType)
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(cleanKey))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return &ProductMediaObject{
		Body:          output.Body,
		ContentType:   contentType,
		ContentLength: aws.Int64Value(output.ContentLength),
	}, nil
}

func (store *s3ProductMediaStore) ensureBucket() error {
	store.once.Do(func() {
		_, err := store.client.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(store.bucket),
		})
		if err == nil {
			return
		}

		_, err = store.client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(store.bucket),
		})
		if err != nil && !isBucketAlreadyOwned(err) {
			logger.Errorf("Unable to create product media bucket %q: %s", store.bucket, err.Error())
			store.onceErr = types.NewInternalServerError()
		}
	})

	return store.onceErr
}

func (store *s3ProductMediaStore) putObject(key string, body []byte) error {
	_, err := store.client.PutObject(&s3.PutObjectInput{
		Bucket:       aws.String(store.bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(body),
		ContentType:  aws.String("image/jpeg"),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	})
	if err != nil {
		logger.Errorf("Unable to put product media object %q: %s", key, err.Error())
		return types.NewInternalServerError()
	}

	return nil
}

func (store *s3ProductMediaStore) publicURL(key string) string {
	return store.publicBase + "/" + strings.TrimLeft(key, "/")
}

func isBucketAlreadyOwned(err error) bool {
	var awsErr awserr.Error
	return errors.As(err, &awsErr) && awsErr.Code() == s3.ErrCodeBucketAlreadyOwnedByYou
}

func sanitizedBaseName(fileName string) string {
	baseName := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	baseName = strings.ToLower(baseName)
	var builder strings.Builder
	for _, char := range baseName {
		if char >= 'a' && char <= 'z' || char >= '0' && char <= '9' {
			builder.WriteRune(char)
			continue
		}
		if builder.Len() > 0 && builder.String()[builder.Len()-1] != '-' {
			builder.WriteByte('-')
		}
	}

	cleanName := strings.Trim(builder.String(), "-")
	if cleanName == "" {
		return "product"
	}

	if len(cleanName) > 60 {
		return cleanName[:60]
	}

	return cleanName
}

func resizeToFit(source image.Image, maxDimension int) image.Image {
	bounds := source.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= maxDimension && height <= maxDimension {
		return source
	}

	scale := float64(maxDimension) / float64(width)
	if height > width {
		scale = float64(maxDimension) / float64(height)
	}
	targetWidth := max(1, int(float64(width)*scale))
	targetHeight := max(1, int(float64(height)*scale))
	target := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	for y := 0; y < targetHeight; y++ {
		sourceY := bounds.Min.Y + y*height/targetHeight
		for x := 0; x < targetWidth; x++ {
			sourceX := bounds.Min.X + x*width/targetWidth
			target.Set(x, y, source.At(sourceX, sourceY))
		}
	}

	return target
}

func encodeJPEG(source image.Image, quality int) ([]byte, error) {
	var buffer bytes.Buffer
	err := jpeg.Encode(&buffer, source, &jpeg.Options{Quality: quality})
	return buffer.Bytes(), err
}
