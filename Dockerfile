FROM golang:1.15 as builder

COPY . /go/src/github.com/kihamo/boggart/

# RUN go mod download
# RUN go build -v -i -race -o boggart

#RUN cd /go/src/github.com/kihamo/boggart/cmd/agent && go build -mod vendor -v -i -race -o boggart
RUN cd /go/src/github.com/kihamo/boggart/cmd/server && \
    BUILD_DATE=$(date +"%y%m%d") && \
    BUILD_TIME=$(date +"%H%M%S") && \
    go build -mod vendor -v -i -race -o boggart -ldflags="-X 'main.Name=Boggart Server' -X 'main.Version=$BUILD_DATE' -X 'main.Build=$BUILD_TIME'" ./

#FROM scratch as production
#FROM golang:latest as agent
#COPY --from=builder /go/src/github.com/kihamo/boggart/cmd/agent/boggart .
#ENTRYPOINT ["./boggart"]
#EXPOSE 80

FROM golang:latest as server
COPY --from=builder /go/src/github.com/kihamo/boggart/cmd/server/boggart .
ENTRYPOINT ["./boggart"]