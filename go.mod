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

require (
	github.com/antonybholmes/go-basemath v0.0.0-20250307171544-f522e91d1448 // indirect
	github.com/antonybholmes/go-dna v0.0.0-20250314164202-2039d499c1f1

)

require (
	github.com/antonybholmes/go-genome v0.0.0-20250304154943-58f0bcc07ad8
	github.com/antonybholmes/go-web v0.0.0-20250314222734-5e8338d8737d
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20250303144057-3b4a79366f5c
	github.com/antonybholmes/go-sys v0.0.0-20250307171544-adff7a5ac327
	github.com/gorilla/sessions v1.4.0 // indirect
)

require (
	github.com/antonybholmes/go-beds v0.0.0-20250303144057-c3c36557b1f8
	github.com/antonybholmes/go-geneconv v0.0.0-20250303144105-bad8002363e0
	github.com/antonybholmes/go-hubs v0.0.0-00010101000000-000000000000
	github.com/antonybholmes/go-math v0.0.0-20250303144104-ebe30afaed11
	github.com/antonybholmes/go-motifs v0.0.0-20250303144101-5e20183355ad
	github.com/antonybholmes/go-mutations v0.0.0-20250303144101-dbfe22a3a9c2
	github.com/antonybholmes/go-pathway v0.0.0-20250303144056-fb8a534e96e1
	github.com/antonybholmes/go-seqs v0.0.0-20250303144058-0f2f53d39239
	github.com/gin-contrib/cors v1.7.3
	github.com/gin-contrib/sessions v1.0.2
	github.com/gin-gonic/gin v1.10.0
	github.com/redis/go-redis/v9 v9.7.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.36.3 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.29.9 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.62 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.43.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.17 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/bytedance/sonic v1.13.1 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v1.0.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.25.0 // indirect
	github.com/go-sql-driver/mysql v1.9.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/segmentio/kafka-go v0.4.47 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/xuri/efp v0.0.0-20250227110027-3491fafc2b79 // indirect
	github.com/xuri/excelize/v2 v2.9.0 // indirect
	github.com/xuri/nfp v0.0.0-20250226145837-86d5fc24b2ba // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/antonybholmes/go-cytobands v0.0.0-20250303144056-0053f66ec6b7
	github.com/antonybholmes/go-gex v0.0.0-20250304155622-243c3c7e7ae8
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/xyproto/randomstring v1.2.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)
