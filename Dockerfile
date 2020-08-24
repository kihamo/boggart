FROM golang:1.15 as builder

COPY . /go/src/github.com/kihamo/boggart/

# RUN go mod download
# RUN go build -v -i -race -o boggart

#RUN cd /go/src/github.com/kihamo/boggart/cmd/agent && go build -mod vendor -v -i -race -o boggart
RUN wget https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip \
    && cd /go/src/github.com/kihamo/boggart/cmd/server \
    && BUILD_DATE=$(date +"%y%m%d") \
    && BUILD_TIME=$(date +"%H%M%S") \
    && go build -mod vendor -v -i -race -o boggart -ldflags="-X 'main.Name=Boggart Server' -X 'main.Version=$BUILD_DATE' -X 'main.Build=$BUILD_TIME'" ./

FROM debian:bullseye-slim as server
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates \
 && update-ca-certificates
COPY --from=builder /go/zoneinfo.zip .
ENV ZONEINFO /zoneinfo.zip
COPY --from=builder /go/src/github.com/kihamo/boggart/cmd/server/boggart .
ENTRYPOINT ["./boggart"]