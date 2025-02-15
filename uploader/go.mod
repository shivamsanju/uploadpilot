module github.com/uploadpilot/uploadpilot/uploader

go 1.23.2

replace github.com/uploadpilot/uploadpilot/go-core/db v0.0.0 => ../go-core/db

replace github.com/uploadpilot/uploadpilot/go-core/pubsub v0.0.0 => ../go-core/pubsub

replace github.com/uploadpilot/uploadpilot/go-core/dsl v0.0.0 => ../go-core/dsl

replace github.com/uploadpilot/uploadpilot/go-core/common v0.0.0 => ../go-core/common

require (
	github.com/aws/aws-sdk-go-v2/config v1.29.6
	github.com/aws/aws-sdk-go-v2/credentials v1.17.59
	github.com/aws/aws-sdk-go-v2/service/s3 v1.76.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/render v1.0.3
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/redis/go-redis/v9 v9.7.0
	github.com/tus/tusd/v2 v2.6.0
	github.com/uploadpilot/uploadpilot/go-core/common v0.0.0
	github.com/uploadpilot/uploadpilot/go-core/db v0.0.0
	github.com/uploadpilot/uploadpilot/go-core/dsl v0.0.0
	github.com/uploadpilot/uploadpilot/go-core/pubsub v0.0.0
	go.temporal.io/sdk v1.32.1
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.67.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2 v1.36.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.8 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.5.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.14 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.24.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nexus-rpc/sdk-go v0.1.0 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/saracen/fastzip v0.1.11 // indirect
	github.com/saracen/zipextra v0.0.0-20220303013732-0187cb0159ea // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.temporal.io/api v1.43.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/exp v0.0.0-20231127185646-65229373498e // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241007155032-5fefd90f89a9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gorm.io/driver/postgres v1.5.11 // indirect
	gorm.io/gorm v1.25.12 // indirect
)
