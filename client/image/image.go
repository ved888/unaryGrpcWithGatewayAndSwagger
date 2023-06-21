package image

import (
	"context"
	"encoding/json"
	"fmt"
	"gRPC2/model"
	"gRPC2/pb/pb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

var Conn3, _ = grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))

var Client3 = pb.NewImageServiceClient(Conn3)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	var image model.Image
	// Parse the multipart form data
	err := r.ParseMultipartForm(32 << 20) // Specify the maximum form size if needed
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("file") // "file" should match the name attribute of the file input field in your HTML form
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a unique file name
	fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), handler.Filename)

	// Create an AWS session with credentials and set the region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAW7S7LYQCP3KHWFLO",
			"i4nMMX/R+BR+OQf7znHMzp3rJQXsjOYrbyHMlUf/",
			""),
	})
	if err != nil {
		http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
		log.Println("Failed to create AWS session:", err)
		return
	}

	// Create an S3 client
	svc := s3.New(sess)

	// Specify the S3 bucket name
	bucketName := "ved-bucket-1"

	// Upload the file to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
		log.Println("Failed to upload file to S3:", err)
		return
	}
	ctx := context.Background()
	res, err := Client3.UploadImage(ctx, &pb.UploadImageRequest{
		ImageId:    image.Id,
		BucketName: image.BucketName,
		ImagePath:  []byte(image.Path),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(res.GetImageUrl())
	if err != nil {
		return
	}
	// Respond with a success message
	fmt.Fprintf(w, "File uploaded successfully")
}
