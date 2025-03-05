package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3_ACCESS_KEY string = "AKIA5V6I7DU7U2PMWZLJ"
var S3_SECRET_KEY string = "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"
var S3_BUCKET_NAME string = "uploadpilottest"
var S3_REGION string = "ap-south-1"

type Repo struct {
	s3Client          *s3.Client
	s3PresignedClient *s3.PresignClient
}

func NewS3Client(accessKey string, secretKey string, s3BucketRegion string) *Repo {
	options := s3.Options{
		Region:      s3BucketRegion,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	}

	client := s3.New(options, func(o *s3.Options) {
		o.Region = s3BucketRegion
		o.UseAccelerate = false
	})

	presignClient := s3.NewPresignClient(client)
	return &Repo{
		s3Client:          client,
		s3PresignedClient: presignClient,
	}
}

func (repo Repo) PutObject(bucketName string, objectKey string, lifetimeSecs int64, contentlength int, hash string) (*v4.PresignedHTTPRequest, error) {
	// 1 mb limit
	if contentlength > 1048576 || contentlength == 0 {
		contentlength = 1048576
	}
	fmt.Println("contentlength: ", contentlength)
	request, err := repo.s3PresignedClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		ContentType:   aws.String("image/png"),
		ContentLength: aws.Int64(int64(contentlength)),
		IfNoneMatch:   aws.String("*"),
	}, func(o *s3.PresignOptions) {
		o.Expires = time.Duration(lifetimeSecs) * time.Second

	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

// Generate a random filename
func getRandomFilename() string {
	return fmt.Sprintf("%d.png", rand.Intn(100000))
}

// Generate a presigned PUT URL
func signPostHandler(w http.ResponseWriter, r *http.Request) {
	// Allow all origins
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Allow specific HTTP methods
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// Allow specific headers
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, x-amz-acl")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	key := getRandomFilename()

	contentlength := r.URL.Query().Get("content-length")
	c, _ := strconv.Atoi(contentlength)
	hash := r.URL.Query().Get("hash")
	client := NewS3Client(S3_ACCESS_KEY, S3_SECRET_KEY, S3_REGION)
	re, err := client.PutObject(S3_BUCKET_NAME, *aws.String(key), 100, c, hash)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("failed"))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(re.URL)
}

func main() {
	http.HandleFunc("/sign_post", signPostHandler)
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
