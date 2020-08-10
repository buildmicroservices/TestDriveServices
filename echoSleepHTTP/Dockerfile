FROM golang:1.14.4-stretch

WORKDIR /go/bin
WORKDIR /go/src
ENV GOBIN=/go/bin
ENV GOPATH=/go/src

COPY ./ $GOBIN
COPY ./devOps/build.sh .
COPY ./cfg/.cobra.yaml ~/.cobra.yaml
RUN ./build.sh
RUN rm ./build.sh
RUN rm -rf ./
COPY ./devOps/runLocal.sh $GOBIN/runLocal.sh
CMD ["$GOBIN/runLocal.sh"]
