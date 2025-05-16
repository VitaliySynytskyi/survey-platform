module github.com/VitaliySynytskyi/survey-platform/backend/services/analytics_service

go 1.21

require (
	github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul v0.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/cors v1.2.1
	github.com/google/uuid v1.6.0
	go.mongodb.org/mongo-driver v1.17.3
)

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/consul/api v1.24.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
)

replace github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul => ../../pkg/consul
