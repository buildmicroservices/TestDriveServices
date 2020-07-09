export GOPATH=${PWD}
export GOBIN=$GOPATH/bin
#go get -d -v ./...
#go install -v ./...
cd ./src
export GO111MODULE=on
go mod vendor
cd ..
export GO111MODULE=auto
go build -o ./bin/echoSleepHTTP ./src
