module github.com/antonybholmes/go-edb-server-gin

go 1.24

replace github.com/antonybholmes/go-web => ../go-web

replace github.com/antonybholmes/go-dna => ../go-dna

replace github.com/antonybholmes/go-genome => ../go-genome

replace github.com/antonybholmes/go-mutations => ../go-mutations

replace github.com/antonybholmes/go-basemath => ../go-basemath

replace github.com/antonybholmes/go-sys => ../go-sys

replace github.com/antonybholmes/go-mailer => ../go-mailer

replace github.com/antonybholmes/go-geneconv => ../go-geneconv

replace github.com/antonybholmes/go-motifs => ../go-motifs

replace github.com/antonybholmes/go-pathway => ../go-pathway

replace github.com/antonybholmes/go-gex => ../go-gex

replace github.com/antonybholmes/go-seqs => ../go-seqs

replace github.com/antonybholmes/go-cytobands => ../go-cytobands

replace github.com/antonybholmes/go-beds => ../go-beds

replace github.com/antonybholmes/go-hubs => ../go-hubs

replace github.com/antonybholmes/go-scrna => ../go-scrna

require (
	github.com/antonybholmes/go-basemath v0.0.0-20250606171604-5853de3754da // indirect
	github.com/antonybholmes/go-dna v0.0.0-20250606171555-80d2c61ab8da

)

require (
	github.com/antonybholmes/go-genome v0.0.0-20250606171547-b52cfc3a1fa1
	github.com/antonybholmes/go-web v0.0.0-20250606171556-2e45aaa13f77
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.28
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20250606171552-7c08a9479c50
	github.com/antonybholmes/go-sys v0.0.0-20250606171605-31639110750b
	github.com/gorilla/sessions v1.4.0 // indirect
)

require (
	github.com/antonybholmes/go-beds v0.0.0-20250606171551-8661cae8ced9
	github.com/antonybholmes/go-geneconv v0.0.0-20250606171606-e65caa62f8db
	github.com/antonybholmes/go-hubs v0.0.0-00010101000000-000000000000
	github.com/antonybholmes/go-math v0.0.0-20250606171604-5853de3754da
	github.com/antonybholmes/go-motifs v0.0.0-20250606171559-5a90b23136d7
	github.com/antonybholmes/go-mutations v0.0.0-20250606171558-b181833245ac
	github.com/antonybholmes/go-pathway v0.0.0-20250606171548-ce8aebd7de15
	github.com/antonybholmes/go-scrna v0.0.0-20250606171553-36fd03e6ed51
	github.com/antonybholmes/go-seqs v0.0.0-20250606171554-8409a199278b
	github.com/gin-contrib/cors v1.7.5
	github.com/gin-contrib/sessions v1.0.4
	github.com/gin-gonic/gin v1.10.1
	github.com/redis/go-redis/v9 v9.10.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.36.3 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.29.15 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.68 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.45.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.20 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/bytedance/sonic v1.13.3 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.26.0 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/segmentio/kafka-go v0.4.48 // indirect
	github.com/tiendc/go-deepcopy v1.6.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.14 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/excelize/v2 v2.9.1 // indirect
	github.com/xuri/nfp v0.0.1 // indirect
	golang.org/x/arch v0.18.0 // indirect
	golang.org/x/exp v0.0.0-20250606033433-dcc06ee1d476 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/antonybholmes/go-cytobands v0.0.0-20250606171549-5392011a60f4
	github.com/antonybholmes/go-gex v0.0.0-20250528204356-0bcd3ac6073d
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/xyproto/randomstring v1.2.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)
