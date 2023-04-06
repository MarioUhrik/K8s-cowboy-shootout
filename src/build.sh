protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

docker build --tag k8s-cowboy-shootout/cowboy .