module github.com/rivelles/rinha-backend-2024-q1/adapters

go 1.22.0

replace github.com/rivelles/rinha-backend-2024-q1/model => ../model

require (
	github.com/gocql/gocql v1.6.0
	github.com/redis/go-redis/v9 v9.4.0
	github.com/rivelles/rinha-backend-2024-q1/model v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/kr/text v0.2.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)
