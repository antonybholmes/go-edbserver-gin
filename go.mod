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
	github.com/antonybholmes/go-basemath v0.0.0-20250507224209-326910455aee // indirect
	github.com/antonybholmes/go-dna v0.0.0-20250509222604-3dd9e6aa9ab6

)

require (
	github.com/antonybholmes/go-genome v0.0.0-20250326211422-64d18aac02e1
	github.com/antonybholmes/go-web v0.0.0-20250513135653-c491eaeedab6
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.28
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20250326211424-cc7dd35c3d80
	github.com/antonybholmes/go-sys v0.0.0-20250512211932-464c953347d7
	github.com/gorilla/sessions v1.4.0 // indirect
)

require (
	github.com/antonybholmes/go-beds v0.0.0-20250326211424-2f11e9fe884e
	github.com/antonybholmes/go-geneconv v0.0.0-20250326211432-83664704cef7
	github.com/antonybholmes/go-hubs v0.0.0-00010101000000-000000000000
	github.com/antonybholmes/go-math v0.0.0-20250307171544-f522e91d1448
	github.com/antonybholmes/go-motifs v0.0.0-20250326211428-91561de16b6d
	github.com/antonybholmes/go-mutations v0.0.0-20250326211428-0c2813bd22dc
	github.com/antonybholmes/go-pathway v0.0.0-20250326211422-45d0c93b0e39
	github.com/antonybholmes/go-scrna v0.0.0-00010101000000-000000000000
	github.com/antonybholmes/go-seqs v0.0.0-20250326211425-a6dc0f9f73fe
	github.com/gin-contrib/cors v1.7.5
	github.com/gin-contrib/sessions v1.0.2
	github.com/gin-gonic/gin v1.10.1
	github.com/redis/go-redis/v9 v9.9.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.36.3 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.29.14 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.67 // indirect
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
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.19 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/bytedance/sonic v1.13.2 // indirect
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
	github.com/tiendc/go-deepcopy v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.14 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/excelize/v2 v2.9.1 // indirect
	github.com/xuri/nfp v0.0.1 // indirect
	golang.org/x/arch v0.17.0 // indirect
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/antonybholmes/go-cytobands v0.0.0-20250326211423-94b98ba66052
	github.com/antonybholmes/go-gex v0.0.0-20250328152020-67df8035632d
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/xyproto/randomstring v1.2.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)
