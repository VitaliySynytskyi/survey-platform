module github.com/VitaliySynytskyi/survey-platform/backend/api_gateway

go 1.21

require (
	github.com/go-chi/chi/v5 v5.0.11
	github.com/go-chi/cors v1.2.1
	github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul v0.0.0
	github.com/VitaliySynytskyi/survey-platform/backend/pkg/tracing v0.0.0
	golang.org/x/time v0.5.0
)

replace github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul => ../pkg/consul
replace github.com/VitaliySynytskyi/survey-platform/backend/pkg/tracing => ../pkg/tracing 