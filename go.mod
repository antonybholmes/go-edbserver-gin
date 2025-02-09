module github.com/antonybholmes/go-edb-server-gin

go 1.23

replace github.com/antonybholmes/go-auth => ../go-auth

replace github.com/antonybholmes/go-dna => ../go-dna

replace github.com/antonybholmes/go-genes => ../go-genes

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

require (
	github.com/antonybholmes/go-basemath v0.0.0-20250205152412-840349f1ca5c // indirect
	github.com/antonybholmes/go-dna v0.0.0-20250205171514-82dd7b47fd97

)

require (
	github.com/antonybholmes/go-auth v0.0.0-20250205171514-f8d791f6a0e3
	github.com/antonybholmes/go-genes v0.0.0-20250205171518-74f5823b64be
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20250205171512-d98bb8e5219c
	github.com/antonybholmes/go-sys v0.0.0-20250205171520-3b66c82ca1f3
	github.com/gorilla/sessions v1.4.0 // indirect
)

require (
	github.com/antonybholmes/go-beds v0.0.0-20250205171511-27c9db460426
	github.com/antonybholmes/go-geneconv v0.0.0-20250205171520-978d384f5d4f
	github.com/antonybholmes/go-math v0.0.0-20250205152412-840349f1ca5c
	github.com/antonybholmes/go-motifs v0.0.0-20250205171516-ec338fab0afc
	github.com/antonybholmes/go-mutations v0.0.0-20250205171516-4a5a25eecb96
	github.com/antonybholmes/go-pathway v0.0.0-20250205171510-54e15ded0a64
	github.com/antonybholmes/go-seqs v0.0.0-20250205171513-452206075530
	github.com/gin-contrib/cors v1.7.3
	github.com/gin-contrib/sessions v1.0.2
	github.com/gin-gonic/gin v1.10.0
	github.com/redis/go-redis/v9 v9.7.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/bytedance/sonic v1.12.8 // indirect
	github.com/bytedance/sonic/loader v0.2.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v1.0.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.24.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/xuri/efp v0.0.0-20241211021726-c4e992084aa6 // indirect
	github.com/xuri/excelize/v2 v2.9.0 // indirect
	github.com/xuri/nfp v0.0.0-20250111060730-82a408b9aa71 // indirect
	golang.org/x/arch v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20250207012021-f9890c6ad9f3 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/antonybholmes/go-cytobands v0.0.0-20250205171511-1b160eec2646
	github.com/antonybholmes/go-gex v0.0.0-20250205171515-dfe453ea0c91
	github.com/aws/aws-sdk-go v1.55.6 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/xyproto/randomstring v1.2.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)
