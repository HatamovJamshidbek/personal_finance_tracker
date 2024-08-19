package firebase

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type FirebaseStorage struct {
	client     *storage.Client
	bucketName string
}

func NewFirebaseStorage(credentialsFile, bucketName string) (*FirebaseStorage, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &FirebaseStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (fs *FirebaseStorage) UploadFile(ctx context.Context, file io.Reader, fileName string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	objectName := "uploads/" + strconv.FormatInt(rand.Int63(), 10) + "-" + fileName

	token, err := generateRandomToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate download token: %w", err)
	}

	metadata := map[string]string{
		"firebaseStorageDownloadTokens": token,
	}

	bucket := fs.client.Bucket(fs.bucketName)
	object := bucket.Object(objectName)
	wc := object.NewWriter(ctx)
	wc.Metadata = metadata

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write to storage: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close storage writer: %w", err)
	}

	encodedObjectName := url.QueryEscape(objectName)
	fileURL := "https://firebasestorage.googleapis.com/v0/b/" + fs.bucketName + "/o/" + encodedObjectName + "?alt=media" + "&token=" + token

	return fileURL, nil
}

func generateRandomToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (fs *FirebaseStorage) DeleteFile(ctx context.Context, objectName string) error {
	bucket := fs.client.Bucket(fs.bucketName)
	object := bucket.Object(objectName)

	log.Printf("Attempting to delete object: %s", objectName) // Логирование имени объекта

	if err := object.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object %s: %w", objectName, err)
	}

	return nil
}
