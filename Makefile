
mock:
	mockgen -destination=internal/mocks/mock_system_service.go -package=mocks metrics/internal/core/service Pinger	

test:
	go test ./...

test-cov:
	go test -v -coverpkg=./... -coverprofile=profile.cov.tmp ./...
	grep -Ev "mocks|migrations" profile.cov.tmp > profile.cov
	go tool cover -func profile.cov

server:
	go run cmd/keeper/main.go -l debug -r http://localhost:8090 -d "host=localhost port=35432 user=username password=password dbname=keeper sslmode=disable"

migrate-new:
	goose create $@ sql

migrate-up:
	goose -dir migrations postgres "user=username dbname=keeper password=password sslmode=disable host=127.0.0.1 port=35432" up

migrate-down:
	goose -dir migrations postgres "user=username dbname=keeper password=password sslmode=disable host=127.0.0.1 port=35432" down

keys:
	openssl genrsa -out private.pem 4096
	openssl rsa -in private.pem -outform PEM -pubout -out public.pem