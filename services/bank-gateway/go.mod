module github.com/omniroute/bank-gateway

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.5.1
	github.com/redis/go-redis/v9 v9.4.0
	github.com/shopspring/decimal v1.3.1
	go.opentelemetry.io/otel v1.22.0
	go.opentelemetry.io/otel/trace v1.22.0
	go.uber.org/zap v1.26.0
)

// Local development - remove for production
replace github.com/omniroute/bank-gateway => ./
