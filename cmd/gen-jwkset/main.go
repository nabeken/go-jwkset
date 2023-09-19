package main

import (
	"crypto"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-jose/go-jose/v3"
)

func main() {
	var flagAlg = flag.String("alg", "", "specify an algorithm")
	flag.Parse()

	if *flagAlg == "" {
		flag.Usage()
		return
	}

	rawPublicKey, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("unable to read a public key from stdin:", err)
	}

	block, _ := pem.Decode(rawPublicKey)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal("unable to load a public key:", err)
	}

	hasher := crypto.SHA256.New()
	fmt.Fprint(hasher, rawPublicKey)
	hashed := hex.EncodeToString(hasher.Sum(nil))

	// TODO(nabeken): be able to specify multiple keys
	jwks := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:       publicKey,
				KeyID:     hashed,
				Algorithm: *flagAlg,
			},
		},
	}

	json.NewEncoder(os.Stdout).Encode(jwks)
}
