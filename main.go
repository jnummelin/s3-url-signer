package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/jnummelin/s3-url-signer/version"
)

var urlSignRequest UrlSignRequest

type UrlSignRequest struct {
	region   string
	bucket   string
	key      string
	verb     string
	duration string
}

// Validates the Signing request struct
func (signReq *UrlSignRequest) Validate() error {
	if len(signReq.region) == 0 {
		return fmt.Errorf("Region must be given")
	}
	if len(signReq.bucket) == 0 {
		return fmt.Errorf("Bucket must be given")
	}
	if len(signReq.key) == 0 {
		return fmt.Errorf("Key must be given")
	}
	if len(signReq.verb) == 0 {
		return fmt.Errorf("Verb must be given")
	}
	if len(signReq.duration) == 0 {
		return fmt.Errorf("Duration must be given")
	}
	// Check that the duration is parseable
	_, err := time.ParseDuration(urlSignRequest.duration)
	if err != nil {
		return err
	}

	return nil
}

func parseFlags() {
	flag.StringVar(&urlSignRequest.region, "region", "", "AWS Region")
	flag.StringVar(&urlSignRequest.bucket, "bucket", "", "Bucket")
	flag.StringVar(&urlSignRequest.key, "key", "", "Key")
	flag.StringVar(&urlSignRequest.verb, "verb", "", "HTTP Verb (GET, PUT, DELETE)")
	flag.StringVar(&urlSignRequest.duration, "duration", "15m", "Expiry time of the URL")

	flag.Parse()

	if err := urlSignRequest.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(13)
	}
}

// Pre-signs S3 object URLs
func main() {
	for _, arg := range os.Args {
		if arg == "-v" {
			fmt.Println(version.BuildVersion())
			os.Exit(0)
		}
	}

	parseFlags()

	// Initialize a session in given region
	// SDK will use to load credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(urlSignRequest.region)},
	)

	// Create S3 service client
	svc := s3.New(sess)

	var req *request.Request

	switch urlSignRequest.verb {
	case "GET":
		req, _ = svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(urlSignRequest.bucket),
			Key:    aws.String(urlSignRequest.key),
		})
	case "PUT":
		req, _ = svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(urlSignRequest.bucket),
			Key:    aws.String(urlSignRequest.key),
		})
	case "DELETE":
		req, _ = svc.DeleteObjectRequest(&s3.DeleteObjectInput{
			Bucket: aws.String(urlSignRequest.bucket),
			Key:    aws.String(urlSignRequest.key),
		})
	default:
		fmt.Fprintf(os.Stderr, "Unsupported HTTP verb: %s\nOnly GET, PUT and DELETE is supported!\n", urlSignRequest.verb)
		os.Exit(1)
	}

	dur, err := time.ParseDuration(urlSignRequest.duration)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse duration: %s\n", err)
		os.Exit(5)
	}

	urlStr, err := req.Presign(dur)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to sign request: %s", err)
		os.Exit(10)
	}

	fmt.Println(urlStr)
}
