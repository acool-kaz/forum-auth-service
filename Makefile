gen:
	protoc --go_out=. \
		--go-grpc_out=. \
		 --experimental_allow_proto3_optional \
		./proto/*.proto

compose-build:
	docker-compose build

compose-up:
	docker-compose up