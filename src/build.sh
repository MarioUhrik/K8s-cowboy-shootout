protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

cd cowboy
docker build --tag k8s-cowboy-shootout/cowboy .
cd -

cd referee
docker build --tag k8s-cowboy-shootout/referee .
cd -

