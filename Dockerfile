FROM golang:1.10 as builder
ARG VERSION
WORKDIR /go/src/github.com/xwjdsh/fy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-X main.version=${VERSION}" -o fy ./cmd/fy 

FROM alpine:latest  
LABEL maintainer="iwendellsun@gmail.com"
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/xwjdsh/fy/fy .
ENTRYPOINT ["/fy"]
