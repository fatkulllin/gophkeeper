PROTO_SRC=api/proto
PROTO_OUT=api/gophkeeperpb

generate-proto:
	protoc \
		--go_out=$(PROTO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)/*.proto

run:
	go run ./cmd/server/main.go

grpc-health:
	grpc-health-probe -addr=localhost:9090

up-docker:
	docker compose -f ./tools/docker-compose.yaml up  -d
