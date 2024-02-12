module github.com/rivelles/rinha-backend-2024-q1/main

go 1.22.0

replace github.com/rivelles/rinha-backend-2024-q1/application => ../application

replace github.com/rivelles/rinha-backend-2024-q1/model => ../model

require (
	github.com/rivelles/rinha-backend-2024-q1/adapters v0.0.0-00010101000000-000000000000
	github.com/rivelles/rinha-backend-2024-q1/application v0.0.0-00010101000000-000000000000
)

require (
	github.com/gocql/gocql v1.6.0 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/rivelles/rinha-backend-2024-q1/model v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace github.com/rivelles/rinha-backend-2024-q1/adapters => ../adapters
