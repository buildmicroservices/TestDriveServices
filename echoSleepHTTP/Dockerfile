FROM golang:1.14.4-stretch

WORKDIR /go/bin

WORKDIR /go/src
COPY ./ /go/src/github.com/buildmicroservices/TestDriveServices
COPY ./devOps/build.sh .
COPY ./cfg/.cobra.yaml ~/.cobra.yaml
RUN ./build.sh
RUN rm ./build.sh
RUN rm -rf ./

CMD ["/go/bin/echoSleepHTTP"]
