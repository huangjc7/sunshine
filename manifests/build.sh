CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o sunshine
docker build -t jcregistry/sunshine:v1.0 .
echo "jcregistry/sunshine:v1.0"
