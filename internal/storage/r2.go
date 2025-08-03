/*
package storage


import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Storage() {
	accountID := ""
	accessKeyID := ""
	accessKeySecret := ""

	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:               r2Endpoint,
				SigningRegion:     "auto", // R2 uses a global region, "auto" is a good default
				HostnameImmutable: true,
			}, nil
		}
		// Fallback to the default resolver for other services
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
		config.WithRegion("auto"), // Set a dummy region, as R2 doesn't use AWS regions
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
}
*/
