package jwkset

import (
	"crypto/ecdsa"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/assert"
)

func TestALBFetcher(t *testing.T) {
	assert := assert.New(t)
	fetcher := &ALBFetcher{
		Client: &http.Client{},
		Region: "ap-northeast-1",
		Algo:   jose.ES256,
	}
	jwksresp, err := fetcher.FetchJWKs("21a3e6e4-c32e-4650-b43d-813ba7628f3b")
	assert.NoError(err)
	assert.Len(jwksresp.Keys, 1)

	key := jwksresp.Keys[0]

	assert.Equal("ES256", key.Algorithm)
	_, ok := jwksresp.Keys[0].Key.(*ecdsa.PublicKey)
	if !assert.True(ok) {
		t.Logf("got '%#v'", jwksresp.Keys[0])
	}
}

func TestFetcher(t *testing.T) {
	assert := assert.New(t)
	fetcher := &HTTPFetcher{
		Client: &http.Client{},
	}

	testURL := testServer(t)

	jwksresp, err := fetcher.FetchJWKs(testURL)
	assert.NoError(err)
	assert.Len(jwksresp.Keys, 2)
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
	jwksresp, err := cacher.FetchJWKs(cacheKey)
	assert.NoError(err)
	assert.Len(jwksresp.Keys, 2)

	cachedResp, found := cacher.cache.Get(cacheKey)
	assert.True(found)

	resp, ok := cachedResp.([]jose.JSONWebKey)
	if assert.True(ok, "cached response should be []jose.JSONWebKey but %#v", cachedResp) {
		assert.Equal(jwksresp.Keys, resp)
	}

	jwksresp, err = cacher.FetchJWKs(cacheKey)
	assert.NoError(err)
	assert.Len(jwksresp.Keys, 2)
}

func testServer(t *testing.T) string {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("_testdata/google.jwk")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}))

	t.Cleanup(func() { ts.Close() })

	return ts.URL
}
