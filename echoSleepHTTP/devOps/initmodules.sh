#export GOPATH=${PWD}
#export GOBIN=$GOPATH/bin
export GO111MODULE=on
cd src
go mod init main.go
go mod vendor

### README

#cd into the ./src director with main.go
#run:
#$go mod init main.go
#$go mod vendor
