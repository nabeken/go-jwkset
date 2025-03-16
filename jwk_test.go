package jwkset

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestALBFetcher(t *testing.T) {
	assert := assert.New(t)
	fetcher := &ALBFetcher{
		Client: &http.Client{},
		Region: "ap-northeast-1",
		Algo:   jose.ES256,
	}
	jwksresp, err := fetcher.FetchJWKs(context.TODO(), "21a3e6e4-c32e-4650-b43d-813ba7628f3b")
	assert.NoError(err)
	assert.Len(jwksresp.Keys, 1)

	key := jwksresp.Keys[0]

	assert.Equal("ES256", key.Algorithm)
	_, ok := jwksresp.Keys[0].Key.(*ecdsa.PublicKey)
	if !assert.True(ok) {
		t.Logf("got '%#v'", jwksresp.Keys[0])
	}
}

func TestInMemoryFetcher(t *testing.T) {
	assert := assert.New(t)

	fetcher := &InMemoryFetcher{
		RAWJWKs: mustReadTestData("google.jwk"),
	}

	jwksresp, err := fetcher.FetchJWKs(context.TODO(), "")
	assert.NoError(err)
	assertTestGoogleJWK(t, jwksresp)
}

func TestS3Fetcher(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	s3m := NewMockS3API(ctrl)

	s3m.EXPECT().
		GetObject(context.TODO(), gomock.Any(), gomock.Any()).
		Return(&s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewReader(mustReadTestData("google.jwk"))),
		}, nil)

	fetcher := &S3Fetcher{
		S3client: s3m,
	}

	jwksresp, err := fetcher.FetchJWKs(context.TODO(), "")
	assert.NoError(err)
	assertTestGoogleJWK(t, jwksresp)
}

func TestFetcher(t *testing.T) {
	assert := assert.New(t)
	fetcher := &HTTPFetcher{
		Client: &http.Client{},
	}

	testURL := testServer(t)

	jwksresp, err := fetcher.FetchJWKs(context.TODO(), testURL)
	assert.NoError(err)
	assertTestGoogleJWK(t, jwksresp)
}

func TestJWKsCacher(t *testing.T) {
	defaultExpiration := 10 * time.Minute
	cleanupInterval := time.Minute

	assert := assert.New(t)
	cacher := NewCacher(defaultExpiration, cleanupInterval, &HTTPFetcher{
		Client: &http.Client{},
	})

	testURL := testServer(t)

	cacheKey := testURL
	jwksresp, err := cacher.FetchJWKs(context.TODO(), cacheKey)
	assert.NoError(err)
	assertTestGoogleJWK(t, jwksresp)

	cachedResp, found := cacher.cache.Get(cacheKey)
	assert.True(found)

	resp, ok := cachedResp.([]jose.JSONWebKey)
	if assert.True(ok, "cached response should be []jose.JSONWebKey but %#v", cachedResp) {
		assert.Equal(jwksresp.Keys, resp)
	}

	jwksresp, err = cacher.FetchJWKs(context.TODO(), cacheKey)
	assert.NoError(err)
	assertTestGoogleJWK(t, jwksresp)
}

func testServer(t *testing.T) string {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := mustReadTestData("google.jwk")
		_, _ = w.Write(data)
	}))

	t.Cleanup(func() { ts.Close() })

	return ts.URL
}

func assertTestGoogleJWK(t *testing.T, resp *Response) bool {
	assert := assert.New(t)

	return assert.Len(resp.Keys, 2) &&
		assert.Equal("c6263d09745b5032e57fa6e1d041b77a54066dbd", resp.Keys[0].KeyID) &&
		assert.Equal("7d334497506acb74cdeedaa66184d15547f83693", resp.Keys[1].KeyID)
}

func mustReadTestData(fn string) []byte {
	data, err := os.ReadFile(filepath.Join("_testdata", fn))
	if err != nil {
		panic(fmt.Errorf("reading a file from _testdata: %w", err))
	}

	return data
}
