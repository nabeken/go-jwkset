package jwkset

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/patrickmn/go-cache"
	"gopkg.in/square/go-jose.v2"
)

// Response represents a response of JWK Set.
// This contains a TTL (Time to Live) for caching purpose.
type Response struct {
	Keys []jose.JSONWebKey

	TTL time.Duration // This would be used as TTL for caching.
}

// ALBFetcher fetchs a public key from AWS's Application Load Balancer and decodes it into JWK.
type ALBFetcher struct {
	Client *http.Client
	Region string
	Algo   jose.SignatureAlgorithm
}

func (f *ALBFetcher) keyURL(kid string) string {
	// https://docs.aws.amazon.com/elasticloadbalancing/latest/application/listener-authenticate-users.html
	return fmt.Sprintf("https://public-keys.auth.elb.%s.amazonaws.com/%s", f.Region, kid)
}

func (f *ALBFetcher) FetchJWKs(kid string) (*Response, error) {
	resp, err := f.Client.Get(f.keyURL(kid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jwks, err := DecodeSigPublicKey(data, kid, f.Algo)
	return &Response{
		Keys: jwks,
	}, err
}

// HTTPFetcher fetches JWKs over HTTP.
type HTTPFetcher struct {
	Client *http.Client
}

// InMemoryFetcher fetches JWKs from its memory.
type InMemoryFetcher struct {
	RAWJWKs []byte
}

// FetchJWKs implements Fetcher interface by using internal JWKs.
func (f *InMemoryFetcher) FetchJWKs(_ string) (*Response, error) {
	jwks, err := Decode(bytes.NewReader(f.RAWJWKs))
	if err != nil {
		return nil, err
	}
	return &Response{
		Keys: jwks,
	}, nil
}

// FetchJWKs implements Fetcher interface by using http.Client.
// FetchJWKs tries to retrieve JWKSet from uri.
func (f *HTTPFetcher) FetchJWKs(uri string) (*Response, error) {
	resp, err := f.Client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jwks, err := Decode(resp.Body)
	return &Response{
		Keys: jwks,
	}, err
}

// S3Fetcher fetches JWKs via S3.
type S3Fetcher struct {
	S3Svc s3iface.S3API
}

// FetchJWKs implements JWKsS3Fetcher by using S3. It tries to retrieve an S3 object from path.
// path must be in s3://<bucket>/<key>.
func (f *S3Fetcher) FetchJWKs(path string) (*Response, error) {
	s3url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	params := &s3.GetObjectInput{
		Bucket: aws.String(s3url.Host),
		Key:    aws.String(s3url.Path),
	}
	resp, err := f.S3Svc.GetObject(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jwks, err := Decode(resp.Body)
	return &Response{
		Keys: jwks,
	}, err
}

// Cacher fetches JWKs via Cache if available.
type Cacher struct {
	fetcher Fetcher
	cache   *cache.Cache

	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

// NewCacher returns Cacher with initializing cache store.
func NewCacher(defaultExpiration, cleanupInterval time.Duration, f Fetcher) *Cacher {
	c := cache.New(defaultExpiration, cleanupInterval)
	return &Cacher{
		fetcher: f,
		cache:   c,

		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
}

// FetchJWKs tries to retrieve JWKs from Cache. If the cache is not available,
// it will call Fetcher.FetchJWKs and cache the result for future request.
func (c *Cacher) FetchJWKs(cacheKey string) (*Response, error) {
	if keys, found := c.cache.Get(cacheKey); found {
		return &Response{Keys: keys.([]jose.JSONWebKey)}, nil
	}
	jwksresp, err := c.fetcher.FetchJWKs(cacheKey)
	if err != nil {
		return nil, err
	}

	ttl := jwksresp.TTL

	// If TTL is larger than cleanupInterval, we should subtract cleanupInterval from TTL to
	// make sure that the latest jwks is obtained.
	if ttl > c.cleanupInterval {
		ttl -= c.cleanupInterval
	}

	c.cache.Set(cacheKey, jwksresp.Keys, ttl)
	return jwksresp, nil
}

// DecodeSigPublicKey decodes the plain public key into JWKs used for sigining.
// https://github.com/square/go-jose/blob/v2.4.1/jose-util/utils.go#L42
func DecodeSigPublicKey(data []byte, kid string, algo jose.SignatureAlgorithm) ([]jose.JSONWebKey, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	pub, err := x509.ParsePKIXPublicKey(input)
	if err != nil {
		return nil, err
	}

	return []jose.JSONWebKey{
		{
			Key:       pub,
			KeyID:     kid,
			Algorithm: string(algo),
			Use:       "sig",
		},
	}, nil
}

// Decode decodes the data with reading from r into JWKs.
func Decode(r io.Reader) ([]jose.JSONWebKey, error) {
	keyset := jose.JSONWebKeySet{}
	if err := json.NewDecoder(r).Decode(&keyset); err != nil && err != io.EOF {
		return nil, err
	}

	return keyset.Keys, nil
}
