CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o sunshine
docker build -t sunshine:v1.0 .
echo "sunshine:v1.0"
