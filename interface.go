//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -package=$GOPACKAGE -destination=interface_mock.go

package jwkset

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Fetcher is an interface that represents JWKs fetcher.
type Fetcher interface {
	// FetchJWKs retrieves JWKSet from path.
	FetchJWKs(ctxt context.Context, path string) (*Response, error)
}

// S3API is an interface that represents S3 API client.
type S3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}
