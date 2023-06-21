package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"uploadImage/database"
	"uploadImage/dbhelper"
	"uploadImage/model"
	"uploadImage/pb/pb"
)

var (
	port = flag.Int("port", 50053, "gRPC server port")
)

type server struct {
	DAO dbhelper.DAO
	pb.UnimplementedImageServiceServer
}

func main() {
	flag.Parse()

	db, err := database.DbConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Close()

	fmt.Println("gRPC server running ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterImageServiceServer(s, &server{DAO: dbhelper.DAO{
		DB: db,
	},
	})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	// Load environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("UploadImage: Error loading .env file: %w", err)
	}

	imageData := req.ImagePath
	fileName := req.ImageId + ".jpg" // Generate a unique file name for the image

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			""),
	})
	if err != nil {
		return nil, fmt.Errorf("UploadImage: Error creating AWS session: %w", err)
	}

	svc := s3.New(sess)

	bucketName := os.Getenv("S3_BUCKET_NAME")

	// Convert the []byte data to an io.Reader
	imageReader := bytes.NewReader(imageData)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   imageReader,
	})
	if err != nil {
		return nil, fmt.Errorf("UploadImage: Error uploading file to S3: %w", err)
	}

	image := model.Image{
		BucketName: bucketName,
		Path:       fileName,
	}

	// Call the UploadImage function in the DAO
	uploadID, err := s.DAO.UploadImage(image)
	if err != nil {
		return nil, fmt.Errorf("UploadImage: Error uploading image: %w", err)
	}

	if uploadID == nil {
		return nil, errors.New("UploadImage: Failed to get upload ID")
	}

	response := &pb.UploadImageResponse{
		ImageUrl: *uploadID,
	}

	return response, nil
}
