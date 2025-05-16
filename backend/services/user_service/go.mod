module github.com/VitaliySynytskyi/survey-platform/backend/services/user_service

go 1.21

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-chi/cors v1.2.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.3.1
	github.com/jackc/pgx/v5 v5.4.3
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.13.0
)

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	golang.org/x/exp v0.0.0-20230817173708-d852ddb80c63 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

require (
	github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul v0.0.0
	github.com/hashicorp/consul/api v1.26.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)

replace github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul => ../../pkg/consul
