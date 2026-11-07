module github.com/nabeken/go-jwkset/v3

go 1.24.0

toolchain go1.27.3

require (
	github.com/aws/aws-sdk-go-v2 v1.47.8
	github.com/aws/aws-sdk-go-v2/service/s3 v1.114.9
	github.com/go-jose/go-jose/v4 v4.1.9
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/stretchr/testify v1.12.1
	go.uber.org/mock v0.6.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.26 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.20.11 // indirect
	github.com/aws/smithy-go v1.28.8 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
)

tool go.uber.org/mock/mockgen
