
cd $GOPATH/src/github.com/liweizhi/containerPool/controller
rm controller
echo build target..
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
echo build completed
cd ..
docker rmi -f li41898/container-pool
docker build -t li41898/container-pool .
docker push li41898/container-pool

