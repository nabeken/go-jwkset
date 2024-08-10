package jwkset

import "context"

// Fetcher is an interface that represents JWKs fetcher.
type Fetcher interface {
	// FetchJWKs retrieves JWKSet from path.
	FetchJWKs(ctxt context.Context, path string) (*Response, error)
}
